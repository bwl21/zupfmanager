<template>
  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
    <div class="relative top-20 mx-auto p-5 border w-11/12 md:w-2/3 lg:w-1/2 shadow-lg rounded-md bg-white">
      <div class="mt-3">
        <!-- Header -->
        <div class="flex items-center justify-between pb-4 border-b">
          <h3 class="text-lg font-medium text-gray-900">
            Build Status - {{ build.build_id.substring(0, 8) }}
          </h3>
          <button
            @click="$emit('close')"
            class="text-gray-400 hover:text-gray-600"
          >
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Build Info -->
        <div class="mt-4 space-y-4">
          <!-- Status -->
          <div class="flex items-center space-x-3">
            <div class="flex-shrink-0">
              <div v-if="build.status === 'completed'" class="h-4 w-4 bg-green-400 rounded-full"></div>
              <div v-else-if="build.status === 'failed'" class="h-4 w-4 bg-red-400 rounded-full"></div>
              <div v-else-if="build.status === 'running'" class="h-4 w-4 bg-blue-400 rounded-full animate-pulse"></div>
              <div v-else class="h-4 w-4 bg-gray-400 rounded-full"></div>
            </div>
            <div>
              <p class="text-sm font-medium text-gray-900">Status</p>
              <p :class="getStatusTextColor(build.status)" class="text-sm font-medium capitalize">
                {{ build.status }}
              </p>
            </div>
          </div>

          <!-- Progress (for running builds) -->
          <div v-if="buildStatus && (build.status === 'running' || build.status === 'pending')" class="space-y-2">
            <div class="flex justify-between items-center">
              <p class="text-sm font-medium text-gray-900">Progress</p>
              <p class="text-sm text-gray-600">{{ buildStatus.progress }}%</p>
            </div>
            <div class="bg-gray-200 rounded-full h-3">
              <div 
                class="bg-blue-600 h-3 rounded-full transition-all duration-300"
                :style="{ width: `${buildStatus.progress}%` }"
              ></div>
            </div>
            <p v-if="buildStatus.message" class="text-sm text-gray-600">{{ buildStatus.message }}</p>
          </div>

          <!-- Timing -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p class="text-sm font-medium text-gray-900">Started</p>
              <p class="text-sm text-gray-600">{{ formatDate(build.started_at) }}</p>
            </div>
            <div v-if="build.completed_at">
              <p class="text-sm font-medium text-gray-900">Completed</p>
              <p class="text-sm text-gray-600">{{ formatDate(build.completed_at) }}</p>
            </div>
          </div>

          <!-- Duration -->
          <div v-if="build.completed_at">
            <p class="text-sm font-medium text-gray-900">Duration</p>
            <p class="text-sm text-gray-600">{{ calculateDuration(build.started_at, build.completed_at) }}</p>
          </div>

          <!-- Output Directory -->
          <div>
            <p class="text-sm font-medium text-gray-900">Output Directory</p>
            <p class="text-sm text-gray-600 font-mono">{{ build.output_dir }}</p>
          </div>

          <!-- Generated Files -->
          <div v-if="build.generated_files && build.generated_files.length > 0">
            <p class="text-sm font-medium text-gray-900 mb-2">Generated Files</p>
            <div class="bg-gray-50 rounded-md p-3">
              <ul class="space-y-1">
                <li v-for="file in build.generated_files" :key="file" 
                    class="text-sm text-gray-700 font-mono flex items-center">
                  <svg class="h-4 w-4 text-gray-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  {{ file }}
                </li>
              </ul>
            </div>
          </div>

          <!-- Error Message -->
          <div v-if="build.error || buildStatus?.error" class="bg-red-50 border border-red-200 rounded-md p-4">
            <div class="flex">
              <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div class="ml-3">
                <h4 class="text-sm font-medium text-red-800">Build Error</h4>
                <p class="mt-1 text-sm text-red-700">{{ build.error || buildStatus?.error }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="mt-6 flex justify-between">
          <div>
            <button
              v-if="build.status === 'running' || build.status === 'pending'"
              @click="refreshStatus"
              :disabled="isRefreshing"
              class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
            >
              <svg class="h-4 w-4 mr-2" :class="{ 'animate-spin': isRefreshing }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              {{ isRefreshing ? 'Refreshing...' : 'Refresh Status' }}
            </button>
          </div>
          
          <div class="flex gap-4">
            <button
              v-if="build.status === 'completed'"
              class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
              disabled
            >
              Download Files (Coming Soon)
            </button>
            <button
              @click="$emit('close')"
              class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { projectBuildApi } from '@/services/api'
import type { BuildResultResponse, BuildStatusResponse } from '@/types/api'

interface Props {
  build: BuildResultResponse
  projectId: number
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
}>()

// State
const buildStatus = ref<BuildStatusResponse | null>(null)
const isRefreshing = ref(false)
let statusInterval: number | null = null

// Methods
const refreshStatus = async () => {
  if (props.build.status === 'completed' || props.build.status === 'failed') {
    return
  }

  isRefreshing.value = true
  try {
    const status = await projectBuildApi.getStatus(props.projectId, props.build.build_id)
    buildStatus.value = status
  } catch (err) {
    console.error('Failed to refresh build status:', err)
  } finally {
    isRefreshing.value = false
  }
}

const getStatusTextColor = (status: string) => {
  const colors = {
    pending: 'text-yellow-600',
    running: 'text-blue-600',
    completed: 'text-green-600',
    failed: 'text-red-600'
  }
  return colors[status as keyof typeof colors] || 'text-gray-600'
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

const calculateDuration = (start: string, end: string) => {
  const startTime = new Date(start).getTime()
  const endTime = new Date(end).getTime()
  const duration = Math.round((endTime - startTime) / 1000)
  
  if (duration < 60) {
    return `${duration} seconds`
  } else if (duration < 3600) {
    const minutes = Math.round(duration / 60)
    return `${minutes} minute${minutes !== 1 ? 's' : ''}`
  } else {
    const hours = Math.round(duration / 3600)
    return `${hours} hour${hours !== 1 ? 's' : ''}`
  }
}

// Lifecycle
onMounted(() => {
  // Initial status fetch
  refreshStatus()
  
  // Set up periodic refresh for active builds
  if (props.build.status === 'running' || props.build.status === 'pending') {
    statusInterval = setInterval(refreshStatus, 3000) // Check every 3 seconds
  }
})

onUnmounted(() => {
  if (statusInterval) {
    clearInterval(statusInterval)
  }
})
</script>
