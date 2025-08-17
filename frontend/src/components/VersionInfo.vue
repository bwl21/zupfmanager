<template>
  <div class="text-xs text-gray-500">
    <span v-if="versionData">
      v{{ versionData.version }}
      <span v-if="versionData.git_commit" class="ml-1">
        ({{ versionData.git_commit.substring(0, 7) }})
      </span>
    </span>
    <span v-else-if="isLoading">Loading version...</span>
    <span v-else-if="error">{{ error }}</span>
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
const isLoading = ref(true)
const error = ref<string | null>(null)

const fetchVersion = async () => {
  isLoading.value = true
  error.value = null
  try {
    console.log('Fetching version from /api/version')
    const response = await fetch('/api/version')
    console.log('Version response status:', response.status)
    if (response.ok) {
      const data = await response.json()
      console.log('Version data:', data)
      versionData.value = data
    } else {
      error.value = `HTTP ${response.status}`
    }
  } catch (err) {
    console.error('Failed to fetch version info:', err)
    error.value = 'Network error'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchVersion()
})
</script>
