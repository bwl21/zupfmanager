<template>
  <div class="space-y-6">
    <!-- Back Button -->
    <div>
      <button
        @click="$router.back()"
        class="flex items-center text-gray-600 hover:text-gray-900 transition-colors"
      >
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        Back to Projects
      </button>
    </div>

    <!-- Project Details -->
    <div v-if="data" class="bg-white rounded-lg shadow">
      <div class="px-6 py-4 border-b border-gray-200">
        <div class="flex items-center justify-between">
          <h1 class="text-2xl font-bold text-gray-900">{{ data.title }}</h1>
          <div class="flex items-center space-x-2">
            <span class="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm font-medium">
              {{ data.short_name }}
            </span>
          </div>
        </div>
      </div>

      <div class="p-6 space-y-6">
        <!-- Basic Information -->
        <div>
          <h3 class="text-lg font-medium text-gray-900 mb-4">Project Information</h3>
          <dl class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <dt class="text-sm font-medium text-gray-500">Project ID</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ data.id }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Short Name</dt>
              <dd class="mt-1 text-sm text-gray-900 font-mono">{{ data.short_name }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Title</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ data.title }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Configuration</dt>
              <dd class="mt-1 text-sm text-gray-900">
                {{ data.config ? 'Custom configuration' : 'Default configuration' }}
              </dd>
            </div>
          </dl>
        </div>

        <!-- Configuration Preview -->
        <div v-if="data.config">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Configuration</h3>
          <div class="bg-gray-50 rounded-lg p-4">
            <pre class="text-xs text-gray-700 overflow-x-auto">{{ JSON.stringify(data.config, null, 2) }}</pre>
          </div>
        </div>

        <!-- Actions -->
        <div>
          <h3 class="text-lg font-medium text-gray-900 mb-4">Actions</h3>
          <div class="flex flex-wrap gap-3">
            <button
              class="flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
              disabled
            >
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              Add Songs (Coming Soon)
            </button>
            <button
              class="flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
              disabled
            >
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
              Build Project (Coming Soon)
            </button>
            <RouterLink
              :to="`/projects`"
              class="flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm font-medium transition-colors"
            >
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
              Edit Project
            </RouterLink>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <p class="mt-2 text-gray-600">Loading project details...</p>
    </div>

    <!-- Error State -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4">
      <div class="flex">
        <svg class="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Error loading project</h3>
          <p class="mt-2 text-sm text-red-700">{{ error.message }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { useRoute, RouterLink } from 'vue-router'
import { projectApi } from '@/services/api'

const route = useRoute()
const projectId = parseInt(route.params.id as string)

// Fetch project details
const { data, isLoading, error } = useQuery({
  queryKey: ['projects', projectId],
  queryFn: () => projectApi.get(projectId),
  enabled: !!projectId
})
</script>
