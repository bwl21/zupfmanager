import axios from 'axios'
import type {
  ProjectResponse,
  ProjectListResponse,
  CreateProjectRequest,
  UpdateProjectRequest,
  SongResponse,
  SongListResponse,
  ImportFileRequest,
  ImportDirectoryRequest,
  ImportResponse,
  HealthResponse,
  ErrorResponse
} from '@/types/api'

// Create axios instance with base configuration
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Add response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.data) {
      throw error.response.data as ErrorResponse
    }
    throw new Error(error.message || 'An unexpected error occurred')
  }
)

// Health Check API
export const healthApi = {
  check: (): Promise<HealthResponse> => api.get('/health').then((res) => res.data)
}

// Project API
export const projectApi = {
  list: (): Promise<ProjectListResponse> => api.get('/api/v1/projects').then((res) => res.data),

  get: (id: number): Promise<ProjectResponse> =>
    api.get(`/api/v1/projects/${id}`).then((res) => res.data),

  create: (data: CreateProjectRequest): Promise<ProjectResponse> =>
    api.post('/api/v1/projects', data).then((res) => res.data),

  update: (id: number, data: UpdateProjectRequest): Promise<ProjectResponse> =>
    api.put(`/api/v1/projects/${id}`, data).then((res) => res.data),

  delete: (id: number): Promise<void> => api.delete(`/api/v1/projects/${id}`)
}

// Song API
export const songApi = {
  list: (): Promise<SongListResponse> => api.get('/api/v1/songs').then((res) => res.data),

  get: (id: number): Promise<SongResponse> =>
    api.get(`/api/v1/songs/${id}`).then((res) => res.data),

  search: (query: string, options?: { title?: boolean; filename?: boolean; genre?: boolean }): Promise<SongListResponse> => {
    const params = new URLSearchParams({ q: query })
    if (options?.title !== undefined) params.append('title', options.title.toString())
    if (options?.filename !== undefined) params.append('filename', options.filename.toString())
    if (options?.genre !== undefined) params.append('genre', options.genre.toString())
    
    return api.get(`/api/v1/songs/search?${params}`).then((res) => res.data)
  }
}

// Import API
export const importApi = {
  file: (data: ImportFileRequest): Promise<ImportResponse> =>
    api.post('/api/v1/import/file', data).then((res) => res.data),

  directory: (data: ImportDirectoryRequest): Promise<ImportResponse> =>
    api.post('/api/v1/import/directory', data).then((res) => res.data)
}

export default api
