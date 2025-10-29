<template>
  <div class="tw-mb-[32px]">
    <v-container fluid class="pb-0"> 
      <!-- Header Section -->
      <div class="mb-3">
        <div class="d-flex align-center mb-2">
          <h2 class="text-h5 font-weight-bold text-grey-darken-3 mr-3">
            Company Stock Prices
          </h2>
          <v-progress-circular
            v-if="loading"
            indeterminate
            size="24"
            width="3"
            color="primary"
          />
        </div>
      </div>

      <!-- Stock Table -->
      <v-card elevation="2" class="rounded-lg">
        <v-card-text class="pa-0">
          <!-- Loading State -->
          <StockTableSkeleton style="max-height: calc(100vh - 200px)" v-if="loading" :rows="6" />
          
          <!-- Error State -->
          <ErrorFetchData v-else-if="error" :error="error" class="pa-8 text-center" @click="tickersStore.fetchTickers()"/> 
       
          <!-- Data State -->
          <TickersTable 
            v-else-if="tableRows.length > 0" 
            style="max-height: calc(100vh - 200px)"
            :rows="tableRows"
            @row-click="handleStockClick"
          />

          <!-- Empty State -->
          <div v-else class="pa-8 text-center">
            <v-icon size="64" color="grey-lighten-2" class="mb-4">
              mdi-chart-line
            </v-icon>
            <h3 class="text-h6 text-grey-darken-2 mb-2">No data Found</h3>
            <p class="text-body-2 text-grey-darken-1">
              No data found. Please try refreshing the page.
            </p>
          </div>
        </v-card-text>
      </v-card>

      <!-- Paginator -->
      <Paginator
        v-if="!loading && !error && tableRows.length > 0"
        :total-pages="totalPages"
        :total-items="totalItems"
        :items-per-page="itemsPerPage"
        :visible-pages="visiblePages"
        :initial-page="currentPage"
        sticky
        class="mt-4"
        @update:current-page="handlePageChange"
      />
    </v-container>
  </div>
</template>

<script setup lang="ts">
import { computed, ComputedRef, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useTickersStore } from '../store/tickersStore'
import StockTableSkeleton from '@/shared/components/StockTableSkeleton.vue'
import ErrorFetchData from '@/shared/components/ErrorFetchData.vue'
import TickersTable, {  TickerRow } from '../components/TickersTable.vue'
import Paginator from '@/shared/components/Paginator.vue'
import { useBreakpoints } from '@/shared/composables/useBreakpoints'
import { storeToRefs } from 'pinia'
import { useRoute } from 'vue-router'
import { RouterNames } from '@/router/names'

const router = useRouter()
const route = useRoute();
const tickersStore = useTickersStore()
const {tickers, totalPages, totalItems,loading,error,currentPage,itemsPerPage,search} = storeToRefs(tickersStore)

onMounted(() => {
  currentPage.value = Number.isNaN(Number(route.query.page)) ? 1 : Number(route.query.page);
  search.value = String(route.query.q ?? '');

  if(tickers.value.length === 0)
    tickersStore.fetchTickers()
})

watch(
  () => route.query,
  (newQuery) => {
    tickersStore.currentPage = Number.isNaN(Number(newQuery.page)) ? 1 : Number(newQuery.page);
    tickersStore.search = String(newQuery.q) || '';
  },
  { deep: true }
);

watch(currentPage, () => {
  router.push({
    name: RouterNames.Tickers,
    query: {
      q: search.value == 'undefined' || !search.value ? undefined : search.value,
      page: currentPage.value == 1 ? undefined : currentPage.value,
    }
  })
})

const tableRows: ComputedRef<TickerRow[]> = computed(() => {
  let rows: TickerRow[] = []
  tickers.value.map(({companyData, ticker,advice}) => {
    rows.push({
      ticker: ticker.id,
      companyName: companyData.companyName,
      price: companyData.price,
      url: companyData.image,
      change: companyData.change,
      changePercentage: companyData.changePercentage,
      sentiment: ticker.sentiment,
      lastRatingDate: ticker.recommendations?.[0].time || 'Not available',
      advice: advice.split('.')[0] || 'Not available'
    })  
  })
  
  return rows
})


const { isMobile, isTablet, isDesktop } = useBreakpoints()
const visiblePages = computed(() => {
  if (isMobile.value)  return 6
  if (isTablet.value)  return 8
  if(isDesktop.value)  return 12
  
  return 20
})


// Event handlers
const handleStockClick = (stock: any) => {
  // Navigate to ticker detail page
  router.push({
    name: RouterNames.TickerHistoric,
    params: {
      id: stock.ticker
    }
  })
}

const handlePageChange = (newPage: number) => {
  tickersStore.currentPage = newPage
}

</script>

<style scoped>
/* Page-specific styles can go here */
</style>
