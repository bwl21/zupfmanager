<template>
  <header class="bg-white shadow-sm border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo and Title -->
        <div class="flex items-center space-x-6">
          <div class="flex-shrink-0">
            <div class="flex flex-col">
              <h1 class="text-2xl font-bold text-gray-900">🎵 Zupfmanager</h1>
              <VersionInfo />
            </div>
          </div>
          
          <!-- Working Directory -->
          <div class="hidden md:flex items-center text-sm text-gray-500 border-l border-gray-200 pl-6 h-auto min-h-8">
            <div class="relative group">
              <div 
                class="max-w-lg whitespace-normal break-all line-clamp-2 overflow-hidden text-ellipsis cursor-pointer hover:bg-gray-100 px-2 py-1 rounded transition-colors"
                @click="openInFileExplorer"
              >
                <span class="font-mono text-xs">📁 {{ workingDir || 'Lade Verzeichnis...' }}</span>
              </div>
              <div class="absolute left-0 top-full mt-1 p-2 bg-white border border-gray-200 rounded shadow-lg z-10 max-w-2xl break-all whitespace-pre-wrap hidden group-hover:block text-xs font-mono">
                {{ workingDir }}
                <div class="mt-1 text-blue-500 text-2xs">Klicken, um im Finder zu öffnen</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Navigation -->
        <nav class="hidden md:flex space-x-8">
          <RouterLink
            to="/"
            class="text-gray-500 hover:text-gray-700 px-3 py-2 rounded-md text-sm font-medium transition-colors"
            active-class="text-blue-600 bg-blue-50"
          >
            Dashboard
          </RouterLink>
          <RouterLink
            to="/projects"
            class="text-gray-500 hover:text-gray-700 px-3 py-2 rounded-md text-sm font-medium transition-colors"
            active-class="text-blue-600 bg-blue-50"
          >
            Projects
          </RouterLink>
          <RouterLink
            to="/songs"
            class="text-gray-500 hover:text-gray-700 px-3 py-2 rounded-md text-sm font-medium transition-colors"
            active-class="text-blue-600 bg-blue-50"
          >
            Songs
          </RouterLink>
          <RouterLink
            to="/import"
            class="text-gray-500 hover:text-gray-700 px-3 py-2 rounded-md text-sm font-medium transition-colors"
            active-class="text-blue-600 bg-blue-50"
          >
            Import
          </RouterLink>
        </nav>

        <!-- Mobile menu button -->
        <div class="md:hidden">
          <button
            @click="mobileMenuOpen = !mobileMenuOpen"
            class="text-gray-500 hover:text-gray-700 focus:outline-none focus:text-gray-700"
          >
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 6h16M4 12h16M4 18h16"
              />
            </svg>
          </button>
        </div>
      </div>

      <!-- Mobile menu -->
      <div v-if="mobileMenuOpen" class="md:hidden">
        <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3">
          <RouterLink
            to="/"
            class="text-gray-500 hover:text-gray-700 block px-3 py-2 rounded-md text-base font-medium"
            active-class="text-blue-600 bg-blue-50"
            @click="mobileMenuOpen = false"
          >
            Dashboard
          </RouterLink>
          <RouterLink
            to="/projects"
            class="text-gray-500 hover:text-gray-700 block px-3 py-2 rounded-md text-base font-medium"
            active-class="text-blue-600 bg-blue-50"
            @click="mobileMenuOpen = false"
          >
            Projects
          </RouterLink>
          <RouterLink
            to="/songs"
            class="text-gray-500 hover:text-gray-700 block px-3 py-2 rounded-md text-base font-medium"
            active-class="text-blue-600 bg-blue-50"
            @click="mobileMenuOpen = false"
          >
            Songs
          </RouterLink>
          <RouterLink
            to="/import"
            class="text-gray-500 hover:text-gray-700 block px-3 py-2 rounded-md text-base font-medium"
            active-class="text-blue-600 bg-blue-50"
            @click="mobileMenuOpen = false"
          >
            Import
          </RouterLink>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import VersionInfo from '../VersionInfo.vue'

const workingDir = ref('')
const isOpening = ref(false)

const openInFileExplorer = async () => {
  if (!workingDir.value) return
  
  try {
    isOpening.value = true
    const response = await fetch('/api/open-directory', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    })
    
    if (!response.ok) {
      const error = await response.json()
      console.error('Failed to open directory:', error)
    }
  } catch (err) {
    console.error('Error opening directory:', err)
  } finally {
    isOpening.value = false
  }
}

onMounted(async () => {
  try {
    const response = await fetch('/api/version')
    if (response.ok) {
      const data = await response.json()
      workingDir.value = data.working_dir || ''
    }
  } catch (err) {
    console.error('Error fetching working directory:', err)
  }
})

const mobileMenuOpen = ref(false)
</script>
