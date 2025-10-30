  <template>

  <v-app-bar
    color="primary"
    dark
    elevation="2"
    app
  >
    <v-container fluid class="pa-0">
      <v-row align="center" no-gutters>
        <!-- Logo Section -->
        <v-col cols="auto" 
        >
          <v-btn
            variant="text"
            class="text-h6 font-weight-bold !tw-hidden sm:!tw-block"
            color="white"
            @click="$router.push('/')"
          >
            <v-icon left color="blue-lighten-2" size="large">
              mdi-trending-up
            </v-icon>
            StockVision
          </v-btn>
        </v-col>
        
        <!-- Search Bar  cols="9" sm="4" md="6" lg="4"-->
        <v-col  :cols="screenWidth < 480 ? 11 : 6" >
          <v-text-field
          v-model="searchQuery"
          placeholder="Search by ticker or company..."
          variant="outlined"
          density="compact"
          hide-details
          clearable
          class="search-field  ml-4 "
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
      </v-row>
    </v-container>
  </v-app-bar>
 </template>
  

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useRoute } from 'vuetify/lib/composables/router.mjs'
import { RouterNames } from '@/router/names'
import { useBreakpoints } from '../composables/useBreakpoints'

const router = useRouter()
const route = useRoute()
const searchQuery = ref(route.value?.query?.q?.toString() || '')
const searchLoading = ref(false)
const {screenWidth} = useBreakpoints();

watch(() => route.value?.query?.q, (newQuery) => {
  searchQuery.value = newQuery?.toString() || ''
})

// Handle search functionality
const handleSearch =  () => {
  if (!searchQuery.value.trim()) return
  
  searchLoading.value = true
  
  try {
    router.push({
      name: RouterNames.Tickers,
      query: {
        q: searchQuery.value,
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
