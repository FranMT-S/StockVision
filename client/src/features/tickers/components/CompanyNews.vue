<template>
  <section class="company-news-sidebar">
    <v-card-title class="pa-4">
      <h3 class="text-h6 font-weight-medium dark:tw-text-[#ffffff] company-news-sidebar-title">
        Company News
      </h3>
    </v-card-title>

    <!-- News Cards -->
    <div class="news-container lg pa-2">
      <v-card
        tag="article"
        v-for="news in companyNews"
        :key="news.id"
        class="news-card mb-3"
        elevation="2"
        :href="news.url"
        target="_blank"
        rel="noopener noreferrer"
      >
        <!-- News Image -->
        <v-img
          v-if="news.image"
          :src="news.image"
          height="140"
          cover
          class="news-image"
        >
        </v-img>
        
        <!-- Fallback for no image -->
        <div v-else class="news-placeholder">
          <v-icon size="48" color="grey-lighten-1">mdi-newspaper-variant-outline</v-icon>

        </div>

        <v-card-text class="pa-3">
          <!-- Headline -->
          <h4 class="text-subtitle-2 font-weight-bold mb-2 news-headline">
            {{ news.headline }}
          </h4>

          <!-- Summary -->
          <p class="text-caption text-grey-darken-1 mb-2 news-summary">
            {{ truncateSummary(news.summary) }}
          </p>

          <!-- Meta Information -->
          <div class="d-flex align-center justify-space-between tw-flex-wrap tw-gap-2">
            <div class="d-flex align-center tw-gap-2">
              <v-chip
                size="x-small"
                variant="text"
                color="secondary"
                class="pa-0"
              >
                <v-icon size="12" class="mr-1">mdi-calendar</v-icon>
                {{ new Date(news.datetime).toLocaleDateString() }}
              </v-chip>
            </div>
            
            <v-chip
              size="x-small"
              variant="outlined"
              color="primary"
              class="source-chip"
            >
              {{ news.source }}
            </v-chip>
          </div>
        </v-card-text>

        <!-- Read More Indicator -->
        <v-card-actions class="pa-2 pt-0">
          <v-spacer />
          <v-btn
            size="x-small"
            variant="text"
            color="primary"
            append-icon="mdi-open-in-new"
          >
            Read More
          </v-btn>
        </v-card-actions>
      </v-card>

      <!-- Empty State -->
      <div v-if="!companyNews || companyNews.length === 0" class="text-center pa-8">
        <v-icon size="64" color="grey-lighten-2" class="mb-4">
          mdi-newspaper-variant-outline
        </v-icon>
        <p class="text-body-2 text-grey-darken-1">No news available</p>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { CompanyNew } from '@/shared/models/recomendations'

interface Props {
  companyNews: CompanyNew[]
}

const props = defineProps<Props>()

const truncateSummary = (summary: string, maxLength: number = 120): string => {
  if (summary.length <= maxLength) return summary
  return summary.substring(0, maxLength).trim() + '...'
}
</script>

<style scoped>
.company-news-sidebar {
  height: 100%;
  overflow-y: auto;
}


.news-container::-webkit-scrollbar {
  width: 6px;
}

.news-container::-webkit-scrollbar-track {
  background: transparent;
}

.news-container::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 3px;
}

.news-container::-webkit-scrollbar-thumb:hover {
  background: #a0aec0;
}

.news-card {
  border-radius: 12px;
  transition: all 0.3s ease;
  cursor: pointer;
  text-decoration: none;
  border: 1px solid rgba(0, 0, 0, 0.08);
}

.news-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.12) !important;
}

.news-image {
  border-radius: 12px 12px 0 0;
}

.news-overlay {
  background: linear-gradient(to bottom, rgba(0, 0, 0, 0.3), transparent);
  height: 100%;
  display: flex;
  align-items: flex-start;
  justify-content: flex-start;
}

.news-placeholder {
  height: 140px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e9f2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  border-radius: 12px 12px 0 0;
}

.news-placeholder .category-chip {
  position: absolute;
  top: 8px;
  left: 8px;
}

.news-headline {
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  color: #2d3748;
}

.news-summary {
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.source-chip {
  font-size: 10px;
  height: 20px;
}

/* Dark theme support */
.v-theme--dark .news-card {
  background: linear-gradient(135deg, #1e1e1e 0%, #2d2d2d 100%);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.v-theme--dark .news-headline {
  color: #ffffff;
}

.v-theme--dark .news-placeholder {
  background: linear-gradient(135deg, #2d2d2d 0%, #3d3d3d 100%);
}
</style>
