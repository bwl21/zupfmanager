<template>
  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
    <div class="relative top-20 mx-auto p-5 border w-11/12 md:w-3/4 lg:w-2/3 shadow-lg rounded-md bg-white">
      <!-- Header -->
      <div class="flex items-center justify-between pb-4 border-b">
        <h3 class="text-lg font-medium text-gray-900">
          Preview PDFs - {{ project.title }}
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
        <!-- Song Selection -->
        <div class="mb-6">
          <h4 class="text-md font-medium text-gray-900 mb-3">Select Song</h4>
          <div v-if="isLoadingSongs" class="text-center py-4">
            <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600 mx-auto"></div>
            <p class="mt-2 text-sm text-gray-500">Loading songs...</p>
          </div>
          <div v-else-if="projectSongs.length === 0" class="text-center py-4">
            <p class="text-sm text-gray-500">No songs in this project</p>
          </div>
          <div v-else class="space-y-2 max-h-32 overflow-y-auto">
            <div
              v-for="projectSong in projectSongs"
              :key="projectSong.id"
              class="flex items-center p-2 border border-gray-200 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
              :class="{ 'bg-blue-50 border-blue-300': selectedSong?.id === projectSong.song?.id }"
              @click="selectSong(projectSong.song!)"
            >
              <div class="flex-1">
                <div class="font-medium text-gray-900">{{ projectSong.song?.title }}</div>
                <div class="text-sm text-gray-500">{{ projectSong.song?.filename }}</div>
              </div>
              <div class="text-xs text-gray-400">
                Priority: {{ projectSong.priority }}
              </div>
            </div>
          </div>
        </div>

        <!-- ABC Directory -->
        <div class="mb-6">
          <h4 class="text-md font-medium text-gray-900 mb-3">ABC File Directory</h4>
          <div class="space-y-3">
            <div>
              <input
                v-model="abcFileDir"
                type="text"
                class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="/path/to/abc/files"
              />
              <p class="mt-1 text-xs text-gray-500">
                Directory containing ABC files and PDFs
                <span v-if="project.abc_file_dir_preference" class="text-blue-600">
                  (Default: {{ project.abc_file_dir_preference }})
                </span>
              </p>
            </div>
            <button
              @click="findPDFs"
              :disabled="isSearching || !selectedSong || !abcFileDir.trim()"
              class="w-full bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors disabled:opacity-50"
            >
              {{ isSearching ? 'Searching...' : 'Find PDFs for Selected Song' }}
            </button>
          </div>
        </div>

        <!-- PDF List Section -->
        <div v-if="selectedSong">
          <div class="flex items-center justify-between mb-3">
            <h4 class="text-md font-medium text-gray-900">
              PDFs for "{{ selectedSong.title }}"
            </h4>
            <button
              @click="refreshPDFs"
              :disabled="isRefreshing || !abcFileDir.trim()"
              class="text-sm text-blue-600 hover:text-blue-800 transition-colors disabled:opacity-50"
            >
              {{ isRefreshing ? 'Refreshing...' : 'Refresh' }}
            </button>
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
            <p class="text-xs text-gray-400">Select a song and search for PDFs in the ABC directory</p>
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
import { useMutation, useQuery } from '@tanstack/vue-query'
import { songApi, projectSongApi } from '@/services/api'
import type { ProjectResponse, SongResponse } from '@/types/api'

interface Props {
  project: ProjectResponse
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
}>()

// State
const abcFileDir = ref(props.project.abc_file_dir_preference || '')
const selectedSong = ref<SongResponse | null>(null)
const error = ref('')
const pdfs = ref<Array<{ filename: string; size: number; created_at: string }>>([])
const isLoadingPDFs = ref(false)

// Queries
const { data: projectSongsData, isLoading: isLoadingSongs } = useQuery({
  queryKey: ['project-songs', props.project.id],
  queryFn: () => projectSongApi.list(props.project.id),
  refetchOnWindowFocus: false
})

const projectSongs = computed(() => projectSongsData.value?.project_songs || [])

// Mutations
const { mutate: findPDFsMutation, isPending: isSearching } = useMutation({
  mutationFn: (data: { songId: number; abc_file_dir: string }) => 
    songApi.generatePreview(data.songId, { abc_file_dir: data.abc_file_dir }),
  onSuccess: (response) => {
    error.value = ''
    pdfs.value = response.pdf_files.map(filename => ({
      filename,
      size: 0,
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
const selectSong = (song: SongResponse) => {
  selectedSong.value = song
  pdfs.value = [] // Clear previous PDFs
  error.value = ''
}

const findPDFs = () => {
  if (!selectedSong.value || !abcFileDir.value.trim()) return
  
  findPDFsMutation({
    songId: selectedSong.value.id,
    abc_file_dir: abcFileDir.value.trim()
  })
}

const refreshPDFs = async () => {
  if (!selectedSong.value || !abcFileDir.value.trim()) return
  
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
  if (!selectedSong.value || !abcFileDir.value.trim()) {
    error.value = 'Song and ABC file directory are required to open PDF'
    return
  }
  
  const url = songApi.getPreviewPDFUrl(selectedSong.value.id, filename, abcFileDir.value.trim())
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
</script>
