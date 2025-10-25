<template>
  <v-app>
    <!-- Navigation Bar -->
    <v-app-bar
      color="primary"
      dark
      elevation="2"
      app
    >
      <v-container fluid class="pa-0">
        <v-row align="center" no-gutters>
          <!-- Logo Section -->
          <v-col cols="auto">
            <v-btn
              variant="text"
              class="text-h6 font-weight-bold"
              color="white"
              @click="$router.push('/')"
            >
              <v-icon left color="blue-lighten-2" size="large">
                mdi-trending-up
              </v-icon>
              StockVision
            </v-btn>
          </v-col>
          
          <!-- Search Bar -->
          <v-col cols="12" sm="8" md="6" lg="4" class="ml-4">
           <v-text-field
            v-model="searchQuery"
            placeholder="Search by id or ticker..."
            variant="outlined"
            density="compact"
            hide-details
            clearable
            class="search-field"
            bg-color="white"
            color="primary"
            @keyup.enter="handleSearch"
          >
            <template #prepend-inner>
              <v-icon
                class="cursor-pointer"
                @click="handleSearch"
              >
                mdi-magnify
              </v-icon>
            </template>
              <template #append-inner>
                <v-progress-circular
                  v-if="searchLoading"
                  indeterminate
                  size="16"
                  width="2"
                  color="primary"
                />
              </template>
            </v-text-field>
          </v-col>
          
          <!-- Spacer and additional actions -->
          <v-spacer />
          
          <!-- Navigation Menu -->
          <v-col cols="auto" class="d-none d-md-flex">
            <v-btn
              variant="text"
              color="white"
              class="mr-2"
              @click="$router.push('/')"
            >
              Home
            </v-btn>
            <v-btn
              variant="text"
              color="white"
              class="mr-2"
              @click="$router.push('/tickers')"
            >
              Tickers
            </v-btn>
          </v-col>
          
          <!-- Theme toggle and other actions can go here -->
          <v-col cols="auto">
            <ToggleTheme />
          </v-col>
        </v-row>
      </v-container>
    </v-app-bar>

    <!-- Main Content Area -->
    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import ToggleTheme from '@shared/components/ToggleTheme.vue'
import { useRouter } from 'vue-router'
import { useRoute } from 'vuetify/lib/composables/router.mjs'


const router = useRouter()
const route = useRoute()
const searchQuery = ref(route.value?.query?.q?.toString() || '')
const searchLoading = ref(false)

watch(() => route.value?.query?.q, (newQuery) => {
  searchQuery.value = newQuery?.toString() || ''
})

// Handle search functionality
const handleSearch =  () => {
  if (!searchQuery.value.trim()) return
  
  searchLoading.value = true
  
  try {
    router.push({
      name: 'Tickers',
      query: {
        q: searchQuery.value
      }
    })

  } catch (error) {
    console.error('Search error:', error)
  } finally {
    searchLoading.value = false
  }
}
</script>

<style scoped>
.search-field {
  max-width: 400px;
}

:deep(.v-field__input) {
  color: #1976d2 !important;
}

:deep(.v-field__outline) {
  color: rgba(255, 255, 255, 0.7) !important;
}
</style>
