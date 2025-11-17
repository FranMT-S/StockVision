<template>
  <v-data-table
    :headers="headers"
    :items="rows"
    :items-per-page="-1"
    :hide-default-footer="true"
    class="recommendations-table"
    density="compact"
    @click:row="handleRowClick"
  >
    <!-- Stock Column -->
    <template v-slot:item.ticker="{ item }">
      <div class="d-flex align-center">
        <v-avatar 
          size="32" 
          class="mr-3"
          :color="getAvatarColor(item.ticker)"
        >
          <v-img
            v-if="item.url"
            :src="item.url"
            :alt="item.companyName"
            class="rounded"
          />
          <span v-else class="text-white font-weight-bold text-caption">
            {{ getCompanyInitials(item.companyName) }}
          </span>
        </v-avatar>
        <div>
          <div class="tw-flex tw-gap-2 tw-items-baseline font-weight-medium text-grey-darken-3">
            <span class="tw-text-[15px]">{{ item.companyName }}</span>
            <span class="tw-text-[12px] tw-text-gray-500">{{ item.ticker }}</span>
          </div>
        </div>
      </div>
    </template>

    <!-- Price Column -->
    <template v-slot:item.price="{ item }">
      <div class="text-body-2 tw-text-[15px] font-weight-medium text-grey-darken-3">
        ${{ item.price.toFixed(2) }}
      </div>
    </template>

    <!-- Change Column -->
    <template v-slot:item.change="{ item }">
      <div 
        class="d-flex align-center"
        :class="getChangeColorClass(item.change)"
      >
        <v-icon 
          :icon="item.change >= 0 ? 'mdi-trending-up' : 'mdi-trending-down'"
          size="x-small"
          class="mr-1"
        />
        <span class="text-body-2 tw-text-[15px] font-weight-medium">
          {{ item.change >= 0 ? '+' : '' }}${{ item.change.toFixed(2) }}
          ({{ item.change >= 0 ? '+' : '' }}{{ item.changePercentage.toFixed(2) }}%)
        </span>
      </div>
    </template>

    <!-- Sentiment Column -->
    <template v-slot:item.sentiment="{ item }">
      <div class="d-flex align-center justify-center ga-2">
        <v-avatar size="24">
          <v-icon 
            :icon="getSentimentIcon(item.sentiment)"
            :color="getSentimentColor(item.sentiment)"
            size="small"
          />
        </v-avatar>
        <v-chip 
          size="small"
          class="tw-capitalize"  
          :color="getSentimentColor(item.sentiment)" 
          variant="outlined"
        >
          {{ item.sentiment }}
        </v-chip>
      </div>
    </template>

    <!-- Action Advice Column -->
    <template v-slot:item.advice="{ item }">
      <div class="d-flex align-center justify-center ga-2">
        <v-avatar size="24">
          <v-icon 
            :icon="getAdviceActionIcon(item.advice)"
            :class="getAdviceActionColor(item.advice).class"
            size="small"
          />
        </v-avatar>
        <v-chip 
          size="small"
          class="tw-capitalize"  
          :color="getAdviceActionColor(item.advice).vuetify" 
          variant="outlined"
        >
          {{ item.advice }}
        </v-chip>
      </div>
    </template>

    <!-- Last Rating Date Column -->
    <template v-slot:item.lastRatingDate="{ item }">
      {{ getDateString(item.lastRatingDate) }}
    </template>
  </v-data-table>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { getAdviceActionColor, getAdviceActionIcon } from '@/shared/helpers/adviceActions'
import { getSentimentColor, getSentimentIcon } from '@/shared/helpers/sentiment'

export interface TickerRow {
  ticker: string
  companyName: string
  price: number
  url?: string
  change: number
  changePercentage: number
  sentiment: string
  lastRatingDate: string | Date
  advice: string
}

interface Props {
  rows: TickerRow[]
}

interface Emits {
  (e: 'rowClick', stock: TickerRow): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

type Header = {
  title: string
  key: string
  align: 'start' | 'center'
  sortable: boolean
}

// Define table headers with sorting
const headers = computed<Header[]>(() => [
  { 
    title: 'Stock', 
    key: 'ticker', 
    align: 'start',
    sortable: true 
  },
  { 
    title: 'Price', 
    key: 'price', 
    align: 'start',
    sortable: true 
  },
  { 
    title: 'Change', 
    key: 'change', 
    align: 'start',
    sortable: true 
  },
  { 
    title: 'Sentiment', 
    key: 'sentiment', 
    align: 'center',
    sortable: true 
  },
  { 
    title: 'Action Advice', 
    key: 'advice', 
    align: 'center',
    sortable: true 
  },
  { 
    title: 'Last Rating Date', 
    key: 'lastRatingDate', 
    align: 'center',
    sortable: true 
  }
])

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

const getDateString = (date: string | Date): string => {
  if (!date) return ''

  const d = new Date(date)

  if (isNaN(d.getTime())) {
    return 'Not available'
  }

  return d.toLocaleDateString()
}

const handleRowClick = (_event: any, { item }: { item: TickerRow }) => {
  emit('rowClick', item)
}
</script>

<style scoped>
.recommendations-table {
  border-collapse: separate;
  border-spacing: 0;
}

:deep(.v-data-table__tr) {
  cursor: pointer;
  transition: background-color 0.2s ease;
}

:deep(.v-data-table__tr:hover) {
  background-color: #f5f5f5;
}

:deep(.v-data-table thead th) {
  background-color: #fafafa !important;
  font-weight: 600 !important;
  color: #616161 !important;
}

:deep(.v-data-table__td) {
  padding: 8px 12px !important;
}
</style>