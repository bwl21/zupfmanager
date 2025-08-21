<template>
  <div class="space-y-6">
    <!-- Header with Add Song Button -->
    <div class="flex justify-between items-center">
      <h3 class="text-lg font-medium text-gray-900">Songs in Project</h3>
      <button
        @click="showAddSongModal = true"
        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
      >
        <svg class="-ml-1 mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
        </svg>
        Add Song
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="flex justify-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-md p-4">
      <div class="flex">
        <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Error loading project songs</h3>
          <p class="mt-1 text-sm text-red-700">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Songs List -->
    <div v-else-if="projectSongs.length > 0" class="bg-white shadow overflow-hidden sm:rounded-md">
      <ul class="divide-y divide-gray-200">
        <li v-for="projectSong in sortedProjectSongs" :key="projectSong.id" class="px-6 py-4">
          <div class="flex items-center justify-between">
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                 <div class="flex-1">
                   <p class="text-sm font-medium text-gray-900 truncate">
                     {{ projectSong.song?.title || 'Unknown Song' }}
                   </p>
                   <p class="text-sm text-gray-500">
                     {{ projectSong.song?.filename }}
                   </p>
                   
                   <!-- Project Badges (excluding current project) -->
                   <div v-if="getOtherProjects(projectSong.song?.projects).length > 0" class="flex flex-wrap gap-1 mt-1">
                     <button
                       v-for="project in getOtherProjects(projectSong.song?.projects)"
                       :key="project.id"
                       @click="navigateToProject(project.id)"
                       class="inline-flex items-center px-2 py-0.5 text-xs font-medium rounded-full bg-blue-100 text-blue-800 hover:bg-blue-200 transition-colors cursor-pointer"
                       :title="`Go to project: ${project.title}`"
                     >
                       {{ project.short_name }}
                     </button>
                   </div>
                 </div>

                <!-- Difficulty Badge/Select -->
                <div class="relative">
                  <select
                    :value="projectSong.difficulty"
                    @change="updateDifficulty(projectSong, ($event.target as HTMLSelectElement).value)"
                    :class="getDifficultyColor(projectSong.difficulty)"
                    class="appearance-none px-2.5 py-0.5 rounded-full text-xs font-medium border-0 cursor-pointer hover:opacity-80 focus:ring-2 focus:ring-offset-1 focus:ring-blue-500 focus:outline-none"
                  >
                    <option value="easy">easy</option>
                    <option value="medium">medium</option>
                    <option value="hard">hard</option>
                    <option value="expert">expert</option>
                  </select>
                </div>
                
                <!-- Priority Badge/Select -->
                <div class="relative">
                  <select
                    :value="projectSong.priority"
                    @change="updatePriority(projectSong, parseInt(($event.target as HTMLSelectElement).value))"
                    class="appearance-none px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 border-0 cursor-pointer hover:bg-blue-200 focus:ring-2 focus:ring-offset-1 focus:ring-blue-500 focus:outline-none"
                  >
                    <option value="1">1</option>
                    <option value="2">2</option>
                    <option value="3">3</option>
                    <option value="4">4</option>
                  </select>
                </div>
                
                <!-- Action Buttons -->
                <button
                  @click="previewSong(projectSong)"
                  class="inline-flex items-center px-2 py-1 border border-transparent text-xs font-medium rounded text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-1 focus:ring-offset-1 focus:ring-green-500"
                >
                  <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 616 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                  Preview
                </button>
                <button
                  @click="editProjectSong(projectSong)"
                  class="inline-flex items-center px-2 py-1 border border-transparent text-xs font-medium rounded text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-1 focus:ring-offset-1 focus:ring-blue-500"
                >
                  <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                  Edit
                </button>
                <button
                  @click="removeProjectSong(projectSong)"
                  class="inline-flex items-center px-2 py-1 border border-transparent text-xs font-medium rounded text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-1 focus:ring-offset-1 focus:ring-red-500"
                >
                  <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                  Remove
                </button>
              </div>

              <!-- Comment -->
              <p v-if="projectSong.comment" class="mt-2 text-sm text-gray-600">
                {{ projectSong.comment }}
              </p>
            </div>
          </div>
        </li>
      </ul>
    </div>

    <!-- Empty State -->
    <div v-else class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19V6l12-3v13M9 19c0 1.105-.895 2-2 2s-2-.895-2-2 .895-2 2-2 2 .895 2 2zm12-3c0 1.105-.895 2-2 2s-2-.895-2-2 .895-2 2-2 2 .895 2 2z" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No songs in project</h3>
      <p class="mt-1 text-sm text-gray-500">Get started by adding a song to this project.</p>
      <div class="mt-6">
        <button
          @click="showAddSongModal = true"
          class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          <svg class="-ml-1 mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          Add Song
        </button>
      </div>
    </div>

    <!-- Add Song Modal -->
    <AddSongModal
      v-if="showAddSongModal"
      :project-id="projectId"
      @close="showAddSongModal = false"
      @song-added="handleSongAdded"
    />

    <!-- Edit Song Modal -->
    <EditProjectSongModal
      v-if="showEditModal && editingProjectSong"
      :project-song="editingProjectSong"
      @close="showEditModal = false"
      @song-updated="handleSongUpdated"
    />

    <!-- Song Preview Modal -->
    <PreviewModal
      v-if="showPreviewModal && previewingSong && previewingSong.song"
      :song="previewingSong.song"
      :project="project"
      @close="handlePreviewClose"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { projectSongApi, projectApi } from '@/services/api'
