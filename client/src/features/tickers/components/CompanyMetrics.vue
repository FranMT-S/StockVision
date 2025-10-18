<template>
  <v-card class="company-metrics" elevation="2">
    <v-card-title class="pa-4">
      <h3 class="text-h6 font-weight-medium text-grey-darken-3">
        Key Metrics
      </h3>
    </v-card-title>
    
    <v-card-text class="pa-4">
      <v-row>
        <v-col cols="6" class="pb-2">
          <div class="metric-item">
            <div class="text-caption text-grey-darken-1 mb-1">Volume</div>
            <div class="text-body-1 font-weight-medium text-grey-darken-3">
              {{ formatVolume(companyData.volume) }}
            </div>
          </div>
        </v-col>
        
        <v-col cols="6" class="pb-2">
          <div class="metric-item">
            <div class="text-caption text-grey-darken-1 mb-1">Avg Volume</div>
            <div class="text-body-1 font-weight-medium text-grey-darken-3">
              {{ formatVolume(companyData.averageVolume) }}
            </div>
          </div>
        </v-col>
        
        <v-col cols="6" class="pb-2">
          <div class="metric-item">
            <div class="text-caption text-grey-darken-1 mb-1">Beta</div>
            <div class="text-body-1 font-weight-medium text-grey-darken-3">
              {{ companyData.beta.toFixed(2) }}
            </div>
          </div>
        </v-col>
        
        <v-col cols="6" class="pb-2">
          <div class="metric-item">
            <div class="text-caption text-grey-darken-1 mb-1">Dividend</div>
            <div class="text-body-1 font-weight-medium text-grey-darken-3">
              ${{ companyData.lastDividend.toFixed(2) }}
            </div>
          </div>
        </v-col>
        
        <v-col cols="6" class="pb-2">
          <div class="metric-item">
            <div class="text-caption text-grey-darken-1 mb-1">Sector</div>
            <div class="text-body-1 font-weight-medium text-grey-darken-3">
              {{ companyData.sector }}
            </div>
          </div>
        </v-col>
        
        <v-col cols="6" class="pb-2">
          <div class="metric-item">
            <div class="text-caption text-grey-darken-1 mb-1">Industry</div>
            <div class="text-body-1 font-weight-medium text-grey-darken-3">
              {{ companyData.industry }}
            </div>
          </div>
        </v-col>
        
        <v-col cols="12" class="pt-2">
          <v-divider class="mb-3" />
          <div class="metric-item">
            <div class="text-caption text-grey-darken-1 mb-1">Market Cap</div>
            <div class="text-h6 font-weight-bold text-grey-darken-3">
              ${{ formatMarketCap(companyData.marketCap) }}
            </div>
          </div>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { CompanyData } from '@/shared/models/recomendations'

interface Props {
  companyData: CompanyData
}

defineProps<Props>()

// Helper functions
const formatVolume = (volume: number): string => {
  if (volume >= 1e9) {
    return (volume / 1e9).toFixed(1) + 'B'
  } else if (volume >= 1e6) {
    return (volume / 1e6).toFixed(1) + 'M'
  } else if (volume >= 1e3) {
    return (volume / 1e3).toFixed(1) + 'K'
  }
  return volume.toString()
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
.company-metrics {
  border-radius: 16px;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.metric-item {
  padding: 8px 0;
}

/* Dark theme support */
.v-theme--dark .company-metrics {
  background: linear-gradient(135deg, #1e1e1e 0%, #2d2d2d 100%);
  border: 1px solid rgba(255, 255, 255, 0.05);
}
</style>
