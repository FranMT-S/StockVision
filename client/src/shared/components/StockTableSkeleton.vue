<template>
  <v-table class="recommendations-table" density="comfortable">
    <thead>
      <tr>
        <th class="text-left font-weight-bold text-grey-darken-2">Stock</th>
        <th class="text-left font-weight-bold text-grey-darken-2">Price</th>
        <th class="text-left font-weight-bold text-grey-darken-2">Change</th>
        <th class="text-left font-weight-bold text-grey-darken-2">Sentiment</th>
        <th class="text-left font-weight-bold text-grey-darken-2">Advice</th>
        <th class="text-left font-weight-bold text-grey-darken-2">Last Rating Date</th>
      </tr>
    </thead>
    <tbody>
      <tr 
        v-for="n in rows" 
        :key="n"
        class="stock-row"
      >
        <!-- Stock Column Skeleton -->
        <td class="pa-4">
          <div class="d-flex    tw-items-center">
            <v-skeleton-loader
              type="avatar"
              class="mr-3"
              width="40"
              
            />
            <div class="flex flex-grow-1 ">
              <v-skeleton-loader
                type="text"
                class="mb-1"
                width="150"
      
              />
            </div>
          </div>
        </td>

        <!-- Price Column Skeleton -->
        <td class="text-left ">
          <div >
            <v-skeleton-loader
              type="text"
              width="120"
          
            />
          </div>
        </td>

        <!-- Change Column Skeleton -->
        <td class="text-left ">
          <div >
            <v-skeleton-loader
              type="text"
              width="120"
          
            />
          </div>
        </td>

        <!-- Sentiment Column Skeleton -->
        <td class="text-left ">
          <v-skeleton-loader
            type="avatar"
              width="40"
            class="mx-auto"
          />
        </td>

        <!-- Advice Column Skeleton -->
        <td class="text-left ">
          <v-skeleton-loader
            type="text"
              width="120"     
          />
        </td>

        <!-- Last Rating Date Column Skeleton -->
        <td class="text-left ">
          <v-skeleton-loader
            type="text"
              width="120"    
          />
        </td>
      </tr>
    </tbody>
  </v-table>
</template>

<script setup lang="ts">
interface Props {
  rows?: number
}

withDefaults(defineProps<Props>(), {
  rows: 6
})
</script>

<style scoped>
.recommendations-table {
  border-collapse: separate;
  border-spacing: 0;
}

.stock-row {
  border-bottom: 1px solid #e0e0e0;
}

.stock-row:last-child {
  border-bottom: none;
}

:deep(.v-table thead th) {
  border-bottom: 1px solid #e0e0e0;
  background-color: #fafafa;
}

:deep(.v-skeleton-loader__text) {
  border-radius: 4px;
  margin: 0px;
  width: 100%;
}

:deep(.v-skeleton-loader__avatar) {
  border-radius: 50%;
  margin: 0;
  width: 10px;
  height: 10px;
}

</style>
