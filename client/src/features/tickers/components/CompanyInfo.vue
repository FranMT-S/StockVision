<template>
  <v-card class="company-info" elevation="2">
    <v-card-title class="pa-4">
      <h3 class="text-h6 font-weight-medium text-grey-darken-3">
        Company Information
      </h3>
    </v-card-title>
    
    <v-card-text class="pa-4">
      <v-list class="pa-0" density="compact">
        <v-list-item class="px-0">
          <template #prepend>
            <v-icon color="grey-darken-1">mdi-domain</v-icon>
          </template>
          <v-list-item-title class="text-body-2 text-grey-darken-1">
            Website
          </v-list-item-title>
          <v-list-item-subtitle>
            <a 
              :href="companyData.website" 
              target="_blank"
              class="text-decoration-none text-primary"
            >
              {{ companyData.website }}
            </a>
          </v-list-item-subtitle>
        </v-list-item>

        <v-list-item class="px-0">
          <template #prepend>
            <v-icon color="grey-darken-1">mdi-account-tie</v-icon>
          </template>
          <v-list-item-title class="text-body-2 text-grey-darken-1">
            CEO
          </v-list-item-title>
          <v-list-item-subtitle class="text-grey-darken-3">
            {{ companyData.ceo }}
          </v-list-item-subtitle>
        </v-list-item>

        <v-list-item class="px-0">
          <template #prepend>
            <v-icon color="grey-darken-1">mdi-map-marker</v-icon>
          </template>
          <v-list-item-title class="text-body-2 text-grey-darken-1">
            Country
          </v-list-item-title>
          <v-list-item-subtitle class="text-grey-darken-3">
            {{ companyData.country }}
          </v-list-item-subtitle>
        </v-list-item>

        <v-list-item class="px-0">
          <template #prepend>
            <v-icon color="grey-darken-1">mdi-bank</v-icon>
          </template>
          <v-list-item-title class="text-body-2 text-grey-darken-1">
            Exchange
          </v-list-item-title>
          <v-list-item-subtitle class="text-grey-darken-3">
            {{ companyData.exchange }}
          </v-list-item-subtitle>
        </v-list-item>

        <v-list-item class="px-0">
          <template #prepend>
            <v-icon color="grey-darken-1">mdi-tag</v-icon>
          </template>
          <v-list-item-title class="text-body-2 text-grey-darken-1">
            Symbol
          </v-list-item-title>
          <v-list-item-subtitle class="text-grey-darken-3">
            {{ companyData.symbol }}
          </v-list-item-subtitle>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { CompanyData } from '@/shared/models/recomendations'

interface Props {
  companyData: CompanyData
}

defineProps<Props>()
</script>

<style scoped>
.company-info {
  border-radius: 16px;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

/* Dark theme support */
.v-theme--dark .company-info {
  background: linear-gradient(135deg, #1e1e1e 0%, #2d2d2d 100%);
  border: 1px solid rgba(255, 255, 255, 0.05);
}
</style>
