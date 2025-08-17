<template>
  <header class="bg-white shadow-sm border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo and Title -->
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="flex flex-col">
              <h1 class="text-2xl font-bold text-gray-900">🎵 Zupfmanager</h1>
              <VersionInfo />
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
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import VersionInfo from '../VersionInfo.vue'

const mobileMenuOpen = ref(false)
</script>
