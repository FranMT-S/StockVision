<template>
  <div class="ticker-historic-page">
    <v-container fluid class="pa-0 lg-pa-4">
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
        <v-row class="mb-0 tw-relative tw-justify-end tw-px-[4px] lg:tw-px-0">
          <v-col :cols="columns.side"
           class="tw-bg-[#ffffff] lg:tw-shadow-sm  md:tw-p-4 tw-border-l-2 tw-rounded-lg dark:tw-bg-[#1e1e1e] dark:tw-text-white lg:tw-sticky md:tw-top-[84px] px-0 !tw-pt-4 !tw-ps-[20px]"
            :class="{ 'tw-h-[calc(100vh_-_100px)]': isLargeDesktop }"
          >
            <!-- Company Information -->
            <CompanyInfo 
              
              :company-data="tickerData.companyData"
            />
            <!-- Company Metrics -->
            <CompanyMetrics 
              :company-data="tickerData.companyData"
              class="md:mb-6 lg:tw-mx-0 tw-mx-3 lg:tw-mb-6"
            />  

            <AdviceBanner 
              v-if="isDesktop"
              
              :advice="tickerData.advice"
              tittle="Friendly Advice"
            />
          </v-col>
          <v-col :cols="columns.main">
            <!-- Stock Chart -->
            <StockChart           
              :historical-data="tickerData.historicalPrices"
              :ticker="tickerId"
              class="mb-6 lg:tw-mx-0 tw-mx-3"
              @update:timeframe="interval = $event"
              @update:predict="fetchPredictions()"
              :predictNextWeek="companyPredictions"
              :predictError="errorPredictions"
              :isPredictLoading="isPredictLoading"
              :isTouchDevice="isTouchable"
            />

            <AdviceBanner 
              v-if="!isDesktop"
              :advice="tickerData.advice"
              tittle="Friendly Advice"
              class="mb-6 lg:tw-mx-0 tw-mx-3"
            />

            <!-- Recommendations Table -->
            <RecommendationsSection 
              :recommendations="recommendations"
              class="mb-6 lg:tw-mx-0 tw-mx-3"
            />
          </v-col>
          <v-col :cols="columns.new" 
          class="section-ticker-extra tw-bg-[#ffffff] tw-shadow-sm tw-border-r-2 tw-rounded-lg dark:tw-bg-[#1e1e1e] dark:tw-text-white lg-h-full lg:tw-sticky xl:tw-top-[84px] xl:tw-h-[calc(100vh_-_100px)]">
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
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import type { CompanyOverview,  Recommendation } from '@/shared/models/recomendations'

// Components
import StockChart from '../components/StockChart.vue'
import CompanyMetrics from '../components/CompanyMetrics.vue'
import RecommendationsSection from '../components/RecommendationsTable.vue'
import CompanyInfo from '../components/CompanyInfo.vue'
import CompanyNews from '../components/CompanyNews.vue'
import AdviceBanner from '../components/AdviceBanner.vue'
import { useTickersStore } from '../store/tickersStore'
import { storeToRefs } from 'pinia'
import { Timeframe } from '@/shared/enums/timeFrame'
import { useBreakpoints } from '@/shared/composables/useBreakpoints';


const route = useRoute()
const {fetchCompanyOverView, fetchCompanyHistoricalPrices, fetchCompanyPredictions} = useTickersStore()
const {companyOverview, error, companyHistoricalPrices,companyPredictions,errorPredictions} = storeToRefs(useTickersStore())
const {isMobile,isTablet,isTouchable,isDesktop,isLargeDesktop,screenWidth} = useBreakpoints()

// State
const loading = ref(true)
const tickerData = ref<CompanyOverview | null>(null)
const interval = ref(Timeframe['1M'])
const isPredictLoading = ref(false)

// Computed
const tickerId = computed(() => route.params.id as string)
const recommendations = computed<Recommendation[]>(() => {
  if(!companyOverview.value){
    return []
  }
  return companyOverview.value.recommendations
})

const columns = computed<{
  side: number,
  main: number,
  new:number
}>(() =>{
  if(screenWidth.value < 1024){
    return { side: 12,
      main: 12,
      new: 12
    }
  }

  if(screenWidth.value < 1280){ 
    return {
      side: 3,
      main: 9,
      new: 9
    }
  }

 
  return {
    side: 3,
    main: 6,
    new: 3
  }
  

}) 

// fetch company overview fill the data in the store
// exception safe
const fetchTickerData = async () => {
  error.value = null
  const from = new Date()
  from.setDate(from.getDate() - interval.value)
  await fetchCompanyOverView(tickerId.value,from)
  
  if(error.value){
    loading.value = false
    return;
  }

  if (!companyOverview.value ) {
    error.value = `Ticker ${tickerId.value} not found`
    loading.value = false
    return
  }

  tickerData.value = {
    companyData: companyOverview.value?.companyData,
    recommendations: companyOverview.value?.recommendations,
    companyNews: companyOverview.value?.companyNews,
    historicalPrices: companyOverview.value?.historicalPrices,
    advice: companyOverview.value?.advice
  }
}

const fetchPredictions = async () => {
  isPredictLoading.value = true
  await fetchCompanyPredictions(tickerId.value)
  isPredictLoading.value = false
}

watch(interval, async () => {
    let now = undefined;
    if(interval.value !== Timeframe['All']){
      now = new Date();
      now.setDate(now.getDate() - interval.value)
    }
    
    await fetchCompanyHistoricalPrices(tickerId.value, now)
})

watch(companyHistoricalPrices, () => {
  if(!tickerData.value){
    return;
  }

  tickerData.value = {
    ...tickerData.value,
    historicalPrices: companyHistoricalPrices.value
  }
})

// Lifecycle
onMounted(async () => {
  loading.value = true
  const now = new Date();
  now.setDate(now.getDate() - interval.value) 
  await fetchTickerData()
  companyPredictions.value = []

  if(error.value){
    loading.value = false
    return;
  }
  
  loading.value = false
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
