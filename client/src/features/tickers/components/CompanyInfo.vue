<template>
  <div class="tw-p-0 tw-m-0 tw-mt-1 md:tw-mt-0 tw-flex tw-flex-col">
    <div class="tw-flex tw-flex-row tw-items-center tw-justify-between lg:tw-justify-start tw-relative tw-ps-[10px]">
      <section class="lg:tw-min-w-[20%] lg:tw-max-w-[60px] lg:tw-max-h-[60px] tw-self-start">
        <img  width="60" height="60" :src="companyData.image" alt="Company Logo" class="company-logo logo tw-self-baseline">
      </section>
      <section class="tw-flex tw-flex-row tw-items-start md:tw-items-center tw-w-full">
        <div class="md:text-h6  font-weight-medium  dark:tw-text-[#ffffff] tw-w-full">
          <section class=" ">
            <div class="tw-flex tw-flex-row tw-items-center tw-gap-2">
              {{ companyData.companyName }}
              <span class="tw-text-[#717171] tw-text-[12px]  dark:tw-text-[#ffffff]">{{ companyData.symbol }}</span>
            </div>
          <section id="company-price" class="tw-flex tw-flex-row tw-items-center tw-gap-1">
              <div 
                class="d-flex align-center tw-gap-3"
              >
                <h3 class="tw-text-[13px] tw-font-medium tw-text-grey-darken-3 text-primary">
                  ${{ companyData.price.toFixed(2) }}
                </h3>
                <div class="tw-text-[12px] font-medium tw-flex tw-flex-row tw-items-baseline" :class="{ 'text-green-darken-2': changePriceStyles.isChangeToUp, 'text-red-darken-2': !changePriceStyles.isChangeToUp }">
                  <v-icon 
                    :icon="changePriceStyles.icon"
                    size="x-small"
                  />
                  <div >{{ changePriceStyles.symbol }}${{ companyData.change.toFixed(2) }}</div>
                  <div>({{ changePriceStyles.symbol }}{{ companyData.changePercentage.toFixed(2) }}%)</div>
                </div>
              </div>
            </section>
          </section>
          <section class=" tw-flex tw-flex-row tw-gap-2 lg:tw-flex-col ">
            <v-chip class="tw-text-[#717171] pa-0" v-if="companyData.sector" color="secondary" size="small" variant="text" >
              <v-tooltip text="Sector" location="top">  
                <template #activator="{ props }">
                  <v-icon v-bind="props" icon="mdi-office-building" size="16" />
                </template>
              </v-tooltip>
              {{ companyData.sector  }}
            </v-chip>
            <v-chip v-if="companyData.industry" class="tw-text-[#717171] pa-0 lg:!tw-h-full" color="secondary" size="small" variant="text">
              <v-tooltip text="Industry" location="top">
                <template #activator="{ props }">
                  <v-icon class="lg:tw-self-center lg:tw-z-30" v-bind="props" icon="mdi-briefcase"  />
                </template>
              </v-tooltip>
               <div class="tw-whitespace-normal tw-break-words lg-tw-max-w-[120px] tw-line-clamp-1">
                {{ companyData.industry }}
              </div>
            </v-chip>
          </section>
        </div>
      </section>    
      <section class="tw-self-baseline tw-absolute tw-right-[10px] tw-top-[0px] lg:tw-top-[-10px] xl:tw-top-[0px]">
        <v-tooltip text="Visit Web" location="top" class="tw-self-baseline ">
          <template #activator="{ props }">
            <a v-bind="props" :href="companyData.website" target="_blank" rel="noopener noreferrer"
            >
              <v-icon size="18" color="primary" icon="mdi-web"></v-icon>
            </a>
          </template>
        </v-tooltip>
      </section>  
    </div>
  </div>
</template>

<script setup lang="ts">
import type { CompanyData } from '@/shared/models/recomendations'
import { computed } from 'vue'

interface Props {
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
