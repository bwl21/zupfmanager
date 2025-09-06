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

        <!-- Configuration Section -->
        <div>
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-medium text-gray-900">Configuration</h3>
            <button
              @click="showConfigModal = true"
              class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              {{ data.config ? 'Edit Configuration' : 'Add Configuration' }}
            </button>
          </div>
          <p class="text-sm text-gray-600">
            {{ data.config ? 'Custom Zupfnoter configuration is active.' : 'Using default Zupfnoter configuration.' }}
            Click the button above to view or modify the configuration.
          </p>
        </div>

        <!-- Project Songs Management -->
        <div>
          <ProjectSongManager :project-id="projectId" />
        </div>

        <!-- Project Build Management -->
        <div>
          <ProjectBuildManager :project-id="projectId" />
        </div>

        <!-- Actions -->
        <div>
          <h3 class="text-lg font-medium text-gray-900 mb-4">Actions</h3>
          <div class="flex flex-wrap gap-3">
            <button
              @click="openEditModal"
              class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              type="button"
            >
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
              Edit Project
            </button>
            <button
              @click="confirmDelete"
              class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
            >
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
              Remove Project
            </button>
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

    <!-- Configuration Modal -->
    <ProjectConfigModal
      v-if="showConfigModal && data"
      :project="data"
      @close="showConfigModal = false"
      @updated="handleConfigUpdated"
    />

    <!-- Edit Project Modal -->
    <EditProjectModal
      v-if="showEditModal && data"
      :project="data"
      @close="showEditModal = false"
      @updated="handleProjectUpdated"
    />

  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { useRoute, RouterLink } from 'vue-router'
import { projectApi } from '@/services/api'
import type { ProjectResponse } from '@/types/api'
import ProjectSongManager from '@/components/ProjectSongManager.vue'
import ProjectBuildManager from '@/components/ProjectBuildManager.vue'
import ProjectConfigModal from '@/components/ProjectConfigModal.vue'
import EditProjectModal from '@/components/EditProjectModal.vue'

const route = useRoute()
const projectId = parseInt(route.params.id as string)

// State
const showConfigModal = ref(false)
const showEditModal = ref(false)

// Fetch project details
const { data, isLoading, error, refetch } = useQuery({
  queryKey: ['projects', projectId],
  queryFn: () => projectApi.get(projectId),
  enabled: !!projectId
})

// Methods
const confirmDelete = () => {
  if (confirm('Are you sure you want to delete this project? This action cannot be undone.')) {
    // TODO: Implement project deletion
    console.log('Delete project:', projectId)
  }
}

const handleConfigUpdated = (updatedProject: ProjectResponse) => {
  // Update the data directly and also refetch to ensure consistency
  if (data.value) {
    data.value.config = updatedProject.config
  }
  refetch()
}

const openEditModal = () => {
  console.log('Opening edit modal')
  showEditModal.value = true
}

const handleProjectUpdated = (updatedProject: ProjectResponse) => {
  // Refetch the project data to update the UI
  refetch()
}
</script>
