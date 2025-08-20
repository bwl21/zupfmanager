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
        Back to Songs
      </button>
    </div>

    <!-- Song Details -->
    <div v-if="data" class="bg-white rounded-lg shadow">
      <div class="px-6 py-4 border-b border-gray-200">
        <div class="flex items-center justify-between">
          <h1 class="text-2xl font-bold text-gray-900">{{ data.title }}</h1>
          <div class="flex items-center space-x-2">
            <span class="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm font-medium">
              ID: {{ data.id }}
            </span>
          </div>
        </div>
      </div>

      <div class="p-6 space-y-6">
        <!-- Basic Information -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <h3 class="text-lg font-medium text-gray-900 mb-4">Basic Information</h3>
            <dl class="space-y-3">
              <div>
                <dt class="text-sm font-medium text-gray-500">Filename</dt>
                <dd class="mt-1 text-sm text-gray-900 font-mono bg-gray-50 px-2 py-1 rounded">
                  {{ data.filename }}
                </dd>
              </div>
              <div v-if="data.genre">
                <dt class="text-sm font-medium text-gray-500">Genre</dt>
                <dd class="mt-1">
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                    {{ data.genre }}
                  </span>
                </dd>
              </div>
              <div v-if="data.copyright">
                <dt class="text-sm font-medium text-gray-500">Copyright</dt>
                <dd class="mt-1 text-sm text-gray-900">{{ data.copyright }}</dd>
              </div>
              <div v-if="data.tocinfo">
                <dt class="text-sm font-medium text-gray-500">Table of Contents Info</dt>
                <dd class="mt-1 text-sm text-gray-900">{{ data.tocinfo }}</dd>
              </div>
            </dl>
          </div>

          <!-- Actions -->
          <div>
            <h3 class="text-lg font-medium text-gray-900 mb-4">Actions</h3>
            <div class="space-y-3">
              <button
                @click="showPreviewModal = true"
                class="w-full flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
              >
                <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                Preview PDFs
              </button>
              <button
                class="w-full flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
                disabled
              >
                <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                Download ABC (Coming Soon)
              </button>
              <button
                @click="showAddToProjectModal = true"
                class="w-full flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
              >
                <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                </svg>
                Add to Project
              </button>
            </div>
          </div>
        </div>

        <!-- Metadata -->
        <div>
          <h3 class="text-lg font-medium text-gray-900 mb-4">Metadata</h3>
          <div class="bg-gray-50 rounded-lg p-4">
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
              <div>
                <span class="font-medium text-gray-500">Song ID:</span>
                <span class="ml-2 text-gray-900">{{ data.id }}</span>
              </div>
              <div>
                <span class="font-medium text-gray-500">File Type:</span>
                <span class="ml-2 text-gray-900">ABC Notation</span>
              </div>
              <div>
                <span class="font-medium text-gray-500">Status:</span>
                <span class="ml-2 text-green-600">Available</span>
              </div>
              <div>
                <span class="font-medium text-gray-500">Format:</span>
                <span class="ml-2 text-gray-900">.abc</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Related Information -->
        <div>
          <h3 class="text-lg font-medium text-gray-900 mb-4">Related</h3>
          <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <div class="flex">
              <svg class="w-5 h-5 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div class="ml-3">
                <h4 class="text-sm font-medium text-blue-800">ABC Notation</h4>
                <p class="mt-1 text-sm text-blue-700">
                  This song is stored in ABC notation format, a text-based music notation system.
                  You can use it with Zupfnoter to generate sheet music for zither instruments.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <p class="mt-2 text-gray-600">Loading song details...</p>
    </div>

    <!-- Error State -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4">
      <div class="flex">
        <svg class="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Error loading song</h3>
          <p class="mt-2 text-sm text-red-700">{{ error.message }}</p>
        </div>
      </div>
    </div>

    <!-- Add to Project Modal -->
    <AddToProjectModal
      v-if="showAddToProjectModal && data"
      :song="data"
      @close="showAddToProjectModal = false"
      @added="handleAddedToProject"
    />

    <!-- Preview Modal -->
    <PreviewModal
      v-if="showPreviewModal && data"
      :song="data"
      @close="showPreviewModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { useRoute } from 'vue-router'
import { songApi } from '@/services/api'
import AddToProjectModal from '@/components/AddToProjectModal.vue'
import PreviewModal from '@/components/PreviewModal.vue'

const route = useRoute()
const songId = parseInt(route.params.id as string)

// State
const showAddToProjectModal = ref(false)
const showPreviewModal = ref(false)

// Fetch song details
const { data, isLoading, error } = useQuery({
  queryKey: ['songs', songId],
  queryFn: () => songApi.get(songId),
  enabled: !!songId
})

// Methods
const handleAddedToProject = () => {
  showAddToProjectModal.value = false
  // Could show a success message here
}
</script>
