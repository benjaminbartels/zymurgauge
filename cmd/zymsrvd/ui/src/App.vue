<template>
  <v-app dark>
    <v-navigation-drawer persistent v-model="drawer" :clipped="true" enable-resize-watcher>
      <v-list>
        <v-list-tile avatar v-for="(item,i) in items" :key="i" :to="{ name: item.route}">
          <v-list-tile-avatar>
            <v-icon class="amber--text" v-html="item.icon"></v-icon>
          </v-list-tile-avatar>
          <v-list-tile-content>
            <v-list-tile-title class="amber--text" v-text="item.title"></v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
    </v-navigation-drawer>
    <v-toolbar class="amber">
      <v-toolbar-side-icon class="black--text" @click.stop="drawer = !drawer"></v-toolbar-side-icon>
      <v-toolbar-title class="black--text" v-text="title"></v-toolbar-title>
    </v-toolbar>
    <main>
      <v-fade-transition>
      <router-view>
      </router-view>
      </v-fade-transition>
    </main>
  </v-app>
</template>

<script>
export default {
  data () {
    return {
      drawer: true,
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
  }
}
</script>