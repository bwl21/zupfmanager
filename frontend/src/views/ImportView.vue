<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div>
      <h1 class="text-3xl font-bold text-gray-900">Import ABC Files</h1>
      <p class="mt-2 text-gray-600">Import ABC notation files into your song database</p>
    </div>

    <!-- Import Options -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Single File Import -->
      <div class="bg-white rounded-lg shadow">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-lg font-medium text-gray-900">Import Single File</h2>
          <p class="mt-1 text-sm text-gray-600">Import a single ABC notation file</p>
        </div>
        <div class="p-6">
          <form @submit.prevent="importSingleFile" class="space-y-4">
            <div>
              <label for="file_path" class="block text-sm font-medium text-gray-700">File Path</label>
              <input
                id="file_path"
                v-model="singleFileForm.file_path"
                type="text"
                required
                class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="/path/to/song.abc"
              />
              <p class="mt-1 text-xs text-gray-500">Enter the full path to the ABC file</p>
            </div>
            <button
              type="submit"
              :disabled="isImportingSingle"
              class="w-full bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors disabled:opacity-50"
            >
              {{ isImportingSingle ? 'Importing...' : 'Import File' }}
            </button>
          </form>
        </div>
      </div>

      <!-- Directory Import -->
      <div class="bg-white rounded-lg shadow">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-lg font-medium text-gray-900">Import Directory</h2>
          <p class="mt-1 text-sm text-gray-600">Import all ABC files from a directory</p>
        </div>
        <div class="p-6">
          <form @submit.prevent="importDirectory" class="space-y-4">
            <div>
              <label for="directory_path" class="block text-sm font-medium text-gray-700">Directory Path</label>
              <input
                id="directory_path"
                v-model="directoryForm.directory_path"
                type="text"
                required
                class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="/path/to/songs/"
              />
              <p class="mt-1 text-xs text-gray-500">Enter the full path to the directory containing ABC files</p>
            </div>
            <button
              type="submit"
              :disabled="isImportingDirectory"
              class="w-full bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-lg font-medium transition-colors disabled:opacity-50"
            >
              {{ isImportingDirectory ? 'Importing...' : 'Import Directory' }}
            </button>
          </form>
        </div>
      </div>
    </div>

    <!-- Quick Import Suggestions -->
    <div class="bg-blue-50 border border-blue-200 rounded-lg p-6">
      <div class="flex">
        <svg class="w-5 h-5 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-blue-800">Quick Import</h3>
          <p class="mt-1 text-sm text-blue-700">
            Try importing the test songs from: <code class="bg-blue-100 px-1 rounded">/workspaces/zupfmanager/test_songs/</code>
          </p>
          <div class="mt-3">
            <button
              @click="importTestSongs"
              :disabled="isImportingDirectory"
              class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1 rounded text-sm font-medium transition-colors disabled:opacity-50"
            >
              Import Test Songs
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Import Results -->
    <div v-if="lastImportResult" class="bg-white rounded-lg shadow">
      <div class="px-6 py-4 border-b border-gray-200">
        <h2 class="text-lg font-medium text-gray-900">Import Results</h2>
      </div>
      <div class="p-6">
        <!-- Summary -->
        <div class="mb-6">
          <div class="grid grid-cols-2 md:grid-cols-5 gap-4">
            <div class="text-center">
              <div class="text-2xl font-bold text-gray-900">{{ lastImportResult.summary.total }}</div>
              <div class="text-sm text-gray-500">Total</div>
            </div>
            <div class="text-center">
              <div class="text-2xl font-bold text-green-600">{{ lastImportResult.summary.created }}</div>
              <div class="text-sm text-gray-500">Created</div>
            </div>
            <div class="text-center">
              <div class="text-2xl font-bold text-blue-600">{{ lastImportResult.summary.updated }}</div>
              <div class="text-sm text-gray-500">Updated</div>
            </div>
            <div class="text-center">
              <div class="text-2xl font-bold text-gray-600">{{ lastImportResult.summary.unchanged }}</div>
              <div class="text-sm text-gray-500">Unchanged</div>
            </div>
            <div class="text-center">
              <div class="text-2xl font-bold text-red-600">{{ lastImportResult.summary.errors }}</div>
              <div class="text-sm text-gray-500">Errors</div>
            </div>
          </div>
        </div>

        <!-- Detailed Results -->
        <div v-if="lastImportResult.results.length > 0">
          <h3 class="text-sm font-medium text-gray-900 mb-3">Detailed Results</h3>
          <div class="space-y-2 max-h-64 overflow-y-auto">
            <div
              v-for="result in lastImportResult.results"
              :key="result.filename"
              class="flex items-center justify-between p-3 border border-gray-200 rounded-lg"
            >
              <div class="flex-1">
                <div class="flex items-center">
                  <span class="font-medium text-gray-900">{{ result.title || result.filename }}</span>
                  <span
                    :class="[
                      'ml-2 px-2 py-1 rounded-full text-xs font-medium',
                      result.action === 'created' ? 'bg-green-100 text-green-800' :
                      result.action === 'updated' ? 'bg-blue-100 text-blue-800' :
                      result.action === 'unchanged' ? 'bg-gray-100 text-gray-800' :
                      'bg-red-100 text-red-800'
                    ]"
                  >
                    {{ result.action }}
                  </span>
                </div>
                <div class="text-sm text-gray-500">{{ result.filename }}</div>
                <div v-if="result.changes && result.changes.length > 0" class="text-xs text-gray-400">
                  Changes: {{ result.changes.join(', ') }}
                </div>
                <div v-if="result.error" class="text-xs text-red-600">
                  Error: {{ result.error }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Success Actions -->
        <div v-if="lastImportResult.success && lastImportResult.summary.created > 0" class="mt-6 flex gap-3">
          <RouterLink
            to="/songs"
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors"
          >
            View Imported Songs
          </RouterLink>
          <button
            @click="clearResults"
            class="bg-gray-100 hover:bg-gray-200 text-gray-700 px-4 py-2 rounded-lg font-medium transition-colors"
          >
            Clear Results
          </button>
        </div>
      </div>
    </div>

    <!-- Error Display -->
    <div v-if="importError" class="bg-red-50 border border-red-200 rounded-lg p-4">
      <div class="flex">
        <svg class="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Import Error</h3>
          <p class="mt-2 text-sm text-red-700">{{ importError.message }}</p>
        </div>
      </div>
    </div>

    <!-- Help Section -->
    <div class="bg-gray-50 rounded-lg p-6">
      <h2 class="text-lg font-medium text-gray-900 mb-4">Import Help</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <h3 class="text-sm font-medium text-gray-900 mb-2">Supported Formats</h3>
          <ul class="text-sm text-gray-600 space-y-1">
            <li>• ABC notation files (.abc)</li>
            <li>• Files must contain valid ABC headers (T: for title)</li>
            <li>• UTF-8 encoding recommended</li>
          </ul>
        </div>
        <div>
          <h3 class="text-sm font-medium text-gray-900 mb-2">Import Behavior</h3>
          <ul class="text-sm text-gray-600 space-y-1">
            <li>• Existing songs are updated if filename matches</li>
            <li>• New songs are created with unique filenames</li>
            <li>• Invalid files are skipped with error messages</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { RouterLink } from 'vue-router'
import { importApi } from '@/services/api'
import type { ImportResponse } from '@/types/api'

const queryClient = useQueryClient()

// Form state
const singleFileForm = reactive({
  file_path: ''
})

const directoryForm = reactive({
  directory_path: ''
})

// Results state
const lastImportResult = ref<ImportResponse | null>(null)
const importError = ref<any>(null)

// Single file import mutation
const { mutate: importFileMutation, isPending: isImportingSingle } = useMutation({
  mutationFn: importApi.file,
  onSuccess: (data) => {
    lastImportResult.value = data
    importError.value = null
    singleFileForm.file_path = ''
    // Invalidate songs query to refresh the list
    queryClient.invalidateQueries({ queryKey: ['songs'] })
  },
  onError: (error) => {
    importError.value = error
    lastImportResult.value = null
  }
})

// Directory import mutation
const { mutate: importDirectoryMutation, isPending: isImportingDirectory } = useMutation({
  mutationFn: importApi.directory,
  onSuccess: (data) => {
    lastImportResult.value = data
    importError.value = null
    directoryForm.directory_path = ''
    // Invalidate songs query to refresh the list
    queryClient.invalidateQueries({ queryKey: ['songs'] })
  },
  onError: (error) => {
    importError.value = error
    lastImportResult.value = null
  }
})

function importSingleFile() {
  importFileMutation({ file_path: singleFileForm.file_path })
}

function importDirectory() {
  importDirectoryMutation({ directory_path: directoryForm.directory_path })
}

function importTestSongs() {
  directoryForm.directory_path = '/workspaces/zupfmanager/test_songs/'
  importDirectory()
}

function clearResults() {
  lastImportResult.value = null
  importError.value = null
}
</script>
