# Zupfmanager Development Guide

Complete development setup and workflow for the Zupfmanager project.

## 🏗️ Architecture Overview

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Vue.js UI     │    │   Vite Proxy     │    │   Go API        │
│   (Port 5173)   │◄──►│   (Port 5173)    │◄──►│   (Port 8080)   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
   Frontend Assets         Proxy Routes              REST API
   - Vue Components        - /api/* → :8080          - Projects CRUD
   - Tailwind CSS          - /health → :8080         - Songs Read
   - TypeScript            - /swagger/* → :8080      - Import Files
   - Vue Query                                       - OpenAPI Docs
```

## 🚀 Quick Start

### 1. Build the API Server
```bash
# From project root
make build
```

### 2. Start Development Servers
```bash
# Terminal 1: Start API server
./dist/zupfmanager api --port 8080

# Terminal 2: Start frontend with proxy
cd frontend
npm install
npm run dev
```

### 3. Access the Application

**Development Mode (with Vite proxy):**
- **Full Application**: http://localhost:5173
- **API Documentation**: http://localhost:5173/swagger/index.html
- **Direct API**: http://localhost:8080 (if needed)

**Production Mode (integrated server):**
```bash
# Build frontend
cd frontend && npm run build

# Start integrated server (API + Frontend from same port)
./dist/zupfmanager api --port 8080 --frontend frontend/dist
```
- **Full Application**: http://localhost:8080
- **API Documentation**: http://localhost:8080/swagger/index.html
- **API Endpoints**: http://localhost:8080/api/v1/...

## 🛠️ Development Workflow

### Frontend Development
```bash
cd frontend

# Install dependencies
npm install

# Start dev server with hot reload
npm run dev

# Type checking
npm run type-check

# Linting and formatting
npm run lint
npm run format

# Build for production
npm run build
```

### Backend Development
```bash
# Build the Go application
make build

# Run with live reload (if using air)
air

# Run tests
go test ./...

# Generate API docs
swag init -g pkg/api/docs.go -o docs/
```

## 🌐 API Integration

### Proxy Configuration
The frontend uses Vite's proxy feature to serve both UI and API from port 5173:

```typescript
// vite.config.ts
server: {
  proxy: {
    '/api': 'http://localhost:8080',
    '/health': 'http://localhost:8080',
    '/swagger': 'http://localhost:8080',
  }
}
```

### API Client
Type-safe API client with automatic error handling:

```typescript
// services/api.ts
import { projectApi, songApi, importApi } from '@/services/api'

// Usage in Vue components
const { data, isLoading, error } = useQuery({
  queryKey: ['projects'],
  queryFn: projectApi.list
})
```

## 📁 Project Structure

```
zupfmanager/
├── cmd/                    # CLI commands
│   ├── api.go             # API server command
│   ├── import.go          # Import command
│   └── project-*.go       # Project commands
├── pkg/
│   ├── api/               # REST API implementation
│   │   ├── handlers/      # HTTP handlers
│   │   ├── models/        # API models
│   │   └── server.go      # Server setup
│   └── core/              # Business logic
│       ├── interfaces.go  # Service interfaces
│       ├── project.go     # Project service
│       ├── song.go        # Song service
│       └── import.go      # Import service
├── internal/
│   ├── database/          # Database client
│   ├── ent/              # ORM generated code
│   └── ui/               # TUI components
├── frontend/              # Vue.js frontend
│   ├── src/
│   │   ├── components/    # Vue components
│   │   ├── views/         # Page components
│   │   ├── services/      # API client
│   │   └── types/         # TypeScript types
│   ├── vite.config.ts     # Vite configuration
│   └── tailwind.config.js # Tailwind CSS config
├── docs/                  # Generated API docs
├── test_songs/           # Sample ABC files
└── dist/                 # Built binaries
```

## 🔧 Configuration

### Environment Variables
```bash
# Frontend (.env.development)
VITE_API_BASE_URL=          # Empty for proxy

# Backend
ZUPFNOTER_PATH=             # Path to zupfnoter CLI
ZUPFNOTER_DEBUG=1           # Enable debug mode
```

### Database
- **Type**: SQLite
- **File**: `zupfmanager.db` (auto-created)
- **ORM**: Ent (Facebook's entity framework)
- **Migrations**: Automatic on startup

## 🧪 Testing

### Frontend Testing
```bash
cd frontend

# Unit tests
npm run test:unit

# E2E tests (if configured)
npm run test:e2e
```

### Backend Testing
```bash
# Run all tests
go test ./...

# Test specific package
go test ./pkg/core

# Test with coverage
go test -cover ./...

# Integration tests
go test ./pkg/api/handlers
```

### API Testing
```bash
# Health check
curl http://localhost:5173/health

# List projects
curl http://localhost:5173/api/v1/projects

# Create project
curl -X POST http://localhost:5173/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{"title": "Test", "short_name": "test", "default_config": true}'

# Import test songs
curl -X POST http://localhost:5173/api/v1/import/directory \
  -H "Content-Type: application/json" \
  -d '{"directory_path": "/workspaces/zupfmanager/test_songs/"}'
```

## 📝 Recent Changes

### 2025-08-21: Performance Optimization for Project Loading
- **Fixed slow search performance** in SongsView caused by inefficient project loading
- **Optimized project-song relationship loading** from O(N×M) to O(M) API calls
- **Implemented parallel processing** for loading all project-song relationships simultaneously
- **Added efficient mapping** using Map data structure for O(1) lookups
- **Applied same optimization** to AddSongModal for consistent performance
- **Improved error handling** so individual project failures don't break the entire process

**Technical Details:**
- Changed from sequential API calls for each song-project combination
- Now loads all projects first, then fetches all project-song relationships in parallel
- Uses a Map to efficiently associate songs with their projects
- Reduces API calls from potentially hundreds to just the number of projects

### 2025-08-21: Project Association Display
- Added project badges to song views (SongsView and AddSongModal)
- Project badges are clickable and navigate to the respective project
- Enhanced song data loading to include project associations

### 2025-08-21: Inline Editing Features
- Implemented inline editing for difficulty (dropdown) and priority (dropdown 1-4) in ProjectSongManager
- Added auto-save functionality for immediate persistence
- Unified spacing for all interactive elements (difficulty, priority, preview, edit, remove buttons)
- Enhanced user experience with immediate feedback

### 2025-08-21: Song Sorting and Preview Enhancements
- Added alphabetical sorting by title for songs in project view with German locale support
- Moved preview functionality from project-level to individual song-level
- Integrated project's abc_file_dir_preference with song previews
- Auto-initializes directory path and searches for PDFs when project directory is available

## 🚀 Deployment

### Production Build
```bash
# Build backend
make build-all

# Build frontend
cd frontend
npm run build

# The dist/ directory contains static files
# Serve with nginx, Apache, or any static file server
```

### Docker Deployment (Future)
```dockerfile
# Multi-stage build
FROM golang:1.23 AS backend
# ... build Go binary

FROM node:18 AS frontend
# ... build Vue.js app

FROM alpine:latest
# ... combine both
```

## 🐛 Debugging

### Frontend Debugging
- **Vue DevTools**: Available in browser
- **Network Tab**: Check API requests
- **Console**: Vue Query DevTools
- **Hot Reload**: Instant feedback

### Backend Debugging
- **Logs**: Check API server output
- **Database**: SQLite browser for `zupfmanager.db`
- **API Docs**: Use Swagger UI for testing
- **Postman**: Import OpenAPI spec

### Common Issues

1. **Proxy Not Working**
   - Ensure API server is running on port 8080
   - Restart frontend dev server after config changes

2. **CORS Errors**
   - Should not occur with proxy setup
   - Check if API server has CORS enabled

3. **Build Failures**
   - Run `npm run type-check` for TypeScript errors
   - Check Go build with `go build`

4. **Database Issues**
   - Delete `zupfmanager.db` to reset
   - Check file permissions

## 📚 Resources

### Documentation
- **Vue 3**: https://vuejs.org/
- **Vite**: https://vitejs.dev/
- **Tailwind CSS**: https://tailwindcss.com/
- **TanStack Query**: https://tanstack.com/query
- **Go**: https://golang.org/
- **Ent ORM**: https://entgo.io/

### Tools
- **VS Code Extensions**: Vue Language Features (Volar), Go
- **Browser DevTools**: Vue DevTools extension
- **API Testing**: Postman, Insomnia, or built-in Swagger UI
- **Database**: SQLite Browser, TablePlus

## 🤝 Contributing

1. **Code Style**: Follow existing patterns
2. **TypeScript**: Add types for new features
3. **Testing**: Write tests for new functionality
4. **Documentation**: Update README and comments
5. **API**: Update OpenAPI spec for new endpoints

### Git Workflow
```bash
# Create feature branch
git checkout -b feature/new-feature

# Make changes and commit
git add .
git commit -m "feat: add new feature"

# Push and create PR
git push origin feature/new-feature
```
