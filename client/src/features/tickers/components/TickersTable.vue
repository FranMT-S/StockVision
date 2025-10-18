<template>
  <v-table class="recommendations-table" >
    <thead>
      <tr class="text-center">
        <th class="text-left font-weight-bold text-grey-darken-2 ps-14 py-1">Stock</th>
        <th class="text-left font-weight-bold text-grey-darken-2 pa-1">Price</th>
        <th class="text-left font-weight-bold text-grey-darken-2 pa-1">Change</th>
        <th class="text-center font-weight-bold text-grey-darken-2 pa-1">Sentiment</th>
      </tr>
    </thead>
    <tbody>
      <tr 
        v-for="stock in stockData" 
        :key="stock.ticker"
        class="stock-row"
        @click="handleRowClick(stock)"
      >
        <!-- Stock Column -->
        <td class="pa-2">
          <div class="d-flex align-center">
            <v-avatar 
              size="32" 
              class="mr-3"
              :color="getAvatarColor(stock.ticker)"
            >
              <span class="text-white font-weight-bold text-caption">
                {{ getCompanyInitials(stock.companyName) }}
              </span>
            </v-avatar>
            <div>
              <div class="text-body-2 font-weight-medium text-grey-darken-3">
                {{ stock.companyName }}
              </div>
              <div class="text-caption text-grey-darken-1">
                {{ stock.ticker }}
              </div>
            </div>
          </div>
        </td>

        <!-- Price Column -->
        <td class="text-left pa-2">
          <div class="text-body-2 font-weight-medium text-grey-darken-3">
            ${{ stock.price.toFixed(2) }}
          </div>
        </td>

        <!-- Change Column -->
        <td class="text-left pa-2">
          <div 
            class="d-flex align-center"
            :class="getChangeColorClass(stock.change)"
          >
            <v-icon 
              :icon="stock.change >= 0 ? 'mdi-trending-up' : 'mdi-trending-down'"
              size="x-small"
              class="mr-1"
            />
            <span class="text-body-2 font-weight-medium">
              {{ stock.change >= 0 ? '+' : '' }}${{ stock.change.toFixed(2) }}
              ({{ stock.change >= 0 ? '+' : '' }}{{ stock.changePercentage.toFixed(2) }}%)
            </span>
          </div>
        </td>

        <!-- Sentiment Column -->
        <td class="text-center pa-2">
          <v-avatar size="24">
            <v-icon 
              :icon="getSentimentIcon(stock.sentiment)"
              :color="getSentimentColor(stock.sentiment)"
              size="small"
            />
          </v-avatar>
        </td>
      </tr>
    </tbody>
  </v-table>
</template>

<script setup lang="ts">
interface StockData {
  ticker: string
  companyName: string
  price: number
  change: number
  changePercentage: number
  sentiment: string
}

interface Props {
  stockData: StockData[]
}

interface Emits {
  (e: 'rowClick', stock: StockData): void
}

defineProps<Props>()
const emit = defineEmits<Emits>()

// Helper functions
const getCompanyInitials = (companyName: string): string => {
  const words = companyName.split(' ')
  if (words.length >= 2) {
    return (words[0][0] + words[1][0]).toUpperCase()
  }
  return companyName.substring(0, 2).toUpperCase()
}

const getAvatarColor = (ticker: string): string => {
  const colors = ['blue', 'green', 'orange', 'purple', 'red', 'teal']
  const index = ticker.length % colors.length
  return colors[index] + '-lighten-1'
}

const getChangeColorClass = (change: number): string => {
  return change >= 0 ? 'text-green-darken-2' : 'text-red-darken-2'
}

const getSentimentIcon = (sentiment: string): string => {
  switch (sentiment.toLowerCase()) {
    case 'positive':
      return 'mdi-emoticon-happy'
    case 'negative':
      return 'mdi-emoticon-sad'
    default:
      return 'mdi-emoticon-neutral'
  }
}

const getSentimentColor = (sentiment: string): string => {
  switch (sentiment.toLowerCase()) {
    case 'positive':
      return 'green'
    case 'negative':
      return 'red'
    default:
      return 'orange'
  }
}

const handleRowClick = (stock: StockData) => {
  emit('rowClick', stock)
}
</script>

<style scoped>
.recommendations-table {
  border-collapse: separate;
  border-spacing: 0;
}

.stock-row {
  border-bottom: 1px solid #e0e0e0;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.stock-row:last-child {
  border-bottom: none;
}

.stock-row:hover {
  background-color: #f5f5f5;
}

:deep(.v-table thead th) {
  border-bottom: 1px solid #e0e0e0;
  background-color: #fafafa;
  padding: 8px 12px !important;
  height: 36px !important;
  font-size: 0.875rem !important;
  line-height: 1.25rem !important;
}
</style>
