<style scoped>
a {
  color: #FFB90D;
  text-decoration: none;
}
</style>

<template>
  <v-container>
    <h3>Fermentations</h3>
    <v-container grid-list-md>
      <v-layout row wrap>
        <v-flex  xs12 sm5 md3 v-for="fermentation of fermentations" :key="fermentation.id">
          <v-card class="grid">
            <router-link :to="{ name: 'editFermentation', params: { id: fermentation.id }} ">
              <v-card-title primary-title class="headline">
                {{fermentation.chamber.name}}
              </v-card-title>
              <v-card-text>{{fermentation.beer.name}}</v-card-text>
            </router-link>
          </v-card>
        </v-flex>
      </v-layout>
    </v-container>
    <v-btn :to="{ name: 'createFermentation', params: { create: true }} ">add</v-btn>
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
    fermentations: [],
    errors: []
  }),

  created () {
    HTTP.get('fermentations')
      .then(response => {
        this.fermentations = response.data
      })
      .catch(e => {
        this.errors.push(e)
      })
  }
}
</script>
