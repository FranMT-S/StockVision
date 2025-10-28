<script setup lang="ts">
import { Anchor } from 'vuetify/lib/types.mjs';
import { useToast } from '../composables/useToast'
const { state, close } = useToast()

</script>

<template>
  <v-snackbar
    v-model="state.show"
    :color="state.color"
    :timeout="state.timeout"
    :show="true"
    :location="state.location"
  >
    {{ state.message }}
    <template #actions>
      <v-btn color="white" text @click="close">Cerrar</v-btn>
    </template>
  </v-snackbar>
</template>
