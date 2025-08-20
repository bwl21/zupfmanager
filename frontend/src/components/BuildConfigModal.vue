<template>
  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
    <div class="relative top-20 mx-auto p-5 border w-11/12 md:w-2/3 lg:w-1/2 shadow-lg rounded-md bg-white">
      <div class="mt-3">
        <!-- Header -->
        <div class="flex items-center justify-between pb-4 border-b">
          <h3 class="text-lg font-medium text-gray-900">Configure Project Build</h3>
          <button
            @click="$emit('close')"
            class="text-gray-400 hover:text-gray-600"
          >
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Build Configuration Form -->
        <div class="mt-6 space-y-4">
          <!-- Output Directory -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Output Directory
            </label>
            <input
              v-model="buildConfig.output_dir"
              type="text"
              placeholder="Leave empty to use project short name"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
            <p class="mt-1 text-xs text-gray-500">
              Directory where build outputs will be generated
            </p>
          </div>
          
          <!-- ABC File Directory -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              ABC File Directory
            </label>
            <div class="space-y-2">
              <!-- Directory Input -->
              <div class="flex space-x-2">
                <div class="flex-1 relative">
                  <input
                    v-model="buildConfig.abc_file_dir"
                    type="text"
                    :placeholder="isLoadingDefaults ? 'Loading defaults...' : 'Full path to ABC files directory (e.g., /home/user/music/abc)'"
                    :disabled="isLoadingDefaults"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:bg-gray-50 disabled:text-gray-500"
                  />
                  <div v-if="isLoadingDefaults" class="absolute right-3 top-1/2 transform -translate-y-1/2">
                    <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-indigo-600"></div>
                  </div>
                </div>
              </div>

            </div>
            <p class="mt-1 text-xs text-gray-500">
              Full path to directory containing ABC notation files.
            </p>
          </div>

          <!-- Priority Threshold -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Priority Threshold
            </label>
            <select
              v-model="buildConfig.priority_threshold"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option :value="1">1 - Only highest priority songs</option>
              <option :value="2">2 - High priority and above</option>
              <option :value="3">3 - Medium priority and above</option>
              <option :value="4">4 - All songs (default)</option>
            </select>
            <p class="mt-1 text-xs text-gray-500">
              Only include songs with priority equal to or higher than this threshold
            </p>
          </div>
          
          <!-- Sample ID -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Sample ID (optional)
            </label>
            <input
              v-model="buildConfig.sample_id"
              type="text"
              placeholder="e.g., demo, preview, final"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
            <p class="mt-1 text-xs text-gray-500">
              Identifier for this build variant
            </p>
          </div>
        </div>

        <!-- Build Preview -->
        <div class="mt-6 p-4 bg-gray-50 rounded-md">
          <h4 class="text-sm font-medium text-gray-900 mb-2">Build Summary</h4>
          <div class="text-sm text-gray-600 space-y-1">
            <p><strong>Output Directory:</strong> {{ buildConfig.output_dir || 'Project default' }}</p>
            <p><strong>Priority Filter:</strong> {{ getPriorityDescription(buildConfig.priority_threshold) }}</p>
            <p v-if="buildConfig.abc_file_dir"><strong>ABC Files:</strong> {{ buildConfig.abc_file_dir }}</p>
            <p v-if="buildConfig.sample_id"><strong>Sample ID:</strong> {{ buildConfig.sample_id }}</p>
          </div>
        </div>

        <!-- Path Warning -->
        <div v-if="buildConfig.abc_file_dir && !isValidPath(buildConfig.abc_file_dir)" class="mt-4 p-3 bg-orange-50 border border-orange-200 rounded-md">
          <div class="flex">
            <svg class="h-5 w-5 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
            <div class="ml-3">
              <p class="text-sm text-orange-700">
                <strong>Path may be incomplete:</strong> Please ensure you enter the complete absolute path to the ABC files directory 
                (e.g., <code class="bg-orange-100 px-1 rounded">/home/user/music/abc</code> or <code class="bg-orange-100 px-1 rounded">C:\Users\User\Music\ABC</code>).
              </p>
            </div>
          </div>
        </div>

        <!-- General Warning -->
        <div class="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-md">
          <div class="flex">
            <svg class="h-5 w-5 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
            <div class="ml-3">
              <p class="text-sm text-yellow-700">
                Building a project may take several minutes depending on the number of songs and complexity.
                You can monitor the progress after starting the build.
              </p>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="mt-6 flex justify-end space-x-3">
          <button
            @click="$emit('close')"
            class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            Cancel
          </button>
          <button
            @click="startBuild"
            :disabled="isStarting"
            class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="isStarting">Starting Build...</span>
            <span v-else>Start Build</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { projectBuildApi, projectApi } from '@/services/api'
