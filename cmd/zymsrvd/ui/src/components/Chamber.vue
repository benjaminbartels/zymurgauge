<template>
  <v-container>
    <div v-if="chamber">
      <v-form v-model="valid">
        <v-text-field label="Name" v-model="chamber.name" required></v-text-field>
        <v-text-field label="Mac Address" v-model="chamber.macAddress" disabled></v-text-field>
        <v-btn @click="submit">submit</v-btn>
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
    valid: false,
    nameRules: [(v) => !!v || 'Name is required']
  }),

  props: ['macAddress'],

  created () {
    if (!this.create) {
      this.fetch()
    } else {
      this.chamber = {}
    }
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
    submit () {
      HTTP.post('chambers', this.chamber)
        .then(response => {
          this.$router.push({ name: 'editChamber', params: { macAddress: this.chamber.macAddress } })
        })
        .catch(e => {
          this.errors.push(e)
        })
    }
  }
}

</script>
