import Vue from 'vue'
import VueRouter from 'vue-router'
import Beer from '@/components/Beer.vue'
import Beers from '@/components/Beers.vue'
import Chamber from '@/components/Chamber.vue'
import Chambers from '@/components/Chambers.vue'

const routes = [
  { path: '/beers', name: 'beers', component: Beers },
  { path: '/beers/:id/edit', name: 'editBeer', component: Beer, props: true },
  { path: '/beers/create', name: 'createBeer', component: Beer, props: { create: true } },
  { path: '/chambers', name: 'chambers', component: Chambers },
  { path: '/chambers/:macAddress/edit', name: 'editChamber', component: Chamber, props: true },
  { path: '*', redirect: { name: 'chambers' } }
]

const router = new VueRouter({
  routes: routes,
  root: './fermentations'
})

Vue.use(VueRouter)

export default router
