<template>
  <div class="tw-mb-[32px]">
    <v-container fluid class="pb-0"> 
      <!-- Header Section -->
      <div class="mb-3">
        <div class="d-flex align-center mb-2">
          <h1 class="text-h5 font-weight-bold text-grey-darken-3 mr-3">
            Company Stock Prices
          </h1>
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
          <StockTableSkeleton v-if="loading" :rows="10" />
          
          <!-- Error State -->
          <ErrorFetchData v-else-if="error" :error="error" class="pa-8 text-center" @click="tickersStore.fetchStockData()"/> 
       
          <!-- Data State -->
          <TickersTable 
            v-else-if="stockData.length > 0" 
            height="calc(100vh - 200px)"
            :stock-data="paginatedStockData"
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
        v-if="!loading && !error && stockData.length > 0"
        :initial-page="currentPage"
        :total-pages="totalPages"
        :total-items="totalItems"
        :items-per-page="itemsPerPage"
        sticky
        class="mt-4"
        @update:current-page="handlePageChange"
      />
    </v-container>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTickersStore } from '../store/tickersStore'
import StockTableSkeleton from '@/shared/components/StockTableSkeleton.vue'
import ErrorFetchData from '@/shared/components/ErrorFetchData.vue'
import TickersTable from '../components/TickersTable.vue'
import Paginator from '@/shared/components/Paginator.vue'

const router = useRouter()
const tickersStore = useTickersStore()

// Initialize data when component mounts
tickersStore.fetchStockData()

// Computed properties
const stockData = computed(() => tickersStore.stockData)
const paginatedStockData = computed(() => tickersStore.paginatedStockData)
const loading = computed(() => tickersStore.loading)
const error = computed(() => tickersStore.error)
const currentPage = computed({
  get: () => tickersStore.currentPage,
  set: (value: number) => tickersStore.setCurrentPage(value)
})
const totalPages = computed(() => tickersStore.totalPages)
const totalItems = computed(() => tickersStore.totalItems)
const itemsPerPage = computed(() => tickersStore.itemsPerPage)

// Event handlers
const handleStockClick = (stock: any) => {
  // Navigate to ticker detail page
  router.push(`/ticker/${stock.ticker}`)
}

const handlePageChange = (page: number) => {
  tickersStore.setCurrentPage(page)
}
</script>

<style scoped>
/* Page-specific styles can go here */
</style>
