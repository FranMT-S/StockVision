<template>
  <div class="pa-0 ma-0 tw-flex tw-flex-col">
    <v-card-title v-if="title" class="pa-4">
      <h3 class="text-h6 font-weight-medium  tw-text-white dark:tw-text-[#ffffff]">
        {{ title }}
      </h3>
    </v-card-title>
    
    <div class="tw-flex tw-flex-row tw-items-center tw-justify-between">
      <div class="tw-flex tw-flex-row tw-items-center">
        <div>
          <img width="60" height="60" :src="companyData.image" alt="Company Logo" class="company-logo logo tw-self-baseline">
        </div>
        <div class="text-h6 font-weight-medium  dark:tw-text-[#ffffff]">
          <section>
            {{ companyData.companyName }}
            <span class="tw-text-[#717171] tw-text-[12px]  dark:tw-text-[#ffffff]">{{ companyData.symbol }}</span>
          </section>
          <section class=" tw-flex tw-flex-row tw-gap-2  ">
            <v-chip class="tw-text-[#717171] pa-0" v-if="companyData.sector" color="secondary" size="small" variant="text" >
              <v-tooltip text="Sector" location="top">
                <template #activator="{ props }">
                  <v-icon v-bind="props" icon="mdi-office-building" size="16" />
                </template>
              </v-tooltip>
              {{ companyData.sector  }}
            </v-chip>
            <v-chip v-if="companyData.industry" class="tw-text-[#717171] pa-0" color="secondary" size="small" variant="text">
              <v-tooltip text="Industry" location="top">
                <template #activator="{ props }">
                  <v-icon v-bind="props" icon="mdi-briefcase" size="16" />
                </template>
              </v-tooltip>
              {{ companyData.industry }}
            </v-chip>
          </section>
        </div>
      </div>    
      <div>
        <v-tooltip text="Website" location="top" class="tw-self-baseline">
          <template #activator="{ props }">
            <a v-bind="props" :href="companyData.website" target="_blank" rel="noopener noreferrer"
            >
              <v-icon size="18" color="primary" icon="mdi-web"></v-icon>
            </a>
          </template>
        </v-tooltip>
      </div>  
      
    </div>
  

    <div class="px-4 tw-flex tw-flex-row tw-items-center tw-gap-2 mt-2">
      <div class="tw-flex tw-flex-row tw-items-center tw-gap-2">
        <div 
          class="d-flex align-center tw-gap-3"
        >
          <h3 class="text-h5 font-weight-medium text-grey-darken-3 text-primary">
            ${{ companyData.price.toFixed(2) }}
          </h3>
          <span class="text-body-2 font-weight-medium" :class="{ 'text-green-darken-2': changePriceStyles.isChangeToUp, 'text-red-darken-2': !changePriceStyles.isChangeToUp }">
            <v-icon 
              :icon="changePriceStyles.icon"
              size="x-small"
            />
            {{ changePriceStyles.symbol }}${{ companyData.change.toFixed(2) }}
            ({{ changePriceStyles.symbol }}{{ companyData.changePercentage.toFixed(2) }}%)
          </span>
        </div>
      </div>
      
    </div>

  </div>
</template>

<script setup lang="ts">
import type { CompanyData } from '@/shared/models/recomendations'
import { computed } from 'vue'

interface Props {
  title?: string
  companyData: CompanyData
}

const props = defineProps<Props>()

const changePriceStyles = computed(() => {
  const symbol = props.companyData.change >= 0 ? '+' : ''
  const icon = props.companyData.change >= 0 ? 'mdi-trending-up' : 'mdi-trending-down'
  const isChangeToUp = props.companyData.change >= 0
  return {
    symbol,
    changePercentage: props.companyData.changePercentage,
    icon,
    isChangeToUp
  }
})

</script>

<style scoped>
.company-info {
  border-radius: 16px;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.logo {
  mix-blend-mode: multiply;
}

/* Dark theme support */
.v-theme--dark .company-info {
  background: linear-gradient(135deg, #1e1e1e 0%, #2d2d2d 100%);
  border: 1px solid rgba(255, 255, 255, 0.05);
}
</style>
