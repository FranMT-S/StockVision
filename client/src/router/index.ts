import { createRouter, createWebHistory } from 'vue-router'
import TickersListPage from '../features/tickers/pages/TickersListPage.vue'
import TickerHistoricalPricePage from '../features/tickers/pages/TickerHistoricalPricePage.vue'
import { RouterNames } from './names'

const routes = [
  {
    path: '/',
    name: RouterNames.Tickers,
    component: TickersListPage
  },
  {
    path: '/ticker/:id',
    name: RouterNames.TickerHistoric,
    component: TickerHistoricalPricePage,
    props: true
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
