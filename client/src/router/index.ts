import { createRouter, createWebHistory } from 'vue-router'
import Home from '../pages/Home.vue'
import TickersListPage from '../features/tickers/pages/TickersListPage.vue'
import TickerHistoricalPricePage from '../features/tickers/pages/TickerHistoricalPricePage.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/tickers',
    name: 'Tickers',
    component: TickersListPage
  },
  {
    path: '/ticker/:id',
    name: 'TickerHistoric',
    component: TickerHistoricalPricePage,
    props: true
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
