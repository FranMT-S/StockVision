<template>
  <div class="pa-8 text-center">
    <v-icon size="64" color="red-lighten-2" class="mb-4">
      mdi-alert-circle-outline
    </v-icon>
    <h3 class="text-h6 text-grey-darken-2 mb-2">Failed to load data</h3>
    <p class="text-body-2 text-grey-darken-1 mb-4">{{ error }}</p>
    <v-btn
      data-testid="try-again-button"
      color="primary"
      variant="outlined"
      @click="() => emit('click')"
    >
      <v-icon left>mdi-refresh</v-icon>
      Try Again
    </v-btn>
  </div>
</template>

<script setup lang="ts">
  import { VIcon } from 'vuetify/components';

  const props = defineProps<{
    error: string | null
  }>()

  const emit = defineEmits<{
    (e: 'click'): () => void
  }>()

</script>
