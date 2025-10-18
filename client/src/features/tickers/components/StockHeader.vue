<template>
  <v-card class="stock-header" elevation="2">
    <v-card-text class="pa-6">
      <v-row align="center" no-gutters>
        <!-- Company Logo and Basic Info -->
        <v-col cols="12" md="6">
          <div class="d-flex align-center">
            <v-avatar size="64" class="mr-4">
              <img :src="tickerData.ticker.logo" :alt="tickerData.companyData.companyName" />
            </v-avatar>
            <div>
              <h1 class="text-h4 font-weight-bold text-grey-darken-3 mb-1">
                {{ tickerData.companyData.companyName }}
              </h1>
              <div class="d-flex align-center">
                <span class="text-h6 font-weight-medium text-grey-darken-1 mr-2">
                  {{ tickerData.ticker.ticker }}
                </span>
                <v-chip
                  :color="getSentimentColor(tickerData.ticker.sentiment)"
                  size="small"
                  variant="flat"
                >
                  {{ getSentimentText(tickerData.ticker.sentiment) }}
                </v-chip>
              </div>
              <div class="text-body-2 text-grey-darken-1 mt-1">
                {{ tickerData.companyData.exchangeFullName }}
              </div>
            </div>
          </div>
        </v-col>

        <!-- Price Information -->
        <v-col cols="12" md="6" class="d-flex justify-end">
          <div class="text-right">
            <div class="text-h3 font-weight-bold text-grey-darken-3 mb-1">
              ${{ tickerData.ticker.price.toFixed(2) }}
            </div>
            <div 
              class="d-flex align-center justify-end"
              :class="getChangeColorClass(tickerData.companyData.change)"
            >
              <v-icon 
                :icon="tickerData.companyData.change >= 0 ? 'mdi-trending-up' : 'mdi-trending-down'"
                size="small"
                class="mr-1"
              />
              <span class="text-h6 font-weight-medium">
                {{ tickerData.companyData.change >= 0 ? '+' : '' }}${{ tickerData.companyData.change.toFixed(2) }}
                ({{ tickerData.companyData.change >= 0 ? '+' : '' }}{{ tickerData.companyData.changePercentage.toFixed(2) }}%)
              </span>
            </div>
            <div class="text-body-2 text-grey-darken-1 mt-1">
              Market Cap: ${{ formatMarketCap(tickerData.companyData.marketCap) }}
            </div>
          </div>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { RecomendationResponse } from '@/shared/models/recomendations'

interface Props {
  tickerData: RecomendationResponse
}

defineProps<Props>()

// Helper functions
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

const getSentimentText = (sentiment: string): string => {
  switch (sentiment.toLowerCase()) {
    case 'positive':
      return 'Bullish'
    case 'negative':
      return 'Bearish'
    default:
      return 'Neutral'
  }
}

const getChangeColorClass = (change: number): string => {
  return change >= 0 ? 'text-green-darken-2' : 'text-red-darken-2'
}

const formatMarketCap = (marketCap: number): string => {
  if (marketCap >= 1e12) {
    return (marketCap / 1e12).toFixed(2) + 'T'
  } else if (marketCap >= 1e9) {
    return (marketCap / 1e9).toFixed(2) + 'B'
  } else if (marketCap >= 1e6) {
    return (marketCap / 1e6).toFixed(2) + 'M'
  }
  return marketCap.toString()
}
</script>

<style scoped>
.stock-header {
  border-radius: 16px;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

/* Dark theme support */
.v-theme--dark .stock-header {
  background: linear-gradient(135deg, #1e1e1e 0%, #2d2d2d 100%);
  border: 1px solid rgba(255, 255, 255, 0.05);
}
</style>
