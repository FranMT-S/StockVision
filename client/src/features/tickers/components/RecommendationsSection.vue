<template>
  <v-card class="recommendations-table" elevation="2">
    <v-card-title class="pa-4">
      <h3 class="text-h6 font-weight-medium text-grey-darken-3">
        Analyst Recommendations
      </h3>
    </v-card-title>
    
    <v-card-text class="pa-0">
      <div class="recommendations-cards">
  <v-card
    v-for="(recommendation, index) in recommendations"
    :key="recommendation?.id || index"
    class="recommendation-card mb-4"
    elevation="1"
  >
    <v-card-text class="pa-4">
      <!-- Header: Brokerage and Date -->
      <div class="d-flex justify-space-between align-start mb-3">
        <div>
          <div class="text-h6 font-weight-bold text-grey-darken-3">
            {{ recommendation?.brokerage.name }}
          </div>
          <div class="text-caption text-grey-darken-1 mt-1">
            {{ formatDate(recommendation?.time || '') }}
          </div>
        </div>
        
        <v-chip
          :color="getActionColor(recommendation?.action || '')"
          size="default"
          variant="flat"
          class="font-weight-bold"
        >
          {{ recommendation?.action }}
        </v-chip>
      </div>

      <v-divider class="mb-3"></v-divider>

      <!-- Rating Change -->
      <div class="mb-3">
        <div class="text-caption text-grey-darken-1 mb-1">Rating</div>
        <div class="d-flex align-center">
          <v-chip
            size="small"
            variant="outlined"
            :class="getRatingColor(recommendation?.rating_from || '')"
          >
            {{ recommendation?.rating_from }}
          </v-chip>
          <v-icon size="small" color="grey-darken-1" class="mx-2">
            mdi-arrow-right
          </v-icon>
          <v-chip
            size="small"
            variant="outlined"
            :class="getRatingColor(recommendation?.rating_to || '')"
          >
            {{ recommendation?.rating_to }}
          </v-chip>
        </div>
      </div>

      <!-- Target Price -->
      <div>
        <div class="text-caption text-grey-darken-1 mb-1">Target Price</div>
        <div class="text-h6 font-weight-bold text-grey-darken-3">
          ${{ recommendation?.target_from }} → ${{ recommendation?.target_to }}
        </div>
      </div>
    </v-card-text>
  </v-card>
</div>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { Recommendation } from '@/shared/models/recomendations'

interface Props {
  recommendations: (Recommendation | undefined)[]
}

defineProps<Props>()

// Helper functions
const getActionColor = (action: string): string => {
  switch (action.toLowerCase()) {
    case 'buy':
      return 'green'
    case 'sell':
      return 'red'
    case 'hold':
      return 'orange'
    default:
      return 'grey'
  }
}

const getRatingColor = (rating: string): string => {
  switch (rating.toLowerCase()) {
    case 'strong buy':
    case 'buy':
      return 'text-green-darken-2'
    case 'strong sell':
    case 'sell':
      return 'text-red-darken-2'
    case 'hold':
      return 'text-orange-darken-2'
    default:
      return 'text-grey-darken-2'
  }
}

const formatDate = (dateString: string): string => {
  if (!dateString) return ''
  
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}
</script>

<style scoped>


.recommendations-cards {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.recommendation-card {
  border-radius: 12px;
  transition: all 0.2s ease;
}



/* Responsive: 2 columnas en tablets y más */
@media (min-width: 768px) {
  .recommendations-cards {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }
}

/* Responsive: 3 columnas en desktop */
@media (min-width: 1200px) {
  .recommendations-cards {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>
