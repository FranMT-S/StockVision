<template>
  <div>
    <div class="tw-flex tw-flex-col tw-gap-1">
      <div class="tw-flex tw-flex-row tw-gap-2 tw-items-center">
        <span class="tw-text-[13px] tw-font-medium ">Price:</span>
        <span class="tw-text-[12px] tw-font-medium">{{ props.data.close?.toFixed(2) }}</span>
     </div>
      <div class="tw-flex tw-flex-row tw-gap-2 tw-items-center">
        <span class="tw-text-[13px] tw-font-medium ">Volume:</span>
        <span class="tw-text-[12px] tw-font-medium">{{ props.data.volume.toLocaleString() }}</span>
      </div>
      <div  class="tw-flex tw-flex-row tw-gap-2 tw-items-center">
        <span class="tw-text-[13px] tw-font-medium ">Open:</span>
        <span class="tw-text-[12px] tw-font-medium">{{ props.data.open?.toFixed(2) }}</span>
      </div>
      <div  class="tw-flex tw-flex-row tw-gap-2 tw-items-center">
        <span class="tw-text-[13px] tw-font-medium ">High:</span>
        <span class="tw-text-[12px] tw-font-medium">{{ props.data.high?.toFixed(2) }}</span>
      </div>
      <div  class="tw-flex tw-flex-row tw-gap-2 tw-items-center">
        <span class="tw-text-[13px] tw-font-medium ">Low:</span>
        <span class="tw-text-[12px] tw-font-medium">{{ props.data.low?.toFixed(2) }}</span>
      </div>  
    </div> 
  </div>
</template>
  
<script setup lang="ts">
import { HistoricalPrice } from '@/shared/models/recomendations';

interface Props {
  data: HistoricalPrice
}

const props = defineProps<Props>()
</script>
