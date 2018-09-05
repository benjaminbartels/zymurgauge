import Vue from 'vue'
import VueRouter from 'vue-router'
import Beer from '@/components/Beer.vue'
import Beers from '@/components/Beers.vue'
import Fermentation from '@/components/Fermentation.vue'
import Fermentations from '@/components/Fermentations.vue'
import Chamber from '@/components/Chamber.vue'
import Chambers from '@/components/Chambers.vue'
import Auth from '@/components/Auth.vue'
import Login from '@/components/Login.vue'

const routes = [
  { path: '/auth', name: 'auth', component: Auth },
  { path: '/login', name: 'login', component: Login },
  { path: '/fermentations', name: 'fermentations', component: Fermentations },
  { path: '/fermentations/:id/edit', name: 'editFermentation', component: Fermentation, props: true },
  { path: '/fermentations/create', name: 'createFermentation', component: Fermentation, props: { create: true } },
  { path: '/beers', name: 'beers', component: Beers },
  { path: '/beers/:id/edit', name: 'editBeer', component: Beer, props: true },
  { path: '/beers/create', name: 'createBeer', component: Beer, props: { create: true } },
  { path: '/chambers', name: 'chambers', component: Chambers },
  { path: '/chambers/:macAddress/edit', name: 'editChamber', component: Chamber, props: true },
  { path: '*', redirect: { name: 'chambers' } }
]

const router = new VueRouter({
  mode: 'history',
  routes: routes
 // root: './fermentations' // ToDo: is this correct?
})

Vue.use(VueRouter)

export default router
