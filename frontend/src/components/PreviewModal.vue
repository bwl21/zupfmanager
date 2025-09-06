<template>
  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
    <div class="relative top-20 mx-auto p-5 border w-11/12 md:w-3/4 lg:w-1/2 shadow-lg rounded-md bg-white">
      <!-- Header -->
      <div class="flex items-center justify-between pb-4 border-b">
        <h3 class="text-lg font-medium text-gray-900">
          Preview PDFs - {{ song.title }}
        </h3>
        <button
          @click="$emit('close')"
          class="text-gray-400 hover:text-gray-600 transition-colors"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Content -->
      <div class="mt-4">
        <!-- Find PDFs Section -->
        <div class="mb-6">
          <h4 class="text-md font-medium text-gray-900 mb-3">Find Existing PDFs</h4>
          <div class="space-y-3">
            <div>
              <label for="abc_file_dir" class="block text-sm font-medium text-gray-700">ABC File Directory</label>
              <input
                id="abc_file_dir"
                v-model="abcFileDir"
                type="text"
                required
                class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="/path/to/abc/files"
              />
              <p class="mt-1 text-xs text-gray-500">Enter the directory containing the ABC file and PDFs</p>
            </div>
            <button
              @click="findPDFs"
              :disabled="isSearching || !abcFileDir.trim()"
              class="w-full bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors disabled:opacity-50"
            >
              {{ isSearching ? 'Searching...' : 'Find PDFs' }}
            </button>
          </div>
        </div>

        <!-- PDF List Section -->
        <div>
          <div class="flex items-center justify-between mb-3">
            <h4 class="text-md font-medium text-gray-900">Available PDFs</h4>
            <div class="flex gap-4">
              <button
                @click="refreshPDFs"
                :disabled="isRefreshing || !abcFileDir.trim()"
                class="text-sm text-blue-600 hover:text-blue-800 transition-colors disabled:opacity-50"
              >
                {{ isRefreshing ? 'Refreshing...' : 'Refresh' }}
              </button>
            </div>
          </div>

          <!-- Loading State -->
          <div v-if="isLoadingPDFs" class="text-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
            <p class="mt-2 text-sm text-gray-500">Loading PDFs...</p>
          </div>

          <!-- Empty State -->
          <div v-else-if="!pdfs.length" class="text-center py-8">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <p class="mt-2 text-sm text-gray-500">No PDFs found</p>
            <p class="text-xs text-gray-400">Enter the ABC directory to search for existing PDFs</p>
          </div>

          <!-- PDF List -->
          <div v-else class="space-y-2 max-h-64 overflow-y-auto">
            <div
              v-for="pdf in pdfs"
              :key="pdf.filename"
              class="flex items-center justify-between p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
            >
              <div class="flex-1">
                <div class="flex items-center">
                  <svg class="w-5 h-5 text-red-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  <span class="font-medium text-gray-900">{{ pdf.filename }}</span>
                </div>
                <div class="text-sm text-gray-500 mt-1">
                  {{ formatFileSize(pdf.size) }} • {{ formatDate(pdf.created_at) }}
                </div>
              </div>
              <button
                @click="openPDF(pdf.filename)"
                class="ml-4 bg-blue-600 hover:bg-blue-700 text-white px-3 py-1 rounded text-sm font-medium transition-colors"
              >
                Open
              </button>
            </div>
          </div>
        </div>

        <!-- Error Display -->
        <div v-if="error" class="mt-4 bg-red-50 border border-red-200 rounded-lg p-4">
          <div class="flex">
            <svg class="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">Error</h3>
              <p class="mt-2 text-sm text-red-700">{{ error }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { songApi } from '@/services/api'
import type { SongResponse } from '@/types/api'

interface Props {
  song: SongResponse
  project?: { abc_file_dir_preference?: string } | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
}>()

const queryClient = useQueryClient()

// State
const abcFileDir = ref('')
const error = ref('')
const pdfs = ref<Array<{ filename: string; size: number; created_at: string }>>([])
const isLoadingPDFs = ref(false)

// Initialize abc_file_dir with project preference
onMounted(() => {
  if (props.project?.abc_file_dir_preference) {
    abcFileDir.value = props.project.abc_file_dir_preference
    // Auto-search for PDFs if we have a default directory
    findPDFs()
  }
})

// Mutations
const { mutate: findPDFsMutation, isPending: isSearching } = useMutation({
  mutationFn: (data: { abc_file_dir: string }) => 
    songApi.generatePreview(props.song.id, data),
  onSuccess: (response) => {
    error.value = ''
    // Convert filenames to PDF objects with mock data since we don't have file stats
    pdfs.value = response.pdf_files.map(filename => ({
      filename,
      size: 0, // We don't have size info when just finding files
      created_at: new Date().toISOString()
    }))
  },
  onError: (err: any) => {
    error.value = err.message || 'Failed to find PDFs'
    pdfs.value = []
  }
})

// State for refresh button
const isRefreshing = ref(false)

// Methods
const findPDFs = () => {
  if (!abcFileDir.value.trim()) return
  
  findPDFsMutation({
    abc_file_dir: abcFileDir.value.trim()
  })
}

const refreshPDFs = async () => {
  if (!abcFileDir.value.trim()) return
  
  isRefreshing.value = true
  try {
    await findPDFs()
    error.value = ''
  } catch (err: any) {
    error.value = err.message || 'Failed to refresh PDFs'
  } finally {
    isRefreshing.value = false
  }
}

const openPDF = (filename: string) => {
  if (!abcFileDir.value.trim()) {
    error.value = 'ABC file directory is required to open PDF'
    return
  }
  
  const url = songApi.getPreviewPDFUrl(props.song.id, filename, abcFileDir.value.trim())
  window.open(url, '_blank')
}

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return 'Unknown size'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

// No automatic loading on mount since we need the ABC directory first
</script>
