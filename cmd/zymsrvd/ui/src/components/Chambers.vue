<style scoped>
a {
  color: #ffb90d;
  text-decoration: none;
}
</style>

<template>
  <v-container>
    <h3>Chambers</h3>
    <v-container grid-list-md>
      <v-layout row wrap>
        <v-flex xs12 sm5 md3 v-for="chamber of chambers" :key="chamber.macAddress">
          <v-card class="grid">
            <router-link :to="{ name: 'editChamber', params: { macAddress: chamber.macAddress }} ">
              <v-card-title primary-title class="headline">
                {{chamber.name}}
              </v-card-title>
              <v-card-text>{{chamber.macAddress}}</v-card-text>
            </router-link>
          </v-card>
        </v-flex>
      </v-layout>
    </v-container>
    <ul v-if=" errors && errors.length ">
      <li v-for="error of errors " :key="error.message ">
        {{error.message}}
      </li>
    </ul>
  </v-container>
</template>

<script>
import { HTTP } from '../http-common'

export default {
  data: () => ({
    chambers: [],
    errors: []
  }),

  created () {
    HTTP.get('chambers')
      .then(response => {
        this.chambers = response.data
      })
      .catch(e => {
        this.errors.push(e)
      })
  }
}
</script>
