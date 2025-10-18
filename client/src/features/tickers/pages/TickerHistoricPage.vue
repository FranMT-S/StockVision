<template>
  <div class="ticker-historic-page">
    <v-container fluid class="pa-4">
      <!-- Loading State -->
      <div v-if="loading" class="d-flex justify-center align-center" style="height: 60vh;">
        <v-progress-circular indeterminate size="64" color="primary" />
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center pa-8">
        <v-icon size="64" color="red-lighten-2" class="mb-4">
          mdi-alert-circle-outline
        </v-icon>
        <h3 class="text-h6 text-grey-darken-2 mb-2">Failed to load ticker data</h3>
        <p class="text-body-2 text-grey-darken-1 mb-4">{{ error }}</p>
        <v-btn color="primary" variant="outlined" @click="fetchTickerData">
          <v-icon left>mdi-refresh</v-icon>
          Try Again
        </v-btn>
      </div>

      <!-- Main Content -->
      <div v-else-if="tickerData" class="ticker-content">
        <!-- Stock Header -->
        <StockHeader 
          :ticker-data="tickerData" 
          class="mb-6"
        />

        <v-row>
          <!-- Left Column - Chart and Recommendations -->
          <v-col cols="12" lg="8">
            <!-- Stock Chart -->
            <StockChart 
              :historical-data="historicalData"
              :ticker="tickerId"
              class="mb-6"
            />

            <!-- Recommendations Table -->
            <RecommendationsTable 
              :recommendations="recommendations"
              class="mb-6"
            />
          </v-col>

          <!-- Right Column - Metrics and Company Info -->
          <v-col cols="12" lg="4">
            <!-- Company Metrics -->
            <CompanyMetrics 
              :company-data="tickerData.companyData"
              class="mb-6"
            />

            <!-- Company Information -->
            <CompanyInfo 
              :company-data="tickerData.companyData"
            />
          </v-col>
        </v-row>
      </div>
    </v-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import type { RecomendationResponse } from '@/shared/models/recomendations'
import type { HistoricalPrice } from '@/shared/models/recomendations'
import { mockRecomendationResponses } from '@/shared/models/__mock_recomendation__'
import { mockAAPLHistorical } from '@/shared/models/__mock_historical__'

// Components
import StockHeader from '../components/StockHeader.vue'
import StockChart from '../components/StockChart.vue'
import CompanyMetrics from '../components/CompanyMetrics.vue'
import RecommendationsTable from '../components/RecommendationsTable.vue'
import CompanyInfo from '../components/CompanyInfo.vue'

const route = useRoute()

// State
const loading = ref(true)
const error = ref<string | null>(null)
const tickerData = ref<RecomendationResponse | null>(null)

// Computed
const tickerId = computed(() => route.params.id as string)

const historicalData = computed(() => {
  return mockAAPLHistorical
})

const recommendations = computed(() => {
  // Mock multiple recommendations for AAPL
  return [
    tickerData.value?.recommendation,
    {
      id: 2,
      ticker_id: "AAPL",
      target_from: "170",
      target_to: "185",
      action: "Hold",
      rating_from: "Buy",
      rating_to: "Hold",
      time: "2024-05-15T14:30:00Z",
      brokerage: "Goldman Sachs"
    },
    {
      id: 3,
      ticker_id: "AAPL",
      target_from: "180",
      target_to: "200",
      action: "Buy",
      rating_from: "Hold",
      rating_to: "Buy",
      time: "2024-05-01T09:15:00Z",
      brokerage: "JP Morgan"
    }
  ].filter(Boolean)
})

// Methods
const fetchTickerData = async () => {
  loading.value = true
  error.value = null

  try {
    // Simulate API call delay
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // Find ticker data by ID
    const foundTicker = mockRecomendationResponses.find(
      item => item.ticker.ticker === tickerId.value
    )
    
    if (!foundTicker) {
      throw new Error(`Ticker ${tickerId.value} not found`)
    }
    
    tickerData.value = foundTicker
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to fetch ticker data'
    console.error('Error fetching ticker data:', err)
  } finally {
    loading.value = false
  }
}

// Lifecycle
onMounted(() => {
  fetchTickerData()
})
</script>

<style scoped>
.ticker-historic-page {
  min-height: 100vh;
  background-color: #fafafa;
}

.ticker-content {
  max-width: 1400px;
  margin: 0 auto;
}

/* Responsive adjustments */
@media (max-width: 960px) {
  .ticker-content {
    padding: 0 8px;
  }
}
</style>
