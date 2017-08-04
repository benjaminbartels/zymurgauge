package internal

import "time"

type Device struct {
	GPIO     int           `json:"gpio"`
	Cooldown time.Duration `json:"cooldown"`
}

func (d *Device) Equals(o *Device) bool {
	if d.GPIO != o.GPIO {
		return false
	} else if d.Cooldown != o.Cooldown {
		return false
	}

	return true
}
