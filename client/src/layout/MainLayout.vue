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
</script>
