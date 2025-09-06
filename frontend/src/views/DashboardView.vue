<template>
  <div class="space-y-6">
    <!-- Page Header -->
    <div>
      <h1 class="text-3xl font-bold text-gray-900">Dashboard</h1>
      <p class="mt-2 text-gray-600">Welcome to Zupfmanager - manage your music projects and ABC notation files</p>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <!-- Projects Card -->
      <div class="bg-white rounded-lg shadow p-6">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-blue-500 rounded-md flex items-center justify-center">
              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">Total Projects</p>
            <p class="text-2xl font-semibold text-gray-900">{{ projectsData?.count || 0 }}</p>
          </div>
        </div>
        <div class="mt-4">
          <RouterLink to="/projects" class="text-blue-600 hover:text-blue-500 text-sm font-medium">
            View all projects →
          </RouterLink>
        </div>
      </div>

      <!-- Songs Card -->
      <div class="bg-white rounded-lg shadow p-6">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-green-500 rounded-md flex items-center justify-center">
              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
              </svg>
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">Total Songs</p>
            <p class="text-2xl font-semibold text-gray-900">{{ songsData?.count || 0 }}</p>
          </div>
        </div>
        <div class="mt-4">
          <RouterLink to="/songs" class="text-green-600 hover:text-green-500 text-sm font-medium">
            Browse songs →
          </RouterLink>
        </div>
      </div>

      <!-- Health Status Card -->
      <div class="bg-white rounded-lg shadow p-6">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-purple-500 rounded-md flex items-center justify-center">
              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">API Status</p>
            <p class="text-2xl font-semibold" :class="healthData?.status === 'ok' ? 'text-green-600' : 'text-red-600'">
              {{ healthData?.status === 'ok' ? 'Online' : 'Offline' }}
            </p>
          </div>
        </div>
        <div class="mt-4">
          <p class="text-sm text-gray-500">
            Version: {{ healthData?.version || 'Unknown' }}
          </p>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="bg-white rounded-lg shadow">
      <div class="px-6 py-4 border-b border-gray-200">
        <h2 class="text-lg font-medium text-gray-900">Quick Actions</h2>
      </div>
      <div class="p-6">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <RouterLink
            to="/projects"
            class="flex items-center p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
          >
            <div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center mr-4">
              <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
            </div>
            <div>
              <p class="font-medium text-gray-900">New Project</p>
              <p class="text-sm text-gray-500">Create a music project</p>
            </div>
          </RouterLink>

          <RouterLink
            to="/import"
            class="flex items-center p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
          >
            <div class="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center mr-4">
              <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
              </svg>
            </div>
            <div>
              <p class="font-medium text-gray-900">Import Songs</p>
              <p class="text-sm text-gray-500">Add ABC files</p>
            </div>
          </RouterLink>

          <RouterLink
            to="/songs"
            class="flex items-center p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
          >
            <div class="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center mr-4">
              <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <div>
              <p class="font-medium text-gray-900">Search Songs</p>
              <p class="text-sm text-gray-500">Find music files</p>
            </div>
          </RouterLink>

          <a
            href="/swagger/index.html"
            target="_blank"
            class="flex items-center p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
          >
            <div class="w-10 h-10 bg-orange-100 rounded-lg flex items-center justify-center mr-4">
              <svg class="w-6 h-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
              </svg>
            </div>
            <div>
              <p class="font-medium text-gray-900">API Docs</p>
              <p class="text-sm text-gray-500">Swagger UI</p>
            </div>
          </a>
        </div>
      </div>
    </div>

    <!-- Loading States -->
    <div v-if="isLoadingProjects || isLoadingSongs || isLoadingHealth" class="text-center py-8">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <p class="mt-2 text-gray-600">Loading dashboard data...</p>
    </div>

    <!-- Error States -->
    <div v-if="projectsError || songsError || healthError" class="bg-red-50 border border-red-200 rounded-lg p-4">
      <div class="flex">
        <svg class="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Error loading dashboard data</h3>
          <div class="mt-2 text-sm text-red-700">
            <p v-if="projectsError">Projects: {{ projectsError.message }}</p>
            <p v-if="songsError">Songs: {{ songsError.message }}</p>
            <p v-if="healthError">Health: {{ healthError.message }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { RouterLink } from 'vue-router'
import { projectApi, songApi, healthApi } from '@/services/api'

// Fetch dashboard data
const { data: projectsData, isLoading: isLoadingProjects, error: projectsError } = useQuery({
  queryKey: ['projects'],
  queryFn: projectApi.list
})

const { data: songsData, isLoading: isLoadingSongs, error: songsError } = useQuery({
  queryKey: ['songs'],
  queryFn: songApi.list
})

const { data: healthData, isLoading: isLoadingHealth, error: healthError } = useQuery({
  queryKey: ['health'],
  queryFn: healthApi.check,
  refetchInterval: 30000 // Refetch every 30 seconds
})
</script>
