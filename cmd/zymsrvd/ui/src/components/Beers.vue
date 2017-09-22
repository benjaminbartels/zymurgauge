<template>
  <v-container>
    <h3>Beers</h3>
    <v-container grid-list-md>
      <v-layout row wrap>
        <v-flex xs12 sm5 md3 v-for="beer of beers" :key="beer.id">
          <v-card class="grid">
            <router-link :to="{ name: 'editBeer', params: { id: beer.id }} ">
              <v-card-title primary-title>
                {{beer.name}}
              </v-card-title>
              <v-card-text>{{beer.style}}</v-card-text>
            </router-link>
          </v-card>

        </v-flex>
      </v-layout>
    </v-container>
    <v-btn :to="{ name: 'createBeer', params: { create: true }} ">add</v-btn>
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
    beers: [],
    errors: []
  }),

  created () {
    HTTP.get('beers')
      .then(response => {
        this.beers = response.data
      })
      .catch(e => {
        this.errors.push(e)
      })
  }
}
</script>
