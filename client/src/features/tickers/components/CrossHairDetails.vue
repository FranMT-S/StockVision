<template>
  <div>
    <div class="tw-flex tw-flex-col tw-gap-1">
      <div class="tw-flex tw-flex-row tw-gap-2 tw-items-center">
        <span class="tw-text-[13px] tw-font-medium ">V:</span>
        <span class="tw-text-[12px] tw-font-medium">{{  humanizeNumberFormat(props.data.volume)}}</span>
      </div>
      <div  class="tw-flex tw-flex-row tw-gap-2 tw-items-center">
        <span class="tw-text-[13px] tw-font-medium ">O:</span>
        <span class="tw-text-[12px] tw-font-medium">{{ humanizeNumberFormat(props.data.open) }}$</span>
      </div>
      <div class="tw-flex tw-flex-row tw-gap-2 tw-items-center">
        <span class="tw-text-[13px] tw-font-medium ">C:</span>
        <span class="tw-text-[12px] tw-font-medium">{{ humanizeNumberFormat(props.data.close) }}$</span>
     </div>
      <div  class="tw-flex tw-flex-row tw-gap-2 tw-items-center">
        <span class="tw-text-[13px] tw-font-medium ">H:</span>
        <span class="tw-text-[12px] tw-font-medium">{{ humanizeNumberFormat(props.data.high) }}$</span>
      </div>
      <div  class="tw-flex tw-flex-row tw-gap-2 tw-items-center">
        <span class="tw-text-[13px] tw-font-medium ">L:</span>
        <span class="tw-text-[12px] tw-font-medium">{{ humanizeNumberFormat(props.data.low) }}$</span>
      </div>  
      <div  class="tw-flex tw-flex-row tw-gap-2 tw-items-center">
        <span class="tw-text-[13px] tw-font-medium ">CH:</span>
        <span class="tw-text-[12px] tw-font-medium">{{ humanizeNumberFormat(props.data.change) }}</span>
      </div>  
      <div  class="tw-flex tw-flex-row tw-gap-2 tw-items-center">
        <span class="tw-text-[13px] tw-font-medium ">%CH:</span>
        <span class="tw-text-[12px] tw-font-medium">{{ humanizeNumberFormat(props.data.changePercent) }}</span>
      </div>  
    </div> 
  </div>
</template>
  
<script setup lang="ts">
import { humanizeNumberFormat } from '@/shared/helpers/formats';
import { StockHLOC } from '@/shared/models/recomendations';

interface Props {
  data: StockHLOC
}

const props = defineProps<Props>()
</script>
