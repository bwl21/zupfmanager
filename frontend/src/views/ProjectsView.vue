<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Projects</h1>
        <p class="mt-2 text-gray-600">Manage your music projects</p>
      </div>
      <button
        @click="showCreateModal = true"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors"
      >
        Create Project
      </button>
    </div>

    <!-- Projects Grid -->
    <div v-if="data?.projects.length" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="project in data.projects"
        :key="project.id"
        class="bg-white rounded-lg shadow hover:shadow-md transition-shadow cursor-pointer"
        @click="$router.push(`/projects/${project.id}`)"
      >
        <div class="p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">{{ project.title }}</h3>
            <div class="flex space-x-2">
              <button
                @click.stop="editProject(project)"
                class="text-gray-400 hover:text-blue-600 transition-colors"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
              </button>
              <button
                @click.stop="deleteProject(project.id)"
                class="text-gray-400 hover:text-red-600 transition-colors"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
          <p class="text-sm text-gray-600 mb-2">Short name: {{ project.short_name }}</p>
          <div class="flex items-center text-sm text-gray-500">
            <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
            Project ID: {{ project.id }}
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="!isLoading" class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No projects</h3>
      <p class="mt-1 text-sm text-gray-500">Get started by creating a new project.</p>
      <div class="mt-6">
        <button
          @click="showCreateModal = true"
          class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors"
        >
          Create Project
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <p class="mt-2 text-gray-600">Loading projects...</p>
    </div>

    <!-- Error State -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4">
      <div class="flex">
        <svg class="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Error loading projects</h3>
          <p class="mt-2 text-sm text-red-700">{{ error.message }}</p>
        </div>
      </div>
    </div>

    <!-- Create/Edit Project Modal -->
    <div v-if="showCreateModal || editingProject" class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900">
            {{ editingProject ? 'Edit Project' : 'Create New Project' }}
          </h3>
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
          <!-- ABC File Directory (only for editing) -->
          <div v-if="editingProject">
            <label for="abc_file_dir" class="block text-sm font-medium text-gray-700">ABC Files Directory</label>
            <div class="mt-1 flex space-x-2">
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
            <!-- Directory selection info -->
            <!-- Path validation warning -->
            <div v-if="projectForm.abc_file_dir_preference && !isValidPath(projectForm.abc_file_dir_preference)" class="mt-2 p-2 bg-orange-50 border border-orange-200 rounded text-xs">
              <div class="flex">
                <svg class="h-4 w-4 text-orange-400 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
                </svg>
                <p class="text-orange-700">
                  <strong>Path may be incomplete:</strong> Please ensure you enter the complete absolute path 
                  (e.g., <code class="bg-orange-100 px-1 rounded">/home/user/music/abc</code> or <code class="bg-orange-100 px-1 rounded">C:\Users\User\Music\ABC</code>).
                </p>
              </div>
            </div>
          </div>
          
          <div>
            <input
              id="default_config"
              v-model="projectForm.default_config"
              type="checkbox"
              class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
            />
            <label for="default_config" class="ml-2 block text-sm text-gray-900">
              Use default configuration
            </label>
          </div>
          <div class="flex justify-end space-x-3 pt-4">
            <button
              type="button"
              @click="cancelEdit"
              class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-md transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="isCreating || isUpdating"
              class="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-md transition-colors disabled:opacity-50"
            >
              {{ isCreating || isUpdating ? 'Saving...' : (editingProject ? 'Update' : 'Create') }}
            </button>
          </div>
        </form>
        
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { projectApi } from '@/services/api'
import type { ProjectResponse, CreateProjectRequest, UpdateProjectRequest } from '@/types/api'

const queryClient = useQueryClient()

// Fetch projects
const { data, isLoading, error } = useQuery({
  queryKey: ['projects'],
  queryFn: projectApi.list
})

// Modal state
const showCreateModal = ref(false)
const editingProject = ref<ProjectResponse | null>(null)

// Form state
const projectForm = reactive({
  title: '',
  short_name: '',
  default_config: true,
  abc_file_dir_preference: ''
})

// Directory picker state

// Create project mutation
const { mutate: createProject, isPending: isCreating } = useMutation({
  mutationFn: (data: CreateProjectRequest) => projectApi.create(data),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['projects'] })
    showCreateModal.value = false
    resetForm()
  }
})

// Update project mutation
const { mutate: updateProject, isPending: isUpdating } = useMutation({
  mutationFn: ({ id, data }: { id: number; data: UpdateProjectRequest }) => projectApi.update(id, data),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['projects'] })
    editingProject.value = null
    resetForm()
  }
})

// Delete project mutation
const { mutate: deleteProjectMutation } = useMutation({
  mutationFn: (id: number) => projectApi.delete(id),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['projects'] })
  }
})

function resetForm() {
  projectForm.title = ''
  projectForm.short_name = ''
  projectForm.default_config = true
  projectForm.abc_file_dir_preference = ''
}

function editProject(project: ProjectResponse) {
  editingProject.value = project
  projectForm.title = project.title
  projectForm.short_name = project.short_name
  projectForm.default_config = false
  projectForm.abc_file_dir_preference = project.abc_file_dir_preference || ''
}

function cancelEdit() {
  showCreateModal.value = false
  editingProject.value = null
  resetForm()
}

async function submitProject() {
  if (editingProject.value) {
    // Update project basic info
    updateProject({
      id: editingProject.value.id,
      data: {
        title: projectForm.title,
        short_name: projectForm.short_name,
        default_config: projectForm.default_config
      }
    })
    
    // Update abc_file_dir_preference separately if it's a valid path
    if (projectForm.abc_file_dir_preference && 
        !projectForm.abc_file_dir_preference.startsWith('[Enter full path to:') &&
        isValidPath(projectForm.abc_file_dir_preference)) {
      try {
        await projectApi.updateAbcFileDir(editingProject.value.id, projectForm.abc_file_dir_preference)
      } catch (err) {
        console.error('Failed to update ABC file directory preference:', err)
      }
    }
  } else {
    createProject({
      title: projectForm.title,
      short_name: projectForm.short_name,
      default_config: projectForm.default_config
    })
  }
}


const isValidPath = (path: string) => {
  if (!path || path.trim() === '') return true // Empty is valid (uses defaults)
  if (path.startsWith('[Enter full path to:')) return false // Placeholder text
  
  // Check for common path patterns
  const hasAbsolutePath = path.startsWith('/') || // Unix/Linux/Mac
                         /^[A-Za-z]:\\/.test(path) || // Windows (C:\)
                         path.startsWith('\\\\') // UNC path (\\server\share)
  
  return hasAbsolutePath && path.length > 3 // Must be more than just root
}

function deleteProject(id: number) {
  if (confirm('Are you sure you want to delete this project?')) {
    deleteProjectMutation(id)
  }
}
</script>
