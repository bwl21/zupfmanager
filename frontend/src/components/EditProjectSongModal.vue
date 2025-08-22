<template>
  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
    <div class="relative top-20 mx-auto p-5 border w-11/12 md:w-2/3 lg:w-1/2 shadow-lg rounded-md bg-white">
      <div class="mt-3">
        <!-- Header -->
        <div class="flex items-center justify-between pb-4 border-b">
          <h3 class="text-lg font-medium text-gray-900">Edit Song Settings</h3>
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
          <h4 class="text-sm font-medium text-gray-900">{{ projectSong.song?.title }}</h4>
          <p class="text-sm text-gray-500">{{ projectSong.song?.filename }}</p>
          <p v-if="projectSong.song?.genre" class="text-xs text-gray-400 mt-1">{{ projectSong.song?.genre }}</p>
        </div>

        <!-- Configuration Form -->
        <div class="mt-6 space-y-4">
          <!-- Difficulty -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Difficulty
            </label>
            <select
              v-model="formData.difficulty"
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
              v-model="formData.priority"
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
              v-model="formData.comment"
              rows="3"
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
            @click="updateProjectSong"
            :disabled="isUpdating"
            class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="isUpdating">Updating...</span>
            <span v-else>Update Song</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { projectSongApi } from '@/services/api'
import type { ProjectSongResponse, UpdateProjectSongRequest } from '@/types/api'

interface Props {
  projectSong: ProjectSongResponse
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
  songUpdated: []
}>()

// State
const isUpdating = ref(false)
const formData = ref<UpdateProjectSongRequest>({
  difficulty: props.projectSong.difficulty,
  priority: props.projectSong.priority,
  comment: props.projectSong.comment || ''
})

// Methods
const updateProjectSong = async () => {
  isUpdating.value = true
  try {
    await projectSongApi.update(
      props.projectSong.project_id,
      props.projectSong.song_id,
      {
        difficulty: formData.value.difficulty,
        priority: formData.value.priority,
        comment: formData.value.comment || undefined
      }
    )
    emit('songUpdated')
  } catch (err) {
    alert('Failed to update song settings')
    console.error('Failed to update song:', err)
  } finally {
    isUpdating.value = false
  }
}
</script>
