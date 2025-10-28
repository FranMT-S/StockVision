import { defineStore } from 'pinia'
import { ref, computed, watch, onMounted, onBeforeMount } from 'vue'
import type { CompanyOverview, HistoricalPrice, TickerListResponse } from '@/shared/models/recomendations'
import { customFetch } from '@/shared/services/customFetch'
import { API_CONFIG } from '@/shared/constants/api'
import { ErrorType } from '@/shared/models/response'
import { useDebounce } from '@/shared/composables/useDebounce'
import { useRoute } from 'vue-router'

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
  const route = useRoute();
  
  const initialPage = Number.isNaN(Number(route.query.page)) ? 1 : Number(route.query.page)
  const initialQuery = route.query.q == undefined || !route.query.q ? '' : route.query.q?.toString()

  // State
  const stockData = ref<StockData[]>([])
  const loading = ref(true)
  const error = ref<string | null>(null)
  const errorPredictions = ref<string | null>(null)
  const errorHistoricalPrices = ref<string | null>(null)
  const tickers = ref<TickerListResponse[]>([])
  const companyOverview = ref<CompanyOverview | null>(null)
  const companyPredictions = ref<HistoricalPrice[]>([])
  const companyHistoricalPrices = ref<HistoricalPrice[]>([])

  // Pagination state
  const currentPage = ref(initialPage)
  const itemsPerPage = ref(10)
  const totalItems = ref(0)
  const sort = ref<'asc' | 'desc'>('asc')
  const search = ref(initialQuery)
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
    const {debounced, cancel} = useDebounce()
    
    // avoid flickering if the request is too fast
    debounced(() => {
      loading.value = true
    }, 200)
    
    abortIfControllerIsActive(tickerCancellationToken.value)
    tickerCancellationToken.value = new AbortController()

    error.value = null
    const params = {
      q: search.value,
      page: currentPage.value,
      size: itemsPerPage.value,
      sort: sort.value
    }
    const url = API_CONFIG.ENDPOINTS.List(params)
    const response = await customFetch<TickerListResponse[]>(url,{signal: tickerCancellationToken.value.signal})

    if(response.ok){
      tickers.value = response.data
      totalItems.value = response.total || 0
    }else{
      error.value = response.errorType !== ErrorType.ABORT_ERROR ? response.error : null
    }
    
    cancel()
    loading.value = false
  }

  const fetchCompanyOverView = async (ticker: string,from:Date) => {
    const {debounced, cancel} = useDebounce()
    
    // avoid flickering if the request is too fast
    debounced(() => {
      loading.value = true
    }, 200)
    
    abortIfControllerIsActive(companyOverviewCancellationToken.value)
    companyOverviewCancellationToken.value = new AbortController()  
    const url = API_CONFIG.ENDPOINTS.Overview(ticker,from)
    
    error.value = null
    const response = await customFetch<CompanyOverview>(url,{signal: companyOverviewCancellationToken.value.signal})

    if(response.ok){
      companyOverview.value = response.data
    }else{
      error.value = response.errorType !== ErrorType.ABORT_ERROR ? response.error : null 
    }
    
    cancel()
    loading.value = false
  }

  /**
   * fetch company historical prices and update the state of companyHistoricalPrices and errorHistoricalPrices
   * abort the request if it is already active
   * this function not use loading state
   */
  const fetchCompanyHistoricalPrices = async (ticker: string,from?:Date,abortController: AbortController | null = null) => {
    abortIfControllerIsActive(abortController)
    const url = API_CONFIG.ENDPOINTS.HistoricalPrices(ticker,from)
    
    error.value = null
    const response = await customFetch<HistoricalPrice[]>(url,{signal: abortController?.signal})

    if(response.ok){
      companyHistoricalPrices.value = response.data
    }else{
      error.value = response.errorType !== ErrorType.ABORT_ERROR ? response.error : null 
    }
  }

  /**
   * fetch company predictions and update the state of companyPredictions and errorPredictions
   * abort the request if it is already active
   * this function not use loading state
   */
  const fetchCompanyPredictions = async (ticker: string) => {
    abortIfControllerIsActive(companyPredictionsCancellationToken.value)
    companyPredictionsCancellationToken.value = new AbortController()  
    const url = API_CONFIG.ENDPOINTS.Predictions(ticker)
    
    const errrosList = [
      "Vision tried to see the future, but the lens fogged up. Try again soon",
      "Vision can’t see right now — try again later",
      "Oops… Vision lost its glasses. Try again later"
    ]

    errorPredictions.value = null
    const response = await customFetch<HistoricalPrice[]>(url,{signal: companyPredictionsCancellationToken.value.signal})
    if(response.ok){
      companyPredictions.value = response.data
    }else{
      errorPredictions.value = errrosList[Math.floor(Math.random() * errrosList.length)] 
      console.error(response.error)
    }  
  }

  onBeforeMount(() => {
    search.value =  route.query.q?.toString() || ''
    currentPage.value = Number.isNaN(Number(route.query.page)) ? 1 : Number(route.query.page)
  })

  return {
    // State
    stockData,
    companyOverview,
    loading,
    error,
    errorPredictions,
    errorHistoricalPrices,
    totalPages,
    companyHistoricalPrices,
    
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
    fetchTickers,
    fetchCompanyHistoricalPrices
  }
})
