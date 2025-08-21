<template>
  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="$emit('close')">
    <div class="relative top-20 mx-auto p-5 border w-11/12 max-w-md shadow-lg rounded-md bg-white" @click.stop>
      <!-- Header -->
      <div class="flex items-center justify-between pb-4 border-b border-gray-200">
        <h3 class="text-lg font-medium text-gray-900">Edit Project</h3>
        <button
          @click="$emit('close')"
          class="text-gray-400 hover:text-gray-600 transition-colors"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Form -->
      <form @submit.prevent="updateProject" class="mt-4">
        <div class="space-y-4">
          <!-- Title -->
          <div>
            <label for="title" class="block text-sm font-medium text-gray-700">
              Title
            </label>
            <input
              id="title"
              v-model="formData.title"
              type="text"
              required
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder="Enter project title"
            />
          </div>

          <!-- Short Name -->
          <div>
            <label for="shortName" class="block text-sm font-medium text-gray-700">
              Short Name
            </label>
            <input
              id="shortName"
              v-model="formData.short_name"
              type="text"
              required
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder="Enter short name (e.g., xmas)"
            />
          </div>

          <!-- ABC File Directory Preference -->
          <div>
            <label for="abcFileDir" class="block text-sm font-medium text-gray-700">
              ABC File Directory Preference
            </label>
            <input
              id="abcFileDir"
              v-model="formData.abc_file_dir_preference"
              type="text"
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder="Enter directory path (optional)"
            />
            <p class="mt-1 text-sm text-gray-500">
              Default directory for ABC files and PDF generation
            </p>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="bg-red-50 border border-red-200 rounded-md p-3">
            <div class="flex">
              <svg class="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div class="ml-3">
                <h3 class="text-sm font-medium text-red-800">Error</h3>
                <p class="mt-1 text-sm text-red-700">{{ error }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="flex items-center justify-end space-x-3 pt-4 border-t border-gray-200 mt-6">
          <button
            type="button"
            @click="$emit('close')"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            Cancel
          </button>
          <button
            type="submit"
            :disabled="isUpdating"
            class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ isUpdating ? 'Updating...' : 'Update Project' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { projectApi } from '@/services/api'
import type { ProjectResponse, UpdateProjectRequest } from '@/types/api'

interface Props {
  project: ProjectResponse
}

const props = defineProps<Props>()

const emit = defineEmits<{
  close: []
  updated: [project: ProjectResponse]
}>()

// State
const isUpdating = ref(false)
const error = ref('')

const formData = reactive<UpdateProjectRequest>({
  title: '',
  short_name: '',
  abc_file_dir_preference: ''
})

// Initialize form data
onMounted(() => {
  formData.title = props.project.title
  formData.short_name = props.project.short_name
  formData.abc_file_dir_preference = props.project.abc_file_dir_preference || ''
})

// Methods
const updateProject = async () => {
  isUpdating.value = true
  error.value = ''

  try {
    const updatedProject = await projectApi.update(props.project.id, {
      title: formData.title,
      short_name: formData.short_name,
      abc_file_dir_preference: formData.abc_file_dir_preference || undefined
    })

    emit('updated', updatedProject)
    emit('close')
  } catch (err) {
    console.error('Failed to update project:', err)
    error.value = err instanceof Error ? err.message : 'Failed to update project. Please try again.'
  } finally {
    isUpdating.value = false
  }
}
</script>
