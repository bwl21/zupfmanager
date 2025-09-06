<template>
  <div class="space-y-6">
    <!-- Header with Build Button -->
    <div class="flex justify-between items-center">
      <h3 class="text-lg font-medium text-gray-900">Project Build</h3>
      <button
        @click="showBuildModal = true"
        :disabled="isBuilding"
        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <svg class="-ml-1 mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012 2v2M7 7h10" />
        </svg>
        <span v-if="isBuilding">Building...</span>
        <span v-else>Start Build</span>
      </button>
    </div>

    <!-- Current Build Status -->
    <div v-if="currentBuild && (currentBuild.status === 'pending' || currentBuild.status === 'running')" 
         class="bg-blue-50 border border-blue-200 rounded-lg p-4">
      <div class="flex items-center">
        <div class="flex-shrink-0">
          <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-blue-600"></div>
        </div>
        <div class="ml-3 flex-1">
          <h4 class="text-sm font-medium text-blue-800">Build in Progress</h4>
          <p class="text-sm text-blue-700">{{ currentBuildStatus?.message || 'Building project...' }}</p>
          
          <!-- Progress Bar -->
          <div class="mt-2">
            <div class="bg-blue-200 rounded-full h-2">
              <div 
                class="bg-blue-600 h-2 rounded-full transition-all duration-300"
                :style="{ width: `${currentBuildStatus?.progress || 0}%` }"
              ></div>
            </div>
            <p class="text-xs text-blue-600 mt-1">{{ currentBuildStatus?.progress || 0 }}% complete</p>
          </div>
        </div>
        <button
          @click="refreshBuildStatus"
          class="ml-3 text-blue-600 hover:text-blue-800"
        >
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Build History -->
    <div>
      <div class="flex justify-between items-center mb-4">
        <h4 class="text-md font-medium text-gray-900">Build History</h4>
        <div class="flex gap-2">
          <button
            @click="clearBuildHistory"
            :disabled="isClearingHistory || builds.length === 0"
            class="text-sm text-red-600 hover:text-red-800 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <svg class="inline h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            {{ isClearingHistory ? 'Clearing...' : 'Clear History' }}
          </button>
          <button
            @click="loadBuilds"
            class="text-sm text-gray-600 hover:text-gray-900"
          >
            <svg class="inline h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Refresh
          </button>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoadingBuilds" class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-600"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="buildsError" class="bg-red-50 border border-red-200 rounded-md p-4">
        <div class="flex">
          <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error loading builds</h3>
            <p class="mt-1 text-sm text-red-700">{{ buildsError }}</p>
          </div>
        </div>
      </div>

      <!-- Builds List -->
      <div v-else-if="builds.length > 0" class="bg-white shadow overflow-hidden sm:rounded-md">
        <ul class="divide-y divide-gray-200">
          <li v-for="build in builds" :key="build.build_id" class="px-6 py-4">
            <div class="flex items-center justify-between">
              <div class="flex-1 min-w-0">
                <div class="flex items-center space-x-3">
                  <!-- Status Icon -->
                  <div class="flex-shrink-0">
                    <div v-if="build.status === 'completed'" class="h-3 w-3 bg-green-400 rounded-full"></div>
                    <div v-else-if="build.status === 'failed'" class="h-3 w-3 bg-red-400 rounded-full"></div>
                    <div v-else-if="build.status === 'running'" class="h-3 w-3 bg-blue-400 rounded-full animate-pulse"></div>
                    <div v-else class="h-3 w-3 bg-gray-400 rounded-full"></div>
                  </div>
                  
                  <div class="flex-1">
                    <div class="flex items-center space-x-2">
                      <p class="text-sm font-medium text-gray-900">
                        Build {{ build.build_id.substring(0, 8) }}
                      </p>
                      <span :class="getStatusColor(build.status)" 
                            class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
                        {{ build.status }}
                      </span>
                    </div>
                    
                    <div class="mt-1 flex items-center space-x-4 text-sm text-gray-500">
                      <span>{{ formatDate(build.started_at) }}</span>
                      <span v-if="build.completed_at">
                        Duration: {{ calculateDuration(build.started_at, build.completed_at) }}
                      </span>
                      <span>Output: {{ build.output_dir }}</span>
                    </div>
                    
                    <!-- Generated Files -->
                    <div v-if="build.generated_files && build.generated_files.length > 0" class="mt-2">
                      <p class="text-xs text-gray-500">Generated files:</p>
                      <div class="flex flex-wrap gap-1 mt-1">
                        <span v-for="file in build.generated_files" :key="file"
                              class="inline-flex items-center px-2 py-1 rounded text-xs bg-gray-100 text-gray-700">
                          {{ file }}
                        </span>
                      </div>
                    </div>
                    
                    <!-- Error Message -->
                    <div v-if="build.error" class="mt-2 p-2 bg-red-50 rounded text-xs text-red-700">
                      {{ build.error }}
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- Actions -->
              <div class="flex items-center space-x-2">
                <button
                  v-if="build.status === 'running' || build.status === 'pending'"
                  @click="viewBuildStatus(build)"
                  class="text-blue-600 hover:text-blue-900 text-sm font-medium"
                >
                  View Status
                </button>
                <button
                  v-if="build.status === 'completed'"
                  class="text-green-600 hover:text-green-900 text-sm font-medium"
                  disabled
                >
                  Download (Coming Soon)
                </button>
              </div>
            </div>
          </li>
        </ul>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012 2v2M7 7h10" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900">No builds yet</h3>
        <p class="mt-1 text-sm text-gray-500">Start your first project build to generate output files.</p>
        <div class="mt-6">
          <button
            @click="showBuildModal = true"
            class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
          >
            <svg class="-ml-1 mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012 2v2M7 7h10" />
            </svg>
            Start Build
          </button>
        </div>
      </div>
    </div>

    <!-- Build Configuration Modal -->
    <BuildConfigModal
      v-if="showBuildModal"
      :project-id="projectId"
      @close="showBuildModal = false"
      @build-started="handleBuildStarted"
    />

    <!-- Build Status Modal -->
    <BuildStatusModal
      v-if="showStatusModal && selectedBuild"
      :build="selectedBuild"
      :project-id="projectId"
      @close="showStatusModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { projectBuildApi } from '@/services/api'
