<template>
  <v-app dark>

    <v-navigation-drawer persistent v-model="drawer" fixed clipped app enable-resize-watcher>
      <v-list>
        <v-list-tile avatar v-for="(item,i) in items" :key="i" :to="{ name: item.route}">
          <v-list-tile-avatar>
            <v-icon class="amber--text" v-html="item.icon"></v-icon>
          </v-list-tile-avatar>
          <v-list-tile-content>
            <v-list-tile-title class="amber--text" v-text="item.title"></v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile v-if="!authenticated" @click="login()">
          <v-list-tile-avatar>
            <v-icon class="amber--text">power_settings_new</v-icon>
          </v-list-tile-avatar>
          <v-list-tile-content>
            <v-list-tile-title class="amber--text">Login</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile v-if="authenticated" @click="logout()">
          <v-list-tile-avatar>
            <v-icon class="amber--text">power_settings_new</v-icon>
          </v-list-tile-avatar>
          <v-list-tile-content>
            <v-list-tile-title class="amber--text">Logout</v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
    </v-navigation-drawer>

    <v-toolbar color="amber" fixed clipped-left app>
      <v-toolbar-side-icon class="black--text" @click.stop="drawer = !drawer"></v-toolbar-side-icon>
      <v-toolbar-title class="black--text" v-text="title"></v-toolbar-title>
    </v-toolbar>
    
    <v-content>
      <v-container>
        <router-view 
          :auth="auth" 
          :authenticated="authenticated">
        </router-view>
      </v-container>
    </v-content>
    
  </v-app>
</template>

<script>

import AuthService from './auth/AuthService'
const auth = new AuthService()
const { login, logout, authenticated, authNotifier } = auth

export default {
  data () {
    authNotifier.on('authChange', authState => {
      this.authenticated = authState.authenticated
    })
    return {
      auth,
      authenticated,
      drawer: null,
      items: [
        {
          icon: 'bubble_chart',
          title: 'Fermentations',
          route: 'fermentations'
        },
        { icon: 'local_drink', title: 'Beers', route: 'beers' },
        {
          icon: 'devices',
          title: 'Chambers',
          route: 'chambers'
        }
      ],
      title: 'Zymurgauge'
    }
  },
  methods: {
    login,
    logout
  }
}
</script>