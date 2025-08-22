<template>
  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
    <div class="relative top-20 mx-auto p-5 border w-11/12 md:w-3/4 lg:w-1/2 shadow-lg rounded-md bg-white">
      <div class="mt-3">
        <!-- Header -->
        <div class="flex items-center justify-between pb-4 border-b">
          <h3 class="text-lg font-medium text-gray-900">Add Song to Project</h3>
          <button
            @click="$emit('close')"
            class="text-gray-400 hover:text-gray-600"
          >
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Song Info -->
        <div class="mt-4 p-4 bg-gray-50 rounded-md">
          <h4 class="text-sm font-medium text-gray-900">{{ song.title }}</h4>
          <p class="text-sm text-gray-500">{{ song.filename }}</p>
          <p v-if="song.genre" class="text-xs text-gray-400 mt-1">{{ song.genre }}</p>
        </div>

        <!-- Select Project -->
        <div class="mt-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Select Project
          </label>
          
          <!-- Loading -->
          <div v-if="isLoadingProjects" class="flex justify-center py-4">
            <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-600"></div>
          </div>
          
          <!-- Projects List -->
          <div v-else-if="availableProjects.length > 0" class="max-h-60 overflow-y-auto border border-gray-200 rounded-md">
            <div
              v-for="project in availableProjects"
              :key="project.id"
              @click="selectProject(project)"
              :class="[
                'p-3 cursor-pointer hover:bg-gray-50 border-b border-gray-200 last:border-b-0',
                selectedProject?.id === project.id ? 'bg-indigo-50 border-indigo-200' : ''
              ]"
            >
              <div class="flex justify-between items-start">
                <div>
                  <p class="text-sm font-medium text-gray-900">{{ project.title }}</p>
                  <p class="text-sm text-gray-500">{{ project.short_name }}</p>
                </div>
                <div v-if="selectedProject?.id === project.id" class="text-indigo-600">
                  <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                </div>
              </div>
            </div>
          </div>
          
          <!-- No Projects -->
          <div v-else class="text-center py-4 text-gray-500">
            No projects found
          </div>
        </div>

        <!-- Song Configuration -->
        <div v-if="selectedProject" class="mt-6 space-y-4 p-4 bg-gray-50 rounded-md">
          <h4 class="text-sm font-medium text-gray-900">Song Configuration</h4>
          
          <!-- Difficulty -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Difficulty
            </label>
            <select
              v-model="songConfig.difficulty"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="easy">Easy</option>
              <option value="medium">Medium</option>
              <option value="hard">Hard</option>
              <option value="expert">Expert</option>
            </select>
          </div>
          
          <!-- Priority -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Priority (1 = highest)
            </label>
            <select
              v-model="songConfig.priority"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option :value="1">1 - Highest</option>
              <option :value="2">2 - High</option>
              <option :value="3">3 - Medium</option>
              <option :value="4">4 - Low</option>
            </select>
          </div>
          
          <!-- Comment -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Comment (optional)
            </label>
            <textarea
              v-model="songConfig.comment"
              rows="2"
              placeholder="Add a comment about this song in the project..."
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            ></textarea>
          </div>
        </div>

        <!-- Actions -->
        <div class="mt-6 flex justify-end gap-4">
          <button
            @click="$emit('close')"
            class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 mr-3"
          >
            Cancel
          </button>
          <button
            @click="addSongToProject"
            :disabled="!selectedProject || isAdding"
            class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="isAdding">Adding...</span>
            <span v-else>Add to Project</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { projectApi, projectSongApi } from '@/services/api'
import type { SongResponse, ProjectResponse, AddSongToProjectRequest } from '@/types/api'

interface Props {
  song: SongResponse
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
  added: []
}>()

// State
const availableProjects = ref<ProjectResponse[]>([])
const selectedProject = ref<ProjectResponse | null>(null)
const isLoadingProjects = ref(false)
const isAdding = ref(false)

const songConfig = ref<AddSongToProjectRequest>({
  difficulty: 'medium',
  priority: 2,
  comment: ''
})

// Methods
const loadProjects = async () => {
  isLoadingProjects.value = true
  try {
    const response = await projectApi.list()
    availableProjects.value = response.projects
  } catch (err) {
    console.error('Failed to load projects:', err)
  } finally {
    isLoadingProjects.value = false
  }
}

const selectProject = (project: ProjectResponse) => {
  selectedProject.value = project
}

const addSongToProject = async () => {
  if (!selectedProject.value) return
  
  isAdding.value = true
  try {
    await projectSongApi.add(selectedProject.value.id, props.song.id, {
      difficulty: songConfig.value.difficulty,
      priority: songConfig.value.priority,
      comment: songConfig.value.comment || undefined
    })
    emit('added')
  } catch (err) {
    alert('Failed to add song to project')
    console.error('Failed to add song:', err)
  } finally {
    isAdding.value = false
  }
}

// Lifecycle
onMounted(() => {
  loadProjects()
})
</script>
