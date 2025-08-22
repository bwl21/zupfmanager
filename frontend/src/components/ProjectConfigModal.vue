<template>
  <div class="fixed inset-0 bg-gray-600 bg-opacity-75 flex items-center justify-center z-50 p-4" @click="$emit('close')">
    <div class="bg-white rounded-lg shadow-xl w-4/5 h-4/5 max-w-6xl max-h-screen flex flex-col" @click.stop>
      <!-- Header -->
      <div class="flex items-center justify-between p-6 border-b border-gray-200 flex-shrink-0">
        <h3 class="text-xl font-semibold text-gray-900">Project Configuration</h3>
        <button
          @click="$emit('close')"
          class="text-gray-400 hover:text-gray-600 transition-colors p-1"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Content -->
      <div class="flex-1 overflow-hidden flex flex-col p-6">
        <div class="mb-4 flex-shrink-0">
          <p class="text-sm text-gray-600">
            This configuration is used by Zupfnoter for rendering ABC notation. 
            Modify with caution as incorrect values may cause rendering issues.
          </p>
        </div>

        <!-- Configuration Editor -->
        <div class="flex-1 flex flex-col space-y-4">
          <div class="flex-1 flex flex-col">
            <label class="block text-sm font-medium text-gray-700 mb-2 flex-shrink-0">
              Configuration JSON
            </label>
            <textarea
              v-model="configText"
              class="flex-1 w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 font-mono text-sm resize-none"
              placeholder="Enter JSON configuration..."
            />
          </div>

          <!-- Validation Error -->
          <div v-if="validationError" class="bg-red-50 border border-red-200 rounded-md p-3 flex-shrink-0">
            <div class="flex">
              <svg class="w-5 h-5 text-red-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div class="ml-3">
                <h3 class="text-sm font-medium text-red-800">Invalid JSON</h3>
                <p class="mt-1 text-sm text-red-700">{{ validationError }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="flex items-center justify-end gap-4 p-6 border-t border-gray-200 flex-shrink-0">
        <button
          @click="resetToDefault"
          class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 mr-3"
        >
          Reset to Default
        </button>
        <button
          @click="$emit('close')"
          class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 mr-3"
        >
          Cancel
        </button>
        <button
          @click="saveConfiguration"
          :disabled="!!validationError || isSaving"
          class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ isSaving ? 'Saving...' : 'Save Configuration' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { projectApi } from '@/services/api'
import type { ProjectResponse } from '@/types/api'

interface Props {
  project: ProjectResponse
}

const props = defineProps<Props>()

const emit = defineEmits<{
  close: []
  updated: [project: ProjectResponse]
}>()

// State
const configText = ref('')
const validationError = ref('')
const isSaving = ref(false)

// Initialize configuration text
onMounted(() => {
  if (props.project.config) {
    configText.value = JSON.stringify(props.project.config, null, 2)
  } else {
    configText.value = JSON.stringify({}, null, 2)
  }
})

// Validate JSON on change
watch(configText, (newValue) => {
  try {
    JSON.parse(newValue)
    validationError.value = ''
  } catch (error) {
    validationError.value = error instanceof Error ? error.message : 'Invalid JSON format'
  }
})

// Methods
const resetToDefault = async () => {
  try {
    // Load default configuration from API
    const defaultConfig = await projectApi.getDefaultConfig()
    configText.value = JSON.stringify(defaultConfig, null, 2)
  } catch (error) {
    console.error('Failed to load default configuration:', error)
    // Fallback to empty config
    configText.value = JSON.stringify({}, null, 2)
  }
}

const saveConfiguration = async () => {
  if (validationError.value) return

  isSaving.value = true
  try {
    const config = JSON.parse(configText.value)
    
    // Update project with new configuration
    const updatedProject = await projectApi.update(props.project.id, {
      title: props.project.title,
      short_name: props.project.short_name,
      config: Object.keys(config).length > 0 ? config : null
    })

    emit('updated', updatedProject)
    emit('close')
  } catch (error) {
    console.error('Failed to save configuration:', error)
    validationError.value = 'Failed to save configuration. Please try again.'
  } finally {
    isSaving.value = false
  }
}
</script>
