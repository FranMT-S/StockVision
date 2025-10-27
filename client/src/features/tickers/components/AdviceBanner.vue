
<template>
  <v-card>
    <v-card-title class="tw-text-lg tw-font-semibold  tw-text-blue-500 ">{{ tittle }}</v-card-title>
    <v-card-text>
      <span 
      class="tw-font-semibold tw-mr-1"
      :class="[advice.color.value.class]"
      >
        {{ advice.action }}<v-icon :icon="advice.icon.value" size="small" />
      </span>
      <span class="tw-text-gray-800">{{ advice.description }}</span>
    </v-card-text>
  </v-card >
</template>

<script setup lang="ts">

import { useGlobalStore } from '@/shared/store/global';
import { computed } from 'vue';
import { useAdvice } from '../composable/useAdvice';
const globalStore = useGlobalStore();


interface props{
  advice: string;
  tittle: string;
}


const props = defineProps<props>();

const advice = computed(() => {
  let [initialAction, description] = props.advice.split('.');
  initialAction = initialAction.trim();
  description = description.trim();

  const { action,icon, color } = useAdvice(initialAction)

  return {
    action,
    description,
    icon,
    color
  }
})


</script>