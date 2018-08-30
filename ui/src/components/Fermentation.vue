<template>
  <v-container>
	<!-- ID               uint64     `json:"id"`
	ChamberMacAddres uint64     `json:"chamberMacAddres"`
	Beer             Beer       `json:"beer"`
	CurrentStep      int        `json:"currentStep,omitempty"`
	StartTime        *time.Time `json:"startTime,omitempty"`
	CompletedTime    *time.Time `json:"completedTime,omitempty"`
	ModTime          time.Time  `json:"modTime"` -->

   <div v-if="fermentation">
      <v-form v-model="valid" ref="form" lazy-validation>
        <v-select v-if="create" label="Chambers" v-model="fermentation.chamber" v-bind:items="chambers" item-text="name" item.value="name" return-object required></v-select>
        <v-select v-if="create" label="Beers" v-model="fermentation.beer" v-bind:items="beers" item-text="name" item.value="name" return-object required></v-select>
        <v-text-field v-if="!create" label="Chambers" v-model="fermentation.chamber.name" disabled required></v-text-field>
        <v-text-field v-if="!create" label="Beers" v-model="fermentation.beer.name" disabled required></v-text-field>
        <v-btn v-if="create" @click="save" :disabled="!valid">save</v-btn>
        <v-btn @click="remove" :disabled="this.create">remove</v-btn>
        <v-btn v-if="!create" :disabled="this.create">start</v-btn>
        <v-btn v-if="!create" :disabled="this.create">stop</v-btn>
      </v-form>
    </div>

    <div v-if="loading">
      <p>Loading...</p>
    </div>

    <ul v-if="errors && errors.length">
      <li v-for="error of errors" :key="error.message">
        {{error.message}}
      </li>
    </ul>
    <temperature-chart :chart-data="chartData" :options="chartOptions"/>
  </v-container>
</template>

<script>
import { HTTP } from '../http-common'
import TemperatureChart from '../TemperatureChart'
import moment from 'moment'

export default {
  components: {
    TemperatureChart
  },
  data: () => ({
    fermentation: null,
    beers: [],
    chambers: [],
    temperaturechanges: [],
    chartData: null,
    chartOptions: {
      responsive: true,
      maintainAspectRatio: false,
      legend: {
        labels: {
          fontColor: '#ffffff'
        }
      },
      scales: {
        xAxes: [{
          type: 'time',
          ticks: {fontColor: '#ffffff'},
          time: {
            unit: 'second',
            unitStepSize: 1000
          }
        }],
        yAxes: [{
          ticks: {fontColor: '#ffffff'}
        }]
      }
    },
    errors: [],
    loading: false,
    valid: true,
    beerRules: [(v) => !!v || 'Beer is required'],
    chamberRules: [(v) => !!v || 'Chamber is required']
  }),

  props: ['id', 'create'],

  created () {
    if (!this.create) {
      this.fetch()
    } else {
      this.fermentation = {}
    }
  },
  mounted () {
    HTTP.get('fermentations/' + this.id + '/temperaturechanges')
      .then(response => {
        this.temperaturechanges = response.data

        var times = []
        var temps = []

        this.temperaturechanges.forEach(function (t) {
          times.push(moment(t.insertTime).format('MM/DD/YYYY h:mm:ss a'))
          temps.push(t.temperature)
        })

        this.chartData = {
          labels: times,
          datasets: [
            {
              label: 'Temperature',
              backgroundColor: '#FFB90D',
              data: temps
            }
          ]
        }
      })
      .catch(e => {
        this.errors.push(e)
      })
  },

  watch: {
    '$route': 'fetch'
  },

  methods: {
    fetch () {
      this.fermentation = null
      this.errors = []
      this.loading = true
      HTTP.get('beers')
        .then(response => {
          this.beers = response.data
        })
        .catch(e => {
          this.errors.push(e)
        })
      HTTP.get('chambers')
        .then(response => {
          this.chambers = response.data
        })
        .catch(e => {
          this.errors.push(e)
        })
      HTTP.get('fermentations/' + this.id)
        .then(response => {
          this.fermentation = response.data
        })
        .catch(e => {
          this.errors.push(e)
        })
        .finally(_ => {
          this.loading = false
        })
    },
    addSchedule () {
      if (this.fermentation.schedule == null) {
        this.fermentation.schedule = []
      }
      this.fermentation.schedule.push({})
    },
    removeSchedule (index) {
      this.fermentation.schedule.splice(index, 1)
    },
    save () {
      if (this.$refs.form.validate()) {
        HTTP.post('fermentations', this.fermentation)
          .then(response => {
            this.$router.push({ name: 'fermentations' })
          })
          .catch(e => {
            this.errors.push(e)
          })
      }
    },
    remove () {
      HTTP.delete('fermentations/' + this.fermentation.id)
        .then(response => {
          this.$router.push({ name: 'fermentations' })
        })
        .catch(e => {
          this.errors.push(e)
        })
    }
  }
}
</script>