import type { BuildProjectRequest, BuildResultResponse } from '@/types/api'

interface Props {
  projectId: number
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
  buildStarted: [build: BuildResultResponse]
}>()

// State
const isStarting = ref(false)
const isLoadingDefaults = ref(false)
const buildConfig = ref<BuildProjectRequest>({
  output_dir: '',
  abc_file_dir: '',
  priority_threshold: 4,
  sample_id: ''
})

// Methods
const startBuild = async () => {
  isStarting.value = true
  try {
    // Clean up empty values
    const config: BuildProjectRequest = {}
    if (buildConfig.value.output_dir?.trim()) {
      config.output_dir = buildConfig.value.output_dir.trim()
    }
    if (buildConfig.value.abc_file_dir?.trim()) {
      config.abc_file_dir = buildConfig.value.abc_file_dir.trim()
      
      // Save the abc_file_dir as preference if it's a valid path (not a placeholder)
      if (!config.abc_file_dir.startsWith('[Enter full path to:')) {
        await saveAbcFileDirPreference(config.abc_file_dir)
      }
    }
    if (buildConfig.value.priority_threshold && buildConfig.value.priority_threshold !== 4) {
      config.priority_threshold = buildConfig.value.priority_threshold
    }
    if (buildConfig.value.sample_id?.trim()) {
      config.sample_id = buildConfig.value.sample_id.trim()
    }

    const result = await projectBuildApi.build(props.projectId, config)
    emit('buildStarted', result)
  } catch (err) {
    alert('Failed to start build')
    console.error('Failed to start build:', err)
  } finally {
    isStarting.value = false
  }
}

const loadDefaults = async () => {
  isLoadingDefaults.value = true
  try {
    const defaults = await projectBuildApi.getDefaults(props.projectId)
    buildConfig.value = {
      output_dir: defaults.output_dir || '',
      abc_file_dir: defaults.abc_file_dir || '',
      priority_threshold: defaults.priority_threshold || 4,
      sample_id: defaults.sample_id || ''
    }
  } catch (err) {
    console.error('Failed to load build defaults:', err)
    // Keep the default values if loading fails
  } finally {
    isLoadingDefaults.value = false
  }
}


const saveAbcFileDirPreference = async (directoryPath: string) => {
  try {
    await projectApi.updateAbcFileDir(props.projectId, directoryPath)
    console.log(`Saved ABC file directory preference: ${directoryPath}`)
  } catch (err) {
    console.error('Failed to save ABC file directory preference:', err)
    // Don't show error to user as this is not critical for the build process
  }
}

const isValidPath = (path: string) => {
  if (!path || path.trim() === '') return true // Empty is valid (uses defaults)
  if (path.startsWith('[Enter full path to:')) return false // Placeholder text
  
  // Check for common path patterns
  const hasAbsolutePath = path.startsWith('/') || // Unix/Linux/Mac
                         /^[A-Za-z]:\\/.test(path) || // Windows (C:\)
                         path.startsWith('\\\\') // UNC path (\\server\share)
  
  return hasAbsolutePath && path.length > 3 // Must be more than just root
}

const getPriorityDescription = (threshold: number | undefined) => {
  const descriptions = {
    1: 'Only highest priority songs (Priority 1)',
    2: 'High priority and above (Priority 1-2)',
    3: 'Medium priority and above (Priority 1-3)',
    4: 'All songs (Priority 1-4)'
  }
  return descriptions[threshold as keyof typeof descriptions] || 'All songs'
}

// Lifecycle
onMounted(() => {
  loadDefaults()
})
</script>