import type { ProjectSongResponse, ProjectResponse, SongResponse } from '@/types/api'
import AddSongModal from './AddSongModal.vue'
import EditProjectSongModal from './EditProjectSongModal.vue'
import PreviewModal from './PreviewModal.vue'

interface Props {
  projectId: number
}

const props = defineProps<Props>()
const router = useRouter()

// State
const projectSongs = ref<ProjectSongResponse[]>([])
const project = ref<ProjectResponse | null>(null)
const isLoading = ref(false)
const error = ref<string | null>(null)
const showAddSongModal = ref(false)
const showEditModal = ref(false)
const showPreviewModal = ref(false)
const editingProjectSong = ref<ProjectSongResponse | null>(null)
const previewingSong = ref<ProjectSongResponse | null>(null)

// Computed property for sorted songs
const sortedProjectSongs = computed(() => {
  return [...projectSongs.value].sort((a, b) => {
    const titleA = a.song?.title || 'Unknown Song'
    const titleB = b.song?.title || 'Unknown Song'
    return titleA.localeCompare(titleB, 'de', { sensitivity: 'base' })
  })
})


// Helper functions
const getOtherProjects = (projects?: Array<{id: number, title: string, short_name: string}>) => {
  if (!projects) return []
  return projects.filter(project => project.id !== props.projectId)
}

const navigateToProject = async (projectId: number) => {
  console.log('Navigating to project:', projectId)
  try {
    // Use replace instead of push to force navigation
    await router.replace(`/projects/${projectId}`)
    console.log('Navigation completed')
  } catch (error) {
    console.error('Navigation failed:', error)
    // Fallback to window.location
    window.location.href = `/projects/${projectId}`
  }
}

// Methods
const loadProjectSongs = async () => {
  isLoading.value = true
  error.value = null
  
  try {
    const response = await projectSongApi.list(props.projectId)
    // The backend now includes project associations directly in the response
    projectSongs.value = response.project_songs
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load project songs'
  } finally {
    isLoading.value = false
  }
}

const loadProject = async () => {
  try {
    project.value = await projectApi.get(props.projectId)
  } catch (err) {
    console.error('Failed to load project:', err)
  }
}


const previewSong = (projectSong: ProjectSongResponse) => {
  previewingSong.value = projectSong
  showPreviewModal.value = true
}

const editProjectSong = (projectSong: ProjectSongResponse) => {
  editingProjectSong.value = projectSong
  showEditModal.value = true
}

const removeProjectSong = async (projectSong: ProjectSongResponse) => {
  if (!confirm(`Remove "${projectSong.song?.title}" from this project?`)) {
    return
  }
  
  try {
    await projectSongApi.remove(props.projectId, projectSong.song_id)
    await loadProjectSongs() // Reload the list
  } catch (err) {
    alert('Failed to remove song from project')
  }
}

const handleSongAdded = () => {
  showAddSongModal.value = false
  loadProjectSongs() // Reload the list
}

const handleSongUpdated = () => {
  showEditModal.value = false
  editingProjectSong.value = null
  loadProjectSongs() // Reload the list
}

const handlePreviewClose = () => {
  showPreviewModal.value = false
  previewingSong.value = null
}

const updateDifficulty = async (projectSong: ProjectSongResponse, newDifficulty: string) => {
  try {
    await projectSongApi.update(props.projectId, projectSong.song_id, {
      difficulty: newDifficulty,
      priority: projectSong.priority,
      comment: projectSong.comment
    })
    // Update local state
    projectSong.difficulty = newDifficulty
  } catch (err) {
    console.error('Failed to update difficulty:', err)
    // Optionally show error message to user
  }
}

const updatePriority = async (projectSong: ProjectSongResponse, newPriority: number) => {
  if (isNaN(newPriority) || newPriority < 1 || newPriority > 4) {
    return // Invalid input, don't update
  }
  
  try {
    await projectSongApi.update(props.projectId, projectSong.song_id, {
      difficulty: projectSong.difficulty,
      priority: newPriority,
      comment: projectSong.comment
    })
    // Update local state
    projectSong.priority = newPriority
  } catch (err) {
    console.error('Failed to update priority:', err)
    // Optionally show error message to user
  }
}

const getDifficultyColor = (difficulty: string) => {
  const colors = {
    easy: 'bg-green-100 text-green-800',
    medium: 'bg-yellow-100 text-yellow-800',
    hard: 'bg-orange-100 text-orange-800',
    expert: 'bg-red-100 text-red-800'
  }
  return colors[difficulty as keyof typeof colors] || 'bg-gray-100 text-gray-800'
}

// Lifecycle
onMounted(() => {
  loadProjectSongs()
  loadProject()
})
</script>
