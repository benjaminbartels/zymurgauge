<template>
    <v-app>
        <v-navigation-drawer persistent v-model="drawer" :clipped="true" enable-resize-watcher>
            <v-list>
                <v-list-tile avatar v-for="(item,i) in items" :key="i" :to="{ name: item.route}">
                    <v-list-tile-avatar>
                        <v-icon v-html="item.icon"></v-icon>
                    </v-list-tile-avatar>
                    <v-list-tile-content>
                        <v-list-tile-title v-text="item.title"></v-list-tile-title>
                    </v-list-tile-content>
                </v-list-tile>
            </v-list>
        </v-navigation-drawer>
        <v-toolbar>
            <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
            <v-toolbar-title v-text="title"></v-toolbar-title>
        </v-toolbar>
        <main>
            <v-fade-transition>
                <router-view></router-view>
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
                { icon: 'bubble_chart', title: 'Fermentations', route: 'fermentations' },
                { icon: 'local_drink', title: 'Beers', route: 'beers' },
                { icon: 'devices', title: 'Chambers', route: 'chambers' }
            ],
            title: 'Zymurgauge'
        }
    }
}
</script>

<style lang="stylus">
  @import './stylus/main'
</style>
