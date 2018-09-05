<template>
  <v-container>
    <div v-if="beer">
      <v-form v-model="valid" ref="form" lazy-validation>
        <v-text-field label="Name" v-model="beer.name" :rules="nameRules" required></v-text-field>
        <v-select label="Styles" v-model="beer.style" :rules="styleRules" v-bind:items="styles" required></v-select>
        <v-list>
          <v-list-tile>
            <v-list-tile-content>
              Order
            </v-list-tile-content>
            <v-list-tile-content>
              Target
            </v-list-tile-content>
            <v-list-tile-content>
              Duration
            </v-list-tile-content>
            <v-btn absolute fab right top small @click="addSchedule">
              <v-icon>add</v-icon>
            </v-btn>
          </v-list-tile>
          <template v-for="(schedule, i) in beer.schedule">
            <v-list-tile v-bind:key="schedule.order">
              <v-text-field v-model.number="schedule.order"></v-text-field>
              <v-text-field v-model.number="schedule.targetTemp"></v-text-field>
              <v-text-field v-model.number="schedule.duration"></v-text-field>
              <v-btn fab small @click="removeSchedule(i)">
                <v-icon>delete</v-icon>
              </v-btn>
            </v-list-tile>
          </template>
        </v-list>
        <v-btn @click="save" :disabled="!valid">save</v-btn>
        <v-btn @click="remove" :disabled="this.create">remove</v-btn>
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
import router from './../router'

export default {
  data: () => ({
    beer: null,
    errors: [],
    loading: false,
    valid: true,
    nameRules: [(v) => !!v || 'Name is required'],
    styleRules: [(v) => !!v || 'Style is required'],
    styles: ['1A - American Light Lager',
      '1B - American Lager',
      '1C - Cream Ale',
      '1D - American Wheat Beer',
      '2A - International Pale Lager',
      '2B - International Amber Lager',
      '2C - International Dark Lager',
      '3A - Czech Pale Lager',
      '3B - Czech Premium Pale Lager',
      '3C - Czech Amber Lager',
      '3D - Czech Dark Lager',
      '4A - Munich Helles',
      '4B - Festbier',
      '4C - Helles Bock',
      '5A - German Leichtbier',
      '5B - Kölsch',
      '5C - German Helles Exportbier',
      '5D - German Pils',
      '6A - Märzen',
      '6B - Rauchbier',
      '6C - Dunkles Bock',
      '7A - Vienna Lager',
      '7B - Altbier',
      '7C - Pale Kellerbier',
      '7C - Amber Kellerbier',
      '8A - Munich Dunkel',
      '8B - Schwarzbier',
      '9A - Doppelbock',
      '9B - Eisbock',
      '9C - Baltic Porter',
      '10A - Weissbier',
      '10B - Dunkles Weissbier',
      '10C - Weizenbock',
      '11A - Ordinary Bitter',
      '11B - Best Bitter',
      '11C - Strong Bitter',
      '12A - British Golden Ale',
      '12B - Australian Sparkling Ale',
      '12C - English IPA',
      '13A - Dark Mild',
      '13B - British Brown Ale',
      '13C - English Porter',
      '14A - Scottish Light',
      '14B - Scottish Heavy',
      '14C - Scottish Export',
      '15A - Irish Red Ale',
      '15B - Irish Stout',
      '15C - Irish Extra Stout',
      '16A - Sweet Stout',
      '16B - Oatmeal Stout',
      '16C - Tropical Stout',
      '16D - Foreign Extra Stout',
      '17A - British Strong Ale',
      '17B - Old Ale',
      '17C - Wee Heavy',
      '17D - English Barleywine',
      '18A - Blonde Ale',
      '18B - American Pale Ale',
      '19A - American Amber Ale',
      '19B - California Common',
      '19C - American Brown Ale',
      '20A - American Porter',
      '20B - American Stout',
      '20C - Imperial Stout',
      '21A - American IPA',
      '21B - Specialty IPA - Belgian IPA',
      '21B - Specialty IPA - Black IPA',
      '21B - Specialty IPA - Brown IPA',
      '21B - Specialty IPA - Red IPA',
      '21B - Specialty IPA - Rye IPA',
      '21B - Specialty IPA - White IPA',
      '22A - Double IPA',
      '22B - American Strong Ale',
      '22C - American Barleywine',
      '22D - Wheatwine',
      '23A - Berliner Weisse',
      '23B - Flanders Red Ale',
      '23C - Oud Bruin',
      '23D - Lambic',
      '23E - Gueuze',
      '23F - Fruit Lambic',
      '24A - Witbier',
      '24B - Belgian Pale Ale',
      '24C - Bière de Garde',
      '25A - Belgian Blond Ale',
      '25B - Saison',
      '25C - Belgian Golden Strong Ale',
      '26A - Trappist Single',
      '26B - Belgian Dubbel',
      '26C - Belgian Tripel',
      '26D - Belgian Dark Strong Ale',
      '27A - Gose',
      '27A - Kentucky Common',
      '27A - Lichtenhainer',
      '27A - London Brown Ale',
      '27A - Piwo Grodziskie',
      '27A - Pre-Prohibition Lager',
      '27A - Pre-Prohibition Porter',
      '27A - Roggenbier',
      '27A - Sahti',
      '28A - Brett Beer',
      '28B - Mixed-Fermentation Sour Beer',
      '28C - Wild Specialty Beer',
      '29A - Fruit Beer',
      '29B - Fruit and Spice Beer',
      '29C - Specialty Fruit Beer',
      '30A - Spice, Herb, or Vegetable Beer',
      '30B - Autumn Seasonal Beer',
      '30C - Winter Seasonal Beer',
      '31A - Alternative Grain Beer',
      '31B - Alternative Sugar Beer',
      '32A - Classic Style Smoked Beer',
      '32B - Specialty Smoked Beer',
      '33A - Wood-Aged Beer',
      '33B - Specialty Wood-Aged Beer',
      '34A - Clone Beer',
      '34B - Mixed-Style Beer',
      '34C - Experimental Beer']
  }),
  props: ['authenticated', 'id', 'create'],
  created () {
    if (!this.authenticated) {
      router.replace('login')
    }

    if (!this.create) {
      this.fetch()
    } else {
      this.beer = {
        name: '',
        style: ''
      }
    }
  },

  watch: {
    '$route': 'fetch'
  },

  methods: {
    fetch () {
      this.beer = null
      this.errors = []
      this.loading = true
      HTTP.get('beers/' + this.id)
        .then(response => {
          this.beer = response.data
        })
        .catch(e => {
          this.errors.push(e)
        })
        .finally(_ => {
          this.loading = false
        })
    },
    addSchedule () {
      if (this.beer.schedule == null) {
        this.beer.schedule = []
      }
      this.beer.schedule.push({})
    },
    removeSchedule (index) {
      this.beer.schedule.splice(index, 1)
    },
    save () {
      if (this.$refs.form.validate()) {
        HTTP.post('beers', this.beer)
          .then(response => {
            this.$router.push({ name: 'beers' })
          })
          .catch(e => {
            this.errors.push(e)
          })
      }
    },
    remove () {
      HTTP.delete('beers/' + this.beer.id)
        .then(response => {
          this.$router.push({ name: 'beers' })
        })
        .catch(e => {
          this.errors.push(e)
        })
    }
  }
}
</script>
