
<template>
  <v-card>
    <v-card-title class="tw-text-lg tw-font-semibold  tw-text-blue-500 ">{{ tittle }}</v-card-title>
    <v-card-text>
      <span 
      data-test="action"
      class="tw-font-semibold tw-mr-1 text-advice"
      :class="[advice.color.value.class]"
      >
        {{ advice.action }}<v-icon data-test="icon" :icon="advice.icon.value" size="small" />
      </span>
      <span data-test="description" class="tw-text-gray-800 text-description">{{ advice.description }}</span>
    </v-card-text>
  </v-card >
</template>

<script setup lang="ts">

import { computed } from 'vue';
import { useAdvice } from '../composable/useAdvice';

interface props{
  advice: string;
  tittle: string;
}


const props = defineProps<props>();

const advice = computed(() => {
  let [initialAction, description] = props.advice.split('.');
  initialAction = initialAction ? initialAction.trim() : '';
  description = description ? description.trim() : '';

  const { action,icon, color } = useAdvice(initialAction)

  return {
    action,
    description,
    icon,
    color
  }
})


</script>