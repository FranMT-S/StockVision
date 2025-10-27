<template>
  <div class="metric-item tw-flex tw-flex-col pa-1  rounded tw-bg-neutral-100">
    <div class="tw-flex tw-flex-row tw-items-center ga-1 justify-center">
      <div class="text-caption text-grey-darken-1 ">{{ title }}</div>
      <v-tooltip interactive  
         open-on-click
        :close-on-content-click="false" max-width="150px" location="start">
        <template v-slot:activator="{ props }">
          <v-icon size="small" v-bind="props" icon="mdi-information-outline" />
        </template>
        <div>
          {{ explanation }}
        </div>
      </v-tooltip>
    </div>
    <div class="text-body-2 font-weight-medium text-grey-darken-3 justify-center text-center">
      {{ value }}
    </div>
  </div>  
</template>

<script setup lang="ts">
import { humanizeNumberFormat } from '@/shared/helpers/formats';
import { computed } from 'vue';

interface Props {
  title: string
  value: number | string
  type: 'number' | 'currency' | "string"
  explanation?: string
}

 const props = withDefaults(defineProps<Props>(), {
  type: 'number' 
})

const value = computed(() => {
  if (props.type === 'number' ) {
    return humanizeNumberFormat(Number(props.value))
  } else if (props.type === 'currency') {
    return '$' + humanizeNumberFormat(Number(props.value))
  } else {
    return props.value
  }
})

</script>

<style scoped>

</style>