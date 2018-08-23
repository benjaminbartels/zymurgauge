<template>
  <v-container>
    <div v-if="chamber">
      <v-form v-model="valid" ref="form" lazy-validation>
        <v-text-field label="Name" v-model="chamber.name" :rules="nameRules"></v-text-field>
        <v-text-field label="Mac Address" v-model="chamber.macAddress" disabled></v-text-field>
        <v-text-field label="Thermometer ID" v-model="chamber.thermostat.thermometerId"></v-text-field>
        <v-text-field label="Cooler Pin" v-model="chamber.thermostat.chillerPin"></v-text-field>
        <v-text-field label="Heater Pin" v-model="chamber.thermostat.heaterPin"></v-text-field>
        <v-btn @click="save" :disabled="!valid">save</v-btn>
        <v-btn @click="remove">remove</v-btn>
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
  </v-container>
</template>

<script>
import { HTTP } from '../http-common'

export default {
  data: () => ({
    chamber: null,
    errors: [],
    loading: false,
    valid: true,
    nameRules: [(v) => !!v || 'Name is required'],
    pinRules: [(v) => isNaN(v) || 'Pin must be numeric']
  }),

  props: ['macAddress'],

  created () {
    this.fetch()
  },

  watch: {
    '$route': 'fetch'
  },

  methods: {
    fetch () {
      this.chamber = null
      this.errors = []
      this.loading = true
      HTTP.get('chambers/' + this.macAddress)
        .then(response => {
          this.chamber = response.data
        })
        .catch(e => {
          this.errors.push(e)
        })
        .finally(_ => {
          this.loading = false
        })
    },
    addSchedule () {
      this.chamber.schedule.push({})
    },
    save () {
      if (this.$refs.form.validate()) {
        HTTP.post('chambers', this.chamber)
          .then(response => {
            this.$router.push({ name: 'chambers' })
          })
          .catch(e => {
            this.errors.push(e)
          })
      }
    },
    remove () {
      HTTP.delete('chambers/' + this.chamber.macAddress)
        .then(response => {
          this.$router.push({ name: 'chambers' })
        })
        .catch(e => {
          this.errors.push(e)
        })
    }
  }
}
</script>