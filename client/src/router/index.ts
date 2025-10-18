import { createRouter, createWebHistory } from 'vue-router'
import Home from '../pages/Home.vue'
import TickersPage from '../features/tickers/pages/TickersPage.vue'
import TickerHistoricPage from '../features/tickers/pages/TickerHistoricPage.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/tickers',
    name: 'Tickers',
    component: TickersPage
  },
  {
    path: '/ticker/:id',
    name: 'TickerHistoric',
    component: TickerHistoricPage,
    props: true
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
