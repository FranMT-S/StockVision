<template>
  <div class="paginator-wrapper">
    <v-card 
      class="paginator-card " 
      elevation="0"
      :class="{ 'sticky bottom-0': sticky}"
    >
      <v-card-text class="pa-3">
        <v-row align="center" justify="center" no-gutters>          
          <!-- Pagination controls -->
          <v-col cols="12" sm="12" class="d-flex justify-center mt-2 mt-sm-0">
            <v-pagination
              v-model="currentPage"
              :length="totalPages"
              :total-visible="visiblePages"
              color="primary"
              density="compact"
              @update:model-value="handlePageChange"
            />
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  initialPage: number
  totalPages: number
  totalItems: number
  itemsPerPage: number
  visiblePages: number
  sticky?: boolean
}

interface Emits {
  (e: 'update:currentPage', page: number): void
}

const props = withDefaults(defineProps<Props>(), {
  sticky: false,
  visiblePages: 6
})

const currentPage = ref(props.initialPage || 1)

const emit = defineEmits<Emits>()


// Handle page changes
const handlePageChange = (page: number) => {
  emit('update:currentPage', page)
}
</script>

<style scoped>
.paginator-wrapper {
  width: 100%;
}

.paginator-card {
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.paginator-card.sticky {
  position: fixed;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 1000;
  max-width: calc(100vw - 32px);
  width: 100%;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
}

/* Dark theme support */
.v-theme--dark .paginator-card {
  background: rgba(33, 33, 33, 0.95);
  border: 1px solid rgba(255, 255, 255, 0.05);
}

/* Responsive adjustments */
@media (max-width: 600px) {
  .paginator-card.sticky {
    bottom: 8px;
    left: 8px;
    right: 8px;
    transform: none;
    max-width: none;
    width: auto;
  }
  
  .paginator-card .v-card-text {
    padding: 12px !important;
  }
}

/* Animation for sticky appearance */
.paginator-card.sticky {
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from {
    transform: translateX(-50%) translateY(100%);
    opacity: 0;
  }
  to {
    transform: translateX(-50%) translateY(0);
    opacity: 1;
  }
}
</style>