import type { BuildResultResponse, BuildStatusResponse } from '@/types/api'
import BuildConfigModal from './BuildConfigModal.vue'
import BuildStatusModal from './BuildStatusModal.vue'

interface Props {
  projectId: number
}

const props = defineProps<Props>()

// State
const builds = ref<BuildResultResponse[]>([])
const isLoadingBuilds = ref(false)
const buildsError = ref<string | null>(null)
const isClearingHistory = ref(false)
const showBuildModal = ref(false)
const showStatusModal = ref(false)
const selectedBuild = ref<BuildResultResponse | null>(null)
const currentBuildStatus = ref<BuildStatusResponse | null>(null)

// Computed
const currentBuild = computed(() => {
  return builds.value.find(build => 
    build.status === 'pending' || build.status === 'running'
  )
})

const isBuilding = computed(() => {
  return !!currentBuild.value
})

// Methods
const loadBuilds = async () => {
  isLoadingBuilds.value = true
  buildsError.value = null
  
  try {
    const response = await projectBuildApi.listBuilds(props.projectId)
    builds.value = response.builds.sort((a, b) => 
      new Date(b.started_at).getTime() - new Date(a.started_at).getTime()
    )
  } catch (err) {
    buildsError.value = err instanceof Error ? err.message : 'Failed to load builds'
  } finally {
    isLoadingBuilds.value = false
  }
}

const refreshBuildStatus = async () => {
  if (!currentBuild.value) return
  
  try {
    const status = await projectBuildApi.getStatus(props.projectId, currentBuild.value.build_id)
    currentBuildStatus.value = status
    
    // If build is completed or failed, refresh the builds list
    if (status.status === 'completed' || status.status === 'failed') {
      await loadBuilds()
    }
  } catch (err) {
    console.error('Failed to refresh build status:', err)
  }
}

const viewBuildStatus = (build: BuildResultResponse) => {
  selectedBuild.value = build
  showStatusModal.value = true
}

const clearBuildHistory = async () => {
  if (!confirm('Are you sure you want to clear all build history? This action cannot be undone.')) {
    return
  }
  
  isClearingHistory.value = true
  
  try {
    await projectBuildApi.clearHistory(props.projectId)
    builds.value = []
    buildsError.value = null
  } catch (err) {
    buildsError.value = err instanceof Error ? err.message : 'Failed to clear build history'
  } finally {
    isClearingHistory.value = false
  }
}

const handleBuildStarted = (build: BuildResultResponse) => {
  showBuildModal.value = false
  builds.value.unshift(build) // Add to beginning of list
  
  // Start polling for status updates
  if (build.status === 'pending' || build.status === 'running') {
    setTimeout(refreshBuildStatus, 2000)
  }
}

const getStatusColor = (status: string) => {
  const colors = {
    pending: 'bg-yellow-100 text-yellow-800',
    running: 'bg-blue-100 text-blue-800',
    completed: 'bg-green-100 text-green-800',
    failed: 'bg-red-100 text-red-800'
  }
  return colors[status as keyof typeof colors] || 'bg-gray-100 text-gray-800'
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

const calculateDuration = (start: string, end: string) => {
  const startTime = new Date(start).getTime()
  const endTime = new Date(end).getTime()
  const duration = Math.round((endTime - startTime) / 1000)
  
  if (duration < 60) {
    return `${duration}s`
  } else if (duration < 3600) {
    return `${Math.round(duration / 60)}m`
  } else {
    return `${Math.round(duration / 3600)}h`
  }
}

// Lifecycle
onMounted(() => {
  loadBuilds()
  
  // Set up periodic refresh for active builds
  const interval = setInterval(() => {
    if (currentBuild.value) {
      refreshBuildStatus()
    }
  }, 5000) // Check every 5 seconds
  
  // Cleanup interval on unmount
  return () => clearInterval(interval)
})
</script>
