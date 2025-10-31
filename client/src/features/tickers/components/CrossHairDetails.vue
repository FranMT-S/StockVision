<template>
  <div>
    <div class="tw-flex tw-flex-col tw-gap-1">
      <v-icon :icon="icon" size="14" :class="textColorByTending" />
      <div v-for="field in fields" :key="field.label" class="tw-flex tw-flex-row tw-gap-2 tw-items-center">
        <span  class="tw-text-[13px] tw-font-medium ">{{ field.label }}</span>
        <span :class="field.extraClass" class="tw-text-[11.5px] tw-font-semibold ">{{ field.value }}</span>
      </div>
      
    </div> 
  </div>
</template>
  
<script setup lang="ts">
interface Props {
  data: StockHLOC
}

import { humanizeNumberFormat } from '@/shared/helpers/formats';
import { StockHLOC } from '@/shared/models/recomendations';
import { computed } from 'vue';

const props = defineProps<Props>()

const icon = computed(() => {
  return props.data.change > 0 ? 'mdi-trending-up' : 'mdi-trending-down';
})

const textColorByTending = computed(() => {
  return props.data.change > 0 ? 'tw-text-green-400' : 'tw-text-red-400';
})

const textColorDefault = 'tw-text-gray-200'

const fields = computed(() => {
  
  return [
    { label: 'V:',value: `${humanizeNumberFormat(props.data.volume,2)}`, extraClass: textColorDefault },
    { label: 'O:',value: `${humanizeNumberFormat(props.data.open,2)}$`, extraClass: textColorDefault},
    { label: 'C:',value: `${humanizeNumberFormat(props.data.close,2)}$`, extraClass: textColorDefault},
    { label: 'H:',value: `${humanizeNumberFormat(props.data.high,2)}$`, extraClass:textColorDefault },
    { label: 'L:',value: `${humanizeNumberFormat(props.data.low,2)}$`, extraClass: textColorDefault},
    { label: 'CH:',value: humanizeNumberFormat(props.data.change,2), extraClass:textColorByTending.value },
    { label: '%CH:',value: `${humanizeNumberFormat(props.data.changePercent)}%`, extraClass:textColorByTending.value },
  ]
})

</script>
