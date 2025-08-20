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
              <div class="mt-1 flex space-x-2">
                <input
                  id="file_path"
                  v-model="singleFileForm.file_path"
                  type="text"
                  required
                  class="flex-1 border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder="/path/to/song.abc"
                  @input="clearFileInfo"
                />
                <button
                  type="button"
                  @click="selectFile"
                  class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                  title="Browse for ABC file"
                >
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                </button>
              </div>
              <p class="mt-1 text-xs text-gray-500">Enter the full path to the ABC file or click Browse to select</p>
              <!-- File selection info -->
              <div v-if="selectedFileInfo" class="mt-2 p-2 bg-green-50 border border-green-200 rounded text-xs">
                <p class="text-green-800 font-medium">✓ File selected: {{ selectedFileInfo.name }}</p>
                <p class="text-green-600">Size: {{ formatFileSize(selectedFileInfo.size) }}</p>
                <p class="text-green-600 mt-1"><strong>Now enter the complete path to this file above</strong> (e.g., /home/user/Documents/{{ selectedFileInfo.name }})</p>
              </div>
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
              <div class="mt-1 flex space-x-2">
                <input
                  id="directory_path"
                  v-model="directoryForm.directory_path"
                  type="text"
                  required
                  class="flex-1 border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder="/path/to/songs/"
                  @input="clearDirectoryInfo"
                />
                <button
                  type="button"
                  @click="selectDirectory"
                  class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                  title="Browse for directory"
                >
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-5l-2-2H6a2 2 0 00-2 2z" />
                  </svg>
                </button>
              </div>
              <p class="mt-1 text-xs text-gray-500">Enter the full path to the directory containing ABC files or click Browse to select</p>
              <!-- Directory selection info -->
              <div v-if="selectedDirectoryInfo" class="mt-2 p-2 bg-blue-50 border border-blue-200 rounded text-xs">
                <p class="text-blue-800 font-medium">✓ Directory selected: {{ selectedDirectoryInfo.name }}</p>
                <p class="text-blue-600">Found {{ selectedDirectoryInfo.abcCount }} ABC files</p>
                <p class="text-blue-600 mt-1"><strong>Now enter the complete path to this directory above</strong> (e.g., /home/user/Documents/{{ selectedDirectoryInfo.name }})</p>
              </div>
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
    
    <!-- Hidden file inputs for file/directory selection -->
    <input
      ref="fileInput"
      type="file"
      accept=".abc"
      @change="handleFileSelection"
      class="hidden"
    />
    <input
      ref="directoryInput"
      type="file"
      webkitdirectory
      directory
      multiple
      @change="handleDirectorySelection"
      class="hidden"
    />
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

// File/Directory picker state
const fileInput = ref<HTMLInputElement | null>(null)
const directoryInput = ref<HTMLInputElement | null>(null)
const selectedFileInfo = ref<{name: string, size: number} | null>(null)
const selectedDirectoryInfo = ref<{name: string, abcCount: number} | null>(null)

// Single file import mutation
const { mutate: importFileMutation, isPending: isImportingSingle } = useMutation({
  mutationFn: importApi.file,
  onSuccess: (data) => {
    lastImportResult.value = data
    importError.value = null
    singleFileForm.file_path = ''
    selectedFileInfo.value = null
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
    selectedDirectoryInfo.value = null
    // Invalidate songs query to refresh the list
    queryClient.invalidateQueries({ queryKey: ['songs'] })
  },
  onError: (error) => {
    importError.value = error
    lastImportResult.value = null
  }
})

function importSingleFile() {
  let filePath = singleFileForm.file_path.trim()
  if (!filePath) {
    alert('Please enter the complete file path.')
    return
  }
  importFileMutation({ file_path: filePath })
}

function importDirectory() {
  let directoryPath = directoryForm.directory_path.trim()
  if (!directoryPath) {
    alert('Please enter the complete directory path.')
    return
  }
  importDirectoryMutation({ directory_path: directoryPath })
}

function selectFile() {
  fileInput.value?.click()
}

function selectDirectory() {
  directoryInput.value?.click()
}

function handleFileSelection(event: Event) {
  const target = event.target as HTMLInputElement
  const files = target.files
  if (files && files.length > 0) {
    const file = files[0]
    
    // Check if it's an ABC file
    if (!file.name.toLowerCase().endsWith('.abc')) {
      alert('Please select an ABC file (.abc extension)')
      target.value = ''
      return
    }

    selectedFileInfo.value = {
      name: file.name,
      size: file.size
    }
    
    // Clear the field and let user enter the full path
    singleFileForm.file_path = ''
    
    // Reset the input
    target.value = ''
  }
}

function handleDirectorySelection(event: Event) {
  const target = event.target as HTMLInputElement
  const files = target.files
  if (files && files.length > 0) {
    const firstFile = files[0]
    const relativePath = firstFile.webkitRelativePath
    
    let directoryName = ''
    let abcCount = 0
    
    if (relativePath) {
      // Extract the root directory name from the relative path
      const pathParts = relativePath.split('/')
      directoryName = pathParts[0]
      
      // Count ABC files found
      const abcFiles = Array.from(files)
        .filter((file: File) => file.name.toLowerCase().endsWith('.abc'))
      abcCount = abcFiles.length
      
      console.log(`Found ${abcCount} ABC files in directory: ${directoryName}`)
    } else {
      // Fallback if webkitRelativePath is not available
      directoryName = 'Selected Directory'
    }
    
    selectedDirectoryInfo.value = {
      name: directoryName,
      abcCount: abcCount
    }
    
    // Clear the field and let user enter the full path
    directoryForm.directory_path = ''
    
    // Reset the input
    target.value = ''
  }
}

function importTestSongs() {
  directoryForm.directory_path = '/workspaces/zupfmanager/test_songs/'
  importDirectory()
}

function clearResults() {
  lastImportResult.value = null
  importError.value = null
}

const clearFileInfo = () => {
  selectedFileInfo.value = null
}


const clearDirectoryInfo = () => {
  selectedDirectoryInfo.value = null
}

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}
</script>
