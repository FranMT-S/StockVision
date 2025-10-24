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
        <v-row class="mb-0 tw-relative">
          <v-col cols="12" md="4" lg="3" class="tw-bg-[#ffffff] tw-shadow-sm tw-border-l-2 tw-rounded-lg dark:tw-bg-[#1e1e1e] dark:tw-text-white md:tw-sticky md:tw-top-[84px] md:tw-h-[calc(100vh_-_100px)]">
          
            <!-- Company Information -->
            <CompanyInfo 
              :company-data="tickerData.companyData"
            />
            <!-- Company Metrics -->
            <CompanyMetrics 
              :company-data="tickerData.companyData"
              class="mb-6"
            />
          </v-col>
          <v-col cols="12" md="4" lg="6">
            <!-- Stock Chart -->
            <StockChart 
            
              :historical-data="tickerData.historicalPrices"
              :ticker="tickerId"
              class="mb-6"
            />

            <!-- Recommendations Table -->
            <RecommendationsSection 
              :recommendations="recommendations"
              class="mb-6"
            />
          </v-col>
          <v-col cols="12" md="4" lg="3" class="section-ticker-extra tw-bg-[#ffffff] tw-shadow-sm tw-border-r-2 tw-rounded-lg dark:tw-bg-[#1e1e1e] dark:tw-text-white md:tw-sticky md:tw-top-[84px] md:tw-h-[calc(100vh_-_100px)]">
            <!-- Company News -->
            <CompanyNews 
              :company-news="tickerData.companyNews"
              class="mb-6"
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
import type { CompanyOverview, Recommendation } from '@/shared/models/recomendations'

// Components
import StockChart from '../components/StockChart.vue'
import CompanyMetrics from '../components/CompanyMetrics.vue'
import RecommendationsSection from '../components/RecommendationsTable.vue'
import CompanyInfo from '../components/CompanyInfo.vue'
import CompanyNews from '../components/CompanyNews.vue'
import { useTickersStore } from '../store/tickersStore'
import { storeToRefs } from 'pinia'


const route = useRoute()
const {fetchCompanyOverView} = useTickersStore()
const {companyOverview, error} = storeToRefs(useTickersStore())

// State
const loading = ref(true)
const tickerData = ref<CompanyOverview | null>(null)

// Computed
const tickerId = computed(() => route.params.id as string)

const recommendations = computed<Recommendation[]>(() => {
  if(!companyOverview.value){
    return []
  }

  return companyOverview.value.recommendations
})
    
  
  

// fetch company overview fill the data in the store
// exception safe
const fetchTickerData = async () => {
  loading.value = true
  error.value = null

  await fetchCompanyOverView(tickerId.value)
  
  if(error.value){
    return;
  }

  if (!companyOverview.value) {
    error.value = `Ticker ${tickerId.value} not found`
    return
  }

  tickerData.value = {
    companyData: companyOverview.value?.companyData,
    recommendations: companyOverview.value?.recommendations,
    companyNews: companyOverview.value?.companyNews,
    historicalPrices: companyOverview.value?.historicalPrices,
    advice: companyOverview.value?.advice
  }

  console.log(tickerData.value)

  loading.value = false
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
