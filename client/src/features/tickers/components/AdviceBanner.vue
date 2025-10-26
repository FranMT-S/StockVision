
<template>
  <v-card>
    <v-card-title class="tw-text-lg tw-font-semibold  tw-text-blue-500 ">{{ tittle }}</v-card-title>
    <v-card-text>{{ advice }}</v-card-text>
  </v-card >
</template>

<script setup lang="ts">

import { useGlobalStore } from '@/shared/store/global';
import { computed } from 'vue';
const globalStore = useGlobalStore();


interface props{
  advice: string;
  tittle: string;
}


const props = defineProps<props>();

</script>