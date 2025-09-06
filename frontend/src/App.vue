<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import { computed } from 'vue'
import AppHeader from './components/Layout/AppHeader.vue'

const route = useRoute()
const routeKey = computed(() => route.fullPath)
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <AppHeader />
    <main class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <RouterView :key="routeKey" />
    </main>
  </div>
</template>
