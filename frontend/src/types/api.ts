// API Types generated from OpenAPI specification

export interface ErrorResponse {
  error: string
  message?: string
  details?: Record<string, string>
}

// Project-Song Types
export interface AddSongToProjectRequest {
  difficulty?: number
  priority?: number
}

export interface UpdateProjectSongRequest {
  difficulty?: number
  priority?: number
}

export interface ProjectSongResponse {
  id: number
  project_id: number
  song_id: number
  difficulty?: number
  priority?: number
  song?: SongResponse
  project?: ProjectResponse
}

export interface ProjectSongsResponse {
  project_songs: ProjectSongResponse[]
  total: number
}

// Import Types
export interface ImportFileRequest {
  file_path: string
}

export interface ImportDirectoryRequest {
  directory_path: string
}

export interface ImportResult {
  filename: string
  title: string
  action: 'created' | 'updated' | 'unchanged'
  changes?: string[]
  error?: string
}

export interface ImportSummary {
  total: number
  created: number
  updated: number
  unchanged: number
  errors: number
}

export interface ImportResponse {
  success: boolean
  results: ImportResult[]
  summary: ImportSummary
}

// Project Types
export interface CreateProjectRequest {
  title: string
  short_name: string
  config_file?: string
  default_config?: boolean
}

export interface UpdateProjectRequest {
  title: string
  short_name: string
  config_file?: string
  default_config?: boolean
}

export interface ProjectResponse {
  id: number
  title: string
  short_name: string
  config?: Record<string, any>
}

export interface ProjectListResponse {
  projects: ProjectResponse[]
  count: number
}

// Song Types
export interface SongResponse {
  id: number
  title: string
  filename: string
  genre?: string
  copyright?: string
  tocinfo?: string
}

export interface SongListResponse {
  songs: SongResponse[]
  count: number
}

// Search Options
export interface SearchOptions {
  title?: boolean
  filename?: boolean
  genre?: boolean
}

// Health Check
export interface HealthResponse {
  status: string
  timestamp: string
  version: string
}
