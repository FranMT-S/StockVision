<template>
  <div class="p-6">
    <v-container>
      <v-row justify="center">
        <v-col cols="12" md="8">
          <v-card class="pa-6" elevation="2">
            <v-card-title class="text-h4 text-center mb-4">
              Welcome to StockVision
            </v-card-title>
            <v-card-text class="text-center">
              <p class="text-h6 mb-4">
                This is a test message to verify that the layout is working correctly.
              </p>
              <p class="text-body-1 mb-6">
                Your stock analysis and recommendation platform is ready to use!
              </p>
              <v-btn
                color="primary"
                size="large"
                variant="elevated"
                @click="$router.push('/tickers')"
              >
                <v-icon left>mdi-trending-up</v-icon>
                View Tickers
              </v-btn>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script setup lang="ts">
// Home page component
</script>
