<template>
  <div class="text-xs text-gray-500">
    <span v-if="versionData">
      v{{ versionData.version }}
      <span v-if="versionData.git_commit" class="ml-1">
        ({{ versionData.git_commit.substring(0, 7) }})
      </span>
    </span>
    <span v-else-if="isLoading">Loading...</span>
    <span v-else>Version unavailable</span>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface VersionData {
  version: string
  git_commit: string
  timestamp: string
}

const versionData = ref<VersionData | null>(null)
const isLoading = ref(false)

const fetchVersion = async () => {
  isLoading.value = true
  try {
    const response = await fetch('/api/version')
    if (response.ok) {
      versionData.value = await response.json()
    }
  } catch (error) {
    console.warn('Failed to fetch version info:', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchVersion()
})
</script>
