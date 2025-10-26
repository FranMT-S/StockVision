<template>
  <div>
    <v-card>
    <v-card-title class="pa-4 py-4">
      <h3 class="text-body-1 font-weight-medium text-black">
        Key Metrics
      </h3>
    </v-card-title>
    
    <v-card-text class="pa-4 py-0">
      <v-row class="mb-0">
        <v-col cols="4" md="6" lg="12" class="pa-1">
          <v-divider class="mb-3" />
          <KeyMetric 
            title="Market Cap" 
            :value="companyData.marketCap" 
            type="number"
            explanation="The total market value of a company's outstanding shares."
          />
        </v-col>
        <v-col cols="4" md="6"  class="pa-1">
          <KeyMetric 
            title="Volume" 
            :value="companyData.volume" 
            type="number"
            explanation="The total number of shares that have been bought and sold in a specific time period."
          />
        </v-col>
        
        <v-col cols="4" md="6"  class="pa-1">
          <KeyMetric 
            title="Avg Volume" 
            :value="companyData.averageVolume" 
            type="number"
            explanation="The average number of shares that have been bought and sold in a specific time period."
          />
        </v-col>
        
        <v-col cols="4" md="6"  class="pa-1">
          <KeyMetric 
            title="Beta" 
            :value="companyData.beta" 
            type="number"
            explanation="A measure of a stock's volatility in relation to the market."
          />
        </v-col>
        
        <v-col cols="4" md="6"  class="pa-1">
          <KeyMetric 
            title="Dividend" 
            :value="companyData.lastDividend" 
            type="number"
            explanation="The amount of money paid to shareholders for each share of stock."
          />
        </v-col> 
      </v-row>
    </v-card-text>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { humanizeNumberFormat } from '@/shared/helpers/formats';
import type { CompanyData } from '@/shared/models/recomendations'
import KeyMetric from './KeyMetric.vue'

interface Props {
  companyData: CompanyData
}

defineProps<Props>()


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
