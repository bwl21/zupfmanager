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
import { useQuery } from '@tanstack/vue-query'
import { RouterLink } from 'vue-router'
import { songApi, projectApi } from '@/services/api'
import { useDebounceFn } from '@vueuse/core'
import type { SongResponse } from '@/types/api'

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


// Enhanced songs with project information
const songsWithProjects = ref<SongResponse[]>([])

// Load project information for songs efficiently
const loadProjectsForSongs = async (songs: SongResponse[]) => {
  try {
    // Get all projects first
    const projectsResponse = await projectApi.list()
    const projects = projectsResponse.projects
    
    // Create a map to store song-to-projects relationships
    const songProjectMap = new Map<number, Array<{id: number, title: string, short_name: string}>>()
    
    // Initialize map with empty arrays for all songs
    songs.forEach(song => {
      songProjectMap.set(song.id, [])
    })
    
    // Load all project-song relationships in parallel
    const projectSongPromises = projects.map(async (project) => {
      try {
        const projectSongs = await projectApi.getSongs(project.id)
        return {
          project: {
            id: project.id,
            title: project.title,
            short_name: project.short_name
          },
          songIds: projectSongs.project_songs.map(ps => ps.song_id)
        }
      } catch (err) {
        console.warn(`Failed to load songs for project ${project.id}:`, err)
        return null
      }
    })
    
    // Wait for all project-song relationships to load
    const projectSongResults = await Promise.all(projectSongPromises)
    
    // Build the song-to-projects map
    projectSongResults.forEach(result => {
      if (result) {
        result.songIds.forEach(songId => {
          const songProjects = songProjectMap.get(songId)
          if (songProjects) {
            songProjects.push(result.project)
          }
        })
      }
    })
    
    // Enhance songs with project information
    const enhancedSongs = songs.map(song => ({
      ...song,
      projects: songProjectMap.get(song.id) || []
    }))
    
    songsWithProjects.value = enhancedSongs
  } catch (err) {
    console.error('Failed to load project information:', err)
    // Fallback to original songs without project info
    songsWithProjects.value = songs.map(song => ({ ...song, projects: [] }))
  }
}

// Watch for changes in song data and load project info
const displayedSongs = computed(() => {
  const songs = searchQuery.value ? searchResults.value : allSongs.value
  if (!songs) return null
  
  return {
    songs: songsWithProjects.value.length > 0 ? songsWithProjects.value : songs.songs,
    count: songs.count
  }
})

// Watch for changes in allSongs and load project info
watch(allSongs, async (newSongs) => {
  if (newSongs?.songs) {
    await loadProjectsForSongs(newSongs.songs)
  }
}, { immediate: true })

// Watch for changes in searchResults and load project info
watch(searchResults, async (newResults) => {
  if (newResults?.songs) {
    await loadProjectsForSongs(newResults.songs)
  }
}, { immediate: true })

// Debounced search function
const debouncedSearch = useDebounceFn(() => {
  // The query will automatically refetch due to reactive dependencies
  // Project info will be loaded by the searchResults watcher
}, 300)

function clearSearch() {
  searchQuery.value = ''
  // Project info will be loaded by the allSongs watcher
}
</script>
