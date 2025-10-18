import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { RecomendationResponse } from '@/shared/models/recomendations'

interface StockData {
  ticker: string
  companyName: string
  price: number
  change: number
  changePercentage: number
  sentiment: string
}

export const useTickersStore = defineStore('tickers', () => {
  // State
  const stockData = ref<StockData[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  
  // Pagination state
  const currentPage = ref(1)
  const itemsPerPage = ref(10)
  const totalItems = ref(0)

  // Mock data - in real app this would come from API
  const mockStockData: StockData[] = [
    {
      ticker: 'AAPL',
      companyName: 'Apple Inc.',
      price: 178.25,
      change: 2.34,
      changePercentage: 1.33,
      sentiment: 'positive'
    },
    {
      ticker: 'MSFT',
      companyName: 'Microsoft Corp.',
      price: 378.91,
      change: -1.23,
      changePercentage: -0.32,
      sentiment: 'neutral'
    },
    {
      ticker: 'GOOGL',
      companyName: 'Alphabet Inc.',
      price: 142.65,
      change: 3.45,
      changePercentage: 2.48,
      sentiment: 'positive'
    },
    {
      ticker: 'AMZN',
      companyName: 'Amazon.com Inc.',
      price: 152.38,
      change: -4.12,
      changePercentage: -2.63,
      sentiment: 'negative'
    },
    {
      ticker: 'TSLA',
      companyName: 'Tesla Inc.',
      price: 248.50,
      change: 5.67,
      changePercentage: 2.33,
      sentiment: 'positive'
    },
    {
      ticker: 'NVDA',
      companyName: 'NVIDIA Corp.',
      price: 875.28,
      change: -12.45,
      changePercentage: -1.40,
      sentiment: 'negative'
    },
    {
      ticker: 'META',
      companyName: 'Meta Platforms Inc.',
      price: 312.45,
      change: 8.32,
      changePercentage: 2.73,
      sentiment: 'positive'
    },
    {
      ticker: 'NFLX',
      companyName: 'Netflix Inc.',
      price: 485.67,
      change: -3.21,
      changePercentage: -0.66,
      sentiment: 'neutral'
    },
    {
      ticker: 'AMD',
      companyName: 'Advanced Micro Devices Inc.',
      price: 142.18,
      change: 5.43,
      changePercentage: 3.97,
      sentiment: 'positive'
    },
    {
      ticker: 'INTC',
      companyName: 'Intel Corporation',
      price: 45.23,
      change: -1.87,
      changePercentage: -3.97,
      sentiment: 'negative'
    },
    {
      ticker: 'CRM',
      companyName: 'Salesforce Inc.',
      price: 267.89,
      change: 4.56,
      changePercentage: 1.73,
      sentiment: 'positive'
    },
    {
      ticker: 'ADBE',
      companyName: 'Adobe Inc.',
      price: 523.12,
      change: -7.89,
      changePercentage: -1.49,
      sentiment: 'negative'
    },
    {
      ticker: 'PYPL',
      companyName: 'PayPal Holdings Inc.',
      price: 68.45,
      change: 2.34,
      changePercentage: 3.54,
      sentiment: 'positive'
    },
    {
      ticker: 'UBER',
      companyName: 'Uber Technologies Inc.',
      price: 78.23,
      change: -1.23,
      changePercentage: -1.55,
      sentiment: 'neutral'
    },
    {
      ticker: 'SPOT',
      companyName: 'Spotify Technology S.A.',
      price: 234.56,
      change: 6.78,
      changePercentage: 2.97,
      sentiment: 'positive'
    },
    {
      ticker: 'ZM',
      companyName: 'Zoom Video Communications Inc.',
      price: 89.34,
      change: -2.45,
      changePercentage: -2.67,
      sentiment: 'negative'
    },
    {
      ticker: 'SHOP',
      companyName: 'Shopify Inc.',
      price: 78.91,
      change: 3.67,
      changePercentage: 4.87,
      sentiment: 'positive'
    }
  ]

  // Actions
  const fetchStockData = async () => {
    loading.value = true
    error.value = null
    
    try {
      // Simulate API call delay
      await new Promise(resolve => setTimeout(resolve, 600))
      
      // In real app, this would be an API call
      // const response = await fetch('/api/tickers')
      // const data = await response.json()
      
      // const number =  Math.random() * 100;
      // console.log(number)
      // if(number > 50)
      //   throw ("test")

      stockData.value = mockStockData
      totalItems.value = mockStockData.length
    } catch (err) {
      error.value = 'Failed to fetch stock data'
      console.error('Error fetching stock data:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchTickers = async () => {
    loading.value = true
    error.value = null
    
    try {
      // Simulate API call delay
      await new Promise(resolve => setTimeout(resolve, 800))
      
      // In real app, this would be an API call
      // const response = await fetch('/api/tickers')
      // const data = await response.json()
      
      // For now, just fetch stock data
      await fetchStockData()
    } catch (err) {
      error.value = 'Failed to fetch tickers'
      console.error('Error fetching tickers:', err)
    } finally {
      loading.value = false
    }
  }

  const searchStocks = async (query: string) => {
    if (!query.trim()) {
      await fetchStockData()
      return
    }

    loading.value = true
    error.value = null
    
    try {
      // Simulate search delay
      await new Promise(resolve => setTimeout(resolve, 500))
      
      // Filter mock data based on query
      const filteredData = mockStockData.filter(stock => 
        stock.ticker.toLowerCase().includes(query.toLowerCase()) ||
        stock.companyName.toLowerCase().includes(query.toLowerCase())
      )
      
      stockData.value = filteredData
    } catch (err) {
      error.value = 'Failed to search stocks'
      console.error('Error searching stocks:', err)
    } finally {
      loading.value = false
    }
  }

  // Pagination actions
  const setCurrentPage = (page: number) => {
    currentPage.value = page
  }

  const setItemsPerPage = (items: number) => {
    itemsPerPage.value = items
    currentPage.value = 1 // Reset to first page when changing items per page
  }

  // Getters
  const totalPages = computed(() => 
    Math.ceil(totalItems.value / itemsPerPage.value)
  )

  const paginatedStockData = computed(() => {
    const start = (currentPage.value - 1) * itemsPerPage.value
    const end = start + itemsPerPage.value
    return stockData.value.slice(start, end)
  })

  const positiveStocks = computed(() => 
    stockData.value.filter(stock => stock.change > 0)
  )

  const negativeStocks = computed(() => 
    stockData.value.filter(stock => stock.change < 0)
  )

  const topPerformers = computed(() => 
    [...stockData.value]
      .sort((a, b) => b.changePercentage - a.changePercentage)
      .slice(0, 3)
  )

  return {
    // State
    stockData,
    loading,
    error,
    
    // Pagination state
    currentPage,
    itemsPerPage,
    totalItems,
    
    // Actions
    fetchStockData,
    fetchTickers,
    searchStocks,
    setCurrentPage,
    setItemsPerPage,
    
    // Getters
    totalPages,
    paginatedStockData,
    positiveStocks,
    negativeStocks,
    topPerformers
  }
})
