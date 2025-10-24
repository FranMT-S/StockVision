import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import type { CompanyOverview, HistoricalPrice, TickerListResponse } from '@/shared/models/recomendations'
import { customFetch } from '@/shared/services/customFetch'
import { API_CONFIG } from '@/shared/constants/api'
import { ErrorType } from '@/shared/models/response'
import { debounce } from 'vuetify/lib/util/helpers.mjs'
import { useDebounce } from '@/shared/composables/useDebounce'

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
  const loading = ref(true)
  const error = ref<string | null>(null)
  const errorPredictions = ref<string | null>(null)
  const tickers = ref<TickerListResponse[]>([])
  const companyOverview = ref<CompanyOverview | null>(null)
  const companyPredictions = ref<HistoricalPrice[] | null>(null)

  // Pagination state
  const currentPage = ref(1)
  const itemsPerPage = ref(10)
  const totalItems = ref(0)
  const sort = ref<'asc' | 'desc'>('asc')
  const search = ref('')
  const tickerCancellationToken = ref<AbortController | null>(null)
  const companyOverviewCancellationToken = ref<AbortController | null>(null)
  const companyPredictionsCancellationToken = ref<AbortController | null>(null)
  const totalPages = computed(() => Math.ceil(totalItems.value / itemsPerPage.value))

  function abortIfControllerIsActive(controller: AbortController | null){
    if(controller && !controller.signal.aborted){
      controller.abort()
    }
  }

  watch([currentPage, sort, search], () => {
    fetchTickers()
  })

  watch(search, (newSearch) => {
    if(!newSearch)
      search.value = ''

    currentPage.value = 1
  })

  watch(itemsPerPage, () => {
    currentPage.value = 1
  })


  const fetchTickers = async () => {
    let showLoader = true; 
    const {debounced, cancel} = useDebounce()
    
    // avoid flickering if the request is too fast
    debounced(() => {
      loading.value = true
    }, 200)
    
    abortIfControllerIsActive(tickerCancellationToken.value)
    tickerCancellationToken.value = new AbortController()

    error.value = null
    const url = API_CONFIG.BASE_URL + API_CONFIG.ENDPOINTS.List(search.value,currentPage.value, itemsPerPage.value, sort.value)
    const response = await customFetch<TickerListResponse[]>(url,{signal: tickerCancellationToken.value.signal})

    if(response.ok){
      tickers.value = response.data
      totalItems.value = response.total || 0
    }else{
      error.value = response.errorType !== ErrorType.ABORT_ERROR ? response.error : null
    }
    
    cancel()
    if(showLoader){
      loading.value = false
    }
  }

  const fetchCompanyOverView = async (ticker: string) => {
    let showLoader = true; 
    const {debounced, cancel} = useDebounce()
    
    // avoid flickering if the request is too fast
    debounced(() => {
      loading.value = true
      showLoader = true
    }, 200)
    
    abortIfControllerIsActive(companyOverviewCancellationToken.value)
    companyOverviewCancellationToken.value = new AbortController()  
    const url = API_CONFIG.BASE_URL + API_CONFIG.ENDPOINTS.Overview(ticker)
    
    error.value = null
    const response = await customFetch<CompanyOverview>(url,{signal: companyOverviewCancellationToken.value.signal})

    if(response.ok){
      companyOverview.value = response.data
    }else{
      error.value = response.errorType !== ErrorType.ABORT_ERROR ? response.error : null 
    }
    
    cancel()
    if(showLoader){
      loading.value = false
    }
  }

  const fetchCompanyPredictions = async (ticker: string) => {
    abortIfControllerIsActive(companyPredictionsCancellationToken.value)
    companyPredictionsCancellationToken.value = new AbortController()  
    const url = API_CONFIG.BASE_URL + API_CONFIG.ENDPOINTS.Predictions(ticker)
    
    errorPredictions.value = null
    const response = await customFetch<HistoricalPrice[]>(url,{signal: companyPredictionsCancellationToken.value.signal})

    if(response.ok){
      companyPredictions.value = response.data
    }else{
      errorPredictions.value = response.errorType !== ErrorType.ABORT_ERROR ? response.error : null 
    }  
  }

  return {
    // State
    stockData,
    companyOverview,
    loading,
    error,
    errorPredictions,
    totalPages,
    
    // Pagination state
    currentPage,
    itemsPerPage,
    totalItems,
    sort,
    search,
    tickers,
    companyPredictions,
    fetchCompanyOverView,
    fetchCompanyPredictions,
    fetchTickers
  }
})
