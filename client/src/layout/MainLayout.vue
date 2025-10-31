<template>
  <v-app>
    <NavBar />
    <!-- Main Content Area -->
    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import NavBar from '@/shared/components/NavBar.vue'
import { useOnboardingStore } from '@/shared/store/onboardingStore'
import { onMounted } from 'vue'

const onboardingStore = useOnboardingStore()

onMounted(async () => {
  await onboardingStore.fetchOnboarding()
})

</script>
