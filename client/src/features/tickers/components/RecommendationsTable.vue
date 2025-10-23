<template>
  <v-card class="recommendations-table" elevation="2">
    <v-card-title class="pa-4">
      <h3 class="text-h6 font-weight-medium text-grey-darken-3">
        Analyst Recommendations
      </h3>
    </v-card-title>
    
    <v-card-text class="pa-0">
      <v-table class="recommendations-table-content">
        <thead>
          <tr>
            <th class="text-left font-weight-bold text-grey-darken-2 pa-3">Brokerage</th>
            <th class="text-left font-weight-bold text-grey-darken-2 pa-3">Action</th>
            <th class="text-left font-weight-bold text-grey-darken-2 pa-3">Rating</th>
            <th class="text-left font-weight-bold text-grey-darken-2 pa-3">Target Price</th>
            <th class="text-left font-weight-bold text-grey-darken-2 pa-3">Date</th>
          </tr>
        </thead>
        <tbody>
          <tr 
            v-for="(recommendation, index) in recommendations" 
            :key="recommendation?.id || index"
            class="recommendation-row"
          >
            <td class="pa-3">
              <div class="text-body-2 font-weight-medium text-grey-darken-3">
                {{ recommendation?.brokerage.name }}
              </div>
            </td>
            
            <td class="pa-3">
              <v-chip
                :color="getActionColor(recommendation?.action || '')"
                size="small"
                variant="flat"
              >
                {{ recommendation?.action }}
              </v-chip>
            </td>
            
            <td class="pa-3">
              <div class="d-flex align-center">
                <span 
                  class="text-body-2 font-weight-medium mr-2"
                  :class="getRatingColor(recommendation?.rating_from || '')"
                >
                  {{ recommendation?.rating_from }}
                </span>
                <v-icon size="small" color="grey-darken-1">
                  mdi-arrow-right
                </v-icon>
                <span 
                  class="text-body-2 font-weight-medium ml-2"
                  :class="getRatingColor(recommendation?.rating_to || '')"
                >
                  {{ recommendation?.rating_to }}
                </span>
              </div>
            </td>
            
            <td class="pa-3">
              <div class="text-body-2 font-weight-medium text-grey-darken-3">
                ${{ recommendation?.target_from }} - ${{ recommendation?.target_to }}
              </div>
            </td>
            
            <td class="pa-3">
              <div class="text-body-2 text-grey-darken-1">
                {{ formatDate(recommendation?.time || '') }}
              </div>
            </td>
          </tr>
        </tbody>
      </v-table>
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
.recommendations-table {
  border-radius: 16px;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.recommendations-table-content {
  border-collapse: separate;
  border-spacing: 0;
}

.recommendation-row {
  border-bottom: 1px solid #f0f0f0;
  transition: background-color 0.2s ease;
}

.recommendation-row:last-child {
  border-bottom: none;
}

.recommendation-row:hover {
  background-color: #f8f9fa;
}

:deep(.v-table thead th) {
  border-bottom: 1px solid #e0e0e0;
  background-color: #fafafa;
  font-size: 0.875rem;
}

/* Dark theme support */
.v-theme--dark .recommendations-table {
  background: linear-gradient(135deg, #1e1e1e 0%, #2d2d2d 100%);
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.v-theme--dark .recommendation-row:hover {
  background-color: #2d2d2d;
}
</style>
