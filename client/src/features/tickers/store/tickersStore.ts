import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import type { CompanyOverview, TickerListResponse } from '@/shared/models/recomendations'
import { customFetch } from '@/shared/services/customFetch'
import { API_CONFIG } from '@/shared/constants/api'

interface StockData {
  ticker: string
  companyName: string
  price: number
  url?: string
  change: number
  changePercentage: number
  sentiment: string
  lastRatingDate: string
}

export const useTickersStore = defineStore('tickers', () => {
  // State
  const stockData = ref<StockData[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const tickers = ref<TickerListResponse[]>([])
  const companyOverview = ref<CompanyOverview | null>(null)

  // Pagination state
  const currentPage = ref(1)
  const itemsPerPage = ref(10)
  const totalItems = ref(0)
  const sort = ref<'asc' | 'desc'>('asc')
  const search = ref('')
  const tickerCancellationToken = ref<AbortController | null>(null)
  const companyOverviewCancellationToken = ref<AbortController | null>(null)

  function abortIfControllerIsActive(controller: AbortController | null){
    if(controller && !controller.signal.aborted){
      controller.abort()
    }
  }

  const fetchTickers = async () => {
    loading.value = true
    error.value = null
    
    abortIfControllerIsActive(tickerCancellationToken.value)
    tickerCancellationToken.value = new AbortController()

    const url = API_CONFIG.BASE_URL + API_CONFIG.ENDPOINTS.List(search.value,currentPage.value, itemsPerPage.value, sort.value)

    const response = await customFetch<TickerListResponse[]>(url,{signal: tickerCancellationToken.value.signal})

    if(response.ok){
      tickers.value = response.data
      totalItems.value = response.total || 0
    }else{
      error.value = response.error
    }

    loading.value = false
  }

  watch([currentPage, sort], () => {
    fetchTickers()
  })

  watch(search, () => {
    currentPage.value = 1
    
    fetchTickers()
  })

  watch(itemsPerPage, () => {
    currentPage.value = 1
    fetchTickers()
  })

  const fetchCompanyOverView = async (ticker: string) => {
    loading.value = true
    error.value = null
    
    abortIfControllerIsActive(companyOverviewCancellationToken.value)
    companyOverviewCancellationToken.value = new AbortController()

    const url = API_CONFIG.BASE_URL + API_CONFIG.ENDPOINTS.Overview(ticker)

    const response = await customFetch<CompanyOverview>(url,{signal: companyOverviewCancellationToken.value.signal})

    if(response.ok){
      companyOverview.value = response.data
    }else{
      error.value = response.error
    }

    loading.value = false
  }

  return {
    // State
    stockData,
    companyOverview,
    loading,
    error,
    
    // Pagination state
    currentPage,
    itemsPerPage,
    totalItems,
    sort,
    search,
    tickers,
    fetchCompanyOverView,
    fetchTickers
  }
})
