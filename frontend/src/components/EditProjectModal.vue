<template>
  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50" @click="$emit('close')">
    <div class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4" @click.stop>
      <div class="px-6 py-4 border-b border-gray-200">
        <h3 class="text-lg font-medium text-gray-900">Edit Project</h3>
      </div>
      <form @submit.prevent="submitProject" class="p-6 space-y-4">
        <div>
          <label for="title" class="block text-sm font-medium text-gray-700">Title</label>
          <input
            id="title"
            v-model="projectForm.title"
            type="text"
            required
            class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            placeholder="My Music Project"
          />
        </div>
        <div>
          <label for="short_name" class="block text-sm font-medium text-gray-700">Short Name</label>
          <input
            id="short_name"
            v-model="projectForm.short_name"
            type="text"
            required
            pattern="[a-zA-Z0-9_-]+"
            class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            placeholder="my-project"
          />
          <p class="mt-1 text-xs text-gray-500">Only letters, numbers, hyphens, and underscores allowed</p>
        </div>
        <!-- ABC File Directory -->
        <div>
          <label for="abc_file_dir" class="block text-sm font-medium text-gray-700">ABC Files Directory</label>
          <div class="mt-1 flex gap-4">
            <input
              id="abc_file_dir"
              v-model="projectForm.abc_file_dir_preference"
              type="text"
              class="flex-1 border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder="Full path to ABC files directory (e.g., /home/user/music/abc)"
            />
          </div>
          <p class="mt-1 text-xs text-gray-500">
            Full path to directory containing ABC notation files.
          </p>
          <!-- Path validation warning -->
          <div v-if="projectForm.abc_file_dir_preference && !isValidPath(projectForm.abc_file_dir_preference)" class="mt-2 p-2 bg-orange-50 border border-orange-200 rounded text-xs">
            <div class="flex">
              <svg class="h-4 w-4 text-orange-400 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
              </svg>
              <p class="text-orange-700">
                <strong>Path may be incomplete:</strong> Please ensure you enter the complete absolute path
              </p>
            </div>
          </div>
        </div>
        <div class="flex justify-end gap-4 pt-4">
          <button
            type="button"
            @click="$emit('close')"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 mr-3"
          >
            Cancel
          </button>
          <button
            type="submit"
            :disabled="isUpdating"
            class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 disabled:opacity-50"
          >
            {{ isUpdating ? 'Updating...' : 'Update Project' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, onMounted, ref } from 'vue'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { projectApi } from '@/services/api'
import type { ProjectResponse, UpdateProjectRequest } from '@/types/api'

interface Props {
  project: ProjectResponse
}

const props = defineProps<Props>()

const emit = defineEmits<{
  close: []
  updated: [project: ProjectResponse]
}>()

const queryClient = useQueryClient()

// Form state
const projectForm = reactive({
  title: '',
  short_name: '',
  abc_file_dir_preference: ''
})

// Update project mutation
const { mutate: updateProject, isPending: isUpdating } = useMutation({
  mutationFn: ({ id, data }: { id: number; data: UpdateProjectRequest }) => projectApi.update(id, data),
  onSuccess: (updatedProject) => {
    queryClient.invalidateQueries({ queryKey: ['projects'] })
    queryClient.invalidateQueries({ queryKey: ['projects', props.project.id] })
    emit('updated', updatedProject)
    emit('close')
  }
})

// Initialize form
onMounted(() => {
  projectForm.title = props.project.title
  projectForm.short_name = props.project.short_name
  projectForm.abc_file_dir_preference = props.project.abc_file_dir_preference || ''
})

// Path validation
const isValidPath = (path: string) => {
  return path.startsWith('/') || path.match(/^[A-Za-z]:\\/)
}

async function submitProject() {
  const updateData: UpdateProjectRequest = {
    title: projectForm.title,
    short_name: projectForm.short_name
  }

  // Include abc_file_dir_preference if it's provided and valid
  if (projectForm.abc_file_dir_preference && isValidPath(projectForm.abc_file_dir_preference)) {
    updateData.abc_file_dir_preference = projectForm.abc_file_dir_preference
  }

  updateProject({
    id: props.project.id,
    data: updateData
  })
}
</script>
