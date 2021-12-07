package main

import (
	"math"
	"time"
)

type PIDAutoTuner struct {
	input            float64
	output           float64
	controlType      int
	noiseBand        float64
	isRunning        bool
	oStep            float64
	lookbackDuration time.Duration
	lastTime         time.Time
	//
	peakType    int
	peakCount   int
	sampleTime  time.Duration
	justchanged bool
	absMax      float64
	absMin      float64
	setpoint    float64
	outputStart float64
	nLookBack   int
	lastInputs  []float64
	peaks       []float64
	peak1       time.Time
	peak2       time.Time
	Ku          float64
	Pu          float64
}

func New(in, out float64) *PIDAutoTuner {
	return &PIDAutoTuner{
		input:            in,
		output:           out, // TODO: why?
		controlType:      0,   // default to PI
		noiseBand:        0.5,
		isRunning:        false,
		oStep:            30.0,
		lookbackDuration: 10 * time.Second,
		// lastTime: // TODO: initialize?
	}
}

func (p *PIDAutoTuner) Run() bool { // TODO: return error?
	if p.peakCount > 9 && p.isRunning {
		p.isRunning = false
		p.finish()

		return true
	}

	now := time.Now()

	if now.Sub(p.lastTime) < p.sampleTime {
		return false
	}

	p.lastTime = now

	refVal := p.input

	if !p.isRunning {
		// initialize working variables the first time around
		p.peakType = 0
		p.peakCount = 0
		p.justchanged = false
		p.absMax = refVal
		p.absMin = refVal
		p.setpoint = refVal
		p.isRunning = true
		p.outputStart = p.output
		p.output = p.outputStart + p.oStep
	} else {
		if refVal > p.absMax {
			p.absMax = refVal
		}
		if refVal < p.absMin {
			p.absMin = refVal
		}
	}

	// oscillate the output base on the input's relation to the setpoint
	if refVal > p.setpoint+p.noiseBand {
		p.output = p.outputStart - p.oStep
	} else if refVal < p.setpoint-p.noiseBand {
		p.output = p.outputStart + p.oStep
	}

	isMax := true
	isMin := true

	for i := p.nLookBack - 1; i >= 0; i-- {
		val := p.lastInputs[i]
		if isMax {
			isMax = refVal > val
		}

		if isMin {
			isMin = refVal < val
		}

		p.lastInputs[i+1] = p.lastInputs[i]
	}

	p.lastInputs[0] = refVal

	if p.nLookBack < 9 {
		// we don't want to trust the maxes or mins until the inputs array has been filled
		return false
	}

	if isMax {
		if p.peakType == 0 {
			p.peakType = 1
		}

		if p.peakType == -1 {
			p.peakType = 1
			p.justchanged = true
			p.peak2 = p.peak1
		}

		p.peak1 = now
		p.peaks[p.peakCount] = refVal
	} else if isMin {
		if p.peakType == 0 {
			p.peakType = -1
		}
		if p.peakType == 1 {
			p.peakType = -1
			p.peakCount++
			p.justchanged = true
		}

		if p.peakCount < 10 {
			p.peaks[p.peakCount] = refVal
		}
	}

	if p.justchanged && p.peakCount > 2 {
		// we've transitioned.  check if we can autotune based on the last peaks
		avgSeparation := (math.Abs(p.peaks[p.peakCount-1]-p.peaks[p.peakCount-2]) +
			math.Abs(p.peaks[p.peakCount-2]-p.peaks[p.peakCount-3])) / 2
		if avgSeparation < 0.05*(p.absMax-p.absMin) {
			p.finish()
			p.isRunning = false

			return true
		}
	}

	p.justchanged = false

	return false
}

func (p *PIDAutoTuner) finish() {
	p.output = p.outputStart
	// we can generate tuning parameters!
	p.Ku = 4 * (2 * p.oStep) / ((p.absMax - p.absMin) * 3.14159)
	p.Pu = (float64)(p.peak1.Second()-p.peak2.Second()) / 1000
}
