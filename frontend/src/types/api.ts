// API Types generated from OpenAPI specification

export interface ErrorResponse {
  error: string
  message?: string
  details?: Record<string, string>
}

// Preview Types
export interface GeneratePreviewRequest {
  abc_file_dir: string
  config?: Record<string, any>
}

export interface GeneratePreviewResponse {
  pdf_files: string[]
  preview_dir: string
}

export interface PreviewPDFResponse {
  filename: string
  size: number
  created_at: string
}

export interface PreviewPDFListResponse {
  pdfs: PreviewPDFResponse[]
  count: number
}

export interface MessageResponse {
  message: string
}

// Project-Song Types
export interface AddSongToProjectRequest {
  difficulty?: string // "easy" | "medium" | "hard" | "expert"
  priority?: number // 1-4
  comment?: string
}

export interface UpdateProjectSongRequest {
  difficulty?: string // "easy" | "medium" | "hard" | "expert"
  priority?: number // 1-4
  comment?: string
}

export interface ProjectSongResponse {
  id: number
  project_id: number
  song_id: number
  difficulty: string // "easy" | "medium" | "hard" | "expert"
  priority: number // 1-4
  comment?: string
  song?: SongResponse
  project?: ProjectResponse
}

export interface ProjectSongsResponse {
  project_songs: ProjectSongResponse[]
  total: number
}

// Project Build Types
export interface BuildProjectRequest {
  output_dir?: string
  abc_file_dir?: string
  priority_threshold?: number // 1-4
  sample_id?: string
}

export interface BuildStatusResponse {
  status: string // "pending" | "running" | "completed" | "failed"
  progress: number // 0-100
  message?: string
  started_at?: string
  completed_at?: string
  error?: string
}

export interface BuildResultResponse {
  build_id: string
  project_id: number
  status: string // "pending" | "running" | "completed" | "failed"
  output_dir: string
  generated_files?: string[]
  started_at: string
  completed_at?: string
  error?: string
}

export interface BuildListResponse {
  builds: BuildResultResponse[]
  total: number
}

export interface BuildDefaultsResponse {
  output_dir: string
  abc_file_dir: string
  priority_threshold: number
  sample_id: string
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
  abc_file_dir_preference?: string
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
  projects?: Array<{
    id: number
    title: string
    short_name: string
  }>
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
