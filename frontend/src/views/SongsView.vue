<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div>
      <h1 class="text-3xl font-bold text-gray-900">Songs</h1>
      <p class="mt-2 text-gray-600">Browse and search your ABC notation files</p>
    </div>

    <!-- Search Bar -->
    <div class="bg-white rounded-lg shadow p-6">
      <div class="flex flex-col md:flex-row gap-4">
        <div class="flex-1">
          <label for="search" class="sr-only">Search songs</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input
              id="search"
              v-model="searchQuery"
              type="text"
              placeholder="Search songs by title, filename, or genre..."
              class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
              @input="debouncedSearch"
            />
          </div>
        </div>
        <div class="flex gap-2">
          <button
            @click="searchOptions.title = !searchOptions.title"
            :class="[
              'px-3 py-2 rounded-md text-sm font-medium transition-colors',
              searchOptions.title
                ? 'bg-blue-100 text-blue-800 border border-blue-200'
                : 'bg-gray-100 text-gray-700 border border-gray-200 hover:bg-gray-200'
            ]"
          >
            Title
          </button>
          <button
            @click="searchOptions.filename = !searchOptions.filename"
            :class="[
              'px-3 py-2 rounded-md text-sm font-medium transition-colors',
              searchOptions.filename
                ? 'bg-blue-100 text-blue-800 border border-blue-200'
                : 'bg-gray-100 text-gray-700 border border-gray-200 hover:bg-gray-200'
            ]"
          >
            Filename
          </button>
          <button
            @click="searchOptions.genre = !searchOptions.genre"
            :class="[
              'px-3 py-2 rounded-md text-sm font-medium transition-colors',
              searchOptions.genre
                ? 'bg-blue-100 text-blue-800 border border-blue-200'
                : 'bg-gray-100 text-gray-700 border border-gray-200 hover:bg-gray-200'
            ]"
          >
            Genre
          </button>
        </div>
      </div>
      <div v-if="searchQuery" class="mt-4 flex items-center justify-between">
        <p class="text-sm text-gray-600">
          {{ isSearching ? 'Searching...' : `Found ${searchResults?.count || 0} songs` }}
        </p>
        <button
          @click="clearSearch"
          class="text-sm text-blue-600 hover:text-blue-500 font-medium"
        >
          Clear search
        </button>
      </div>
    </div>

    <!-- Songs List -->
    <div class="bg-white rounded-lg shadow">
      <div class="px-6 py-4 border-b border-gray-200">
        <h2 class="text-lg font-medium text-gray-900">
          {{ searchQuery ? 'Search Results' : 'All Songs' }}
          <span class="text-sm font-normal text-gray-500 ml-2">
            ({{ displayedSongs?.count || 0 }} songs)
          </span>
        </h2>
      </div>

      <!-- Songs Grid -->
      <div v-if="displayedSongs?.songs.length" class="p-6">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div
            v-for="song in displayedSongs.songs"
            :key="song.id"
            class="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow cursor-pointer"
            @click="$router.push(`/songs/${song.id}`)"
          >
            <div class="flex items-start justify-between mb-2">
              <h3 class="font-semibold text-gray-900 truncate">{{ song.title }}</h3>
              <svg class="w-5 h-5 text-gray-400 flex-shrink-0 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
              </svg>
            </div>
            <p class="text-sm text-gray-600 mb-1">{{ song.filename }}</p>
            
            <!-- Project Badges -->
            <div v-if="song.projects && song.projects.length > 0" class="flex flex-wrap gap-1 mb-2">
              <button
                v-for="project in song.projects"
                :key="project.id"
                @click.stop="$router.push(`/projects/${project.id}`)"
                class="inline-flex items-center px-2 py-1 text-xs font-medium rounded-full bg-blue-100 text-blue-800 hover:bg-blue-200 transition-colors cursor-pointer"
                :title="`Go to project: ${project.title}`"
              >
                {{ project.short_name }}
              </button>
            </div>
            
            <div class="flex items-center justify-between text-xs text-gray-500">
              <span v-if="song.genre" class="bg-gray-100 px-2 py-1 rounded">{{ song.genre }}</span>
              <span v-else class="bg-gray-100 px-2 py-1 rounded">No genre</span>
              <span>ID: {{ song.id }}</span>
            </div>
            <div v-if="song.copyright || song.tocinfo" class="mt-2 text-xs text-gray-500">
              <p v-if="song.copyright" class="truncate">© {{ song.copyright }}</p>
              <p v-if="song.tocinfo" class="truncate">{{ song.tocinfo }}</p>
            </div>
            
            <!-- Action Buttons -->
            <div class="mt-3 flex justify-end gap-2">
              <button
                @click.stop="deleteSong(song)"
                :disabled="(song.projects && song.projects.length > 0) || isDeleting"
                class="inline-flex items-center px-2 py-1 text-xs font-medium rounded text-red-600 hover:text-red-800 hover:bg-red-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                :title="song.projects && song.projects.length > 0 ? 'Cannot delete: song is used in projects' : 'Delete song'"
              >
                <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
                Delete
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="!isLoading && !isSearching" class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900">
          {{ searchQuery ? 'No songs found' : 'No songs available' }}
        </h3>
        <p class="mt-1 text-sm text-gray-500">
          {{ searchQuery ? 'Try adjusting your search terms or filters.' : 'Import some ABC files to get started.' }}
        </p>
        <div v-if="!searchQuery" class="mt-6">
          <RouterLink
            to="/import"
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors"
          >
            Import Songs
          </RouterLink>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading || isSearching" class="text-center py-12">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <p class="mt-2 text-gray-600">{{ isSearching ? 'Searching songs...' : 'Loading songs...' }}</p>
      </div>

      <!-- Error State -->
      <div v-if="error || searchError" class="p-6">
        <div class="bg-red-50 border border-red-200 rounded-lg p-4">
          <div class="flex">
            <svg class="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">Error loading songs</h3>
              <p class="mt-2 text-sm text-red-700">{{ (error || searchError)?.message }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { RouterLink } from 'vue-router'
import { songApi, projectApi } from '@/services/api'
import { useDebounceFn } from '@vueuse/core'
import type { SongResponse } from '@/types/api'

const queryClient = useQueryClient()

// Search state
const searchQuery = ref('')
const searchOptions = reactive({
  title: true,
  filename: false,
  genre: false
})

// Fetch all songs
const { data: allSongs, isLoading, error } = useQuery({
  queryKey: ['songs'],
  queryFn: songApi.list
})

// Search songs
const { data: searchResults, isLoading: isSearching, error: searchError } = useQuery({
  queryKey: ['songs', 'search', searchQuery, searchOptions],
  queryFn: () => songApi.search(searchQuery.value, searchOptions),
  enabled: computed(() => searchQuery.value.length > 0)
})


// Display songs directly from API responses (now include project info)
const displayedSongs = computed(() => {
  return searchQuery.value ? searchResults.value : allSongs.value
})

// Debounced search function
const debouncedSearch = useDebounceFn(() => {
  // The query will automatically refetch due to reactive dependencies
  // Project info will be loaded by the searchResults watcher
}, 300)

function clearSearch() {
  searchQuery.value = ''
  // Project info will be loaded by the allSongs watcher
}

// Delete song mutation
const { mutate: deleteSongMutation, isPending: isDeleting } = useMutation({
  mutationFn: songApi.delete,
  onSuccess: (_, songId) => {
    // Invalidate all songs-related queries to refresh the list
    queryClient.invalidateQueries({ queryKey: ['songs'] })
    queryClient.invalidateQueries({ predicate: (query) => query.queryKey[0] === 'songs' })
    
    // Show success message
    console.log(`Song ${songId} deleted successfully`)
  },
  onError: (error: any) => {
    let errorMessage = 'Failed to delete song'
    
    if (error.response?.status === 409) {
      errorMessage = error.response.data.message || 'Song is used in projects and cannot be deleted'
    } else if (error.response?.status === 404) {
      errorMessage = 'Song not found'
    } else if (error.response?.data?.message) {
      errorMessage = error.response.data.message
    } else if (error.message) {
      errorMessage = error.message
    }
    
    alert(`Error: ${errorMessage}`)
  }
})

// Delete song function with confirmation
function deleteSong(song: SongResponse) {
  // Double-check project associations
  if (song.projects && song.projects.length > 0) {
    const projectNames = song.projects.map(p => p.title).join(', ')
    alert(`Cannot delete song "${song.title}": it is used in the following project(s):\n\n${projectNames}\n\nRemove the song from all projects first.`)
    return
  }

  // Confirmation dialog
  const message = `Are you sure you want to delete the song "${song.title}"?\n\nFilename: ${song.filename}\nID: ${song.id}\n\nThis action cannot be undone.`
  const confirmed = confirm(message)
  
  if (confirmed) {
    deleteSongMutation(song.id)
  }
}
</script>
