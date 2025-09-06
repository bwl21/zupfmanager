# Zupfmanager Frontend

Modern Vue.js frontend for the Zupfmanager API - a tool for managing music projects and ABC notation files.

## 🚀 Features

- **Vue 3** with Composition API and TypeScript
- **Tailwind CSS** for modern, responsive design
- **TanStack Vue Query** for efficient API state management
- **Vue Router** for client-side routing
- **Vite** for fast development and building
- **API Proxy** - serves both frontend and API from the same port

## 🛠️ Development Setup

### Prerequisites
- Node.js 18+ 
- npm or yarn
- Zupfmanager API server running on port 8080

### Quick Start

```bash
# 1. Start the API server (from project root)
./dist/zupfmanager api --port 8080

# 2. Start the frontend dev server (from frontend directory)
cd frontend
npm install
npm run dev
```

### Access Points
- **Frontend**: http://localhost:5173
- **API (proxied)**: http://localhost:5173/api/v1/...
- **Swagger UI (proxied)**: http://localhost:5173/swagger/index.html

## 🌐 Proxy Configuration

The Vite development server proxies API requests to the Go backend, allowing you to serve both frontend and API from the same port (5173).

### Proxy Routes
- `/api/*` → `http://localhost:8080/api/*`
- `/health` → `http://localhost:8080/health`
- `/swagger/*` → `http://localhost:8080/swagger/*`

### Benefits
- **Single Port**: Access everything from http://localhost:5173
- **No CORS Issues**: Same-origin requests
- **Simplified Development**: One URL for frontend and API
- **Production Ready**: Easy to deploy with reverse proxy

## 📁 Project Structure

```
frontend/
├── src/
│   ├── components/Layout/   # Layout components (Header, etc.)
│   ├── views/              # Page components
│   │   ├── DashboardView.vue    # Main dashboard
│   │   ├── ProjectsView.vue     # Project management
│   │   ├── SongsView.vue        # Song browser & search
│   │   └── ImportView.vue       # File import interface
│   ├── services/api.ts     # Type-safe API client
│   ├── types/api.ts        # TypeScript API types
│   └── router/index.ts     # Vue Router configuration
├── vite.config.ts          # Vite config with proxy setup
└── tailwind.config.js      # Tailwind CSS configuration
```

## 🎨 UI Features

### Dashboard
- Real-time statistics (projects, songs, API status)
- Quick action cards for common tasks
- Health monitoring with auto-refresh

### Project Management
- Full CRUD operations with modal forms
- Project cards with hover effects
- Input validation and error handling

### Song Browser
- Advanced search with filters (title, filename, genre)
- Debounced search for performance
- Detailed song information pages

### Import Interface
- Single file and directory import
- Progress tracking and detailed results
- Quick import for test files

## 🔧 Available Scripts

```bash
# Development
npm run dev          # Start dev server with proxy
npm run build        # Build for production
npm run preview      # Preview production build

# Code Quality
npm run type-check   # TypeScript type checking
npm run lint         # ESLint linting
npm run format       # Prettier formatting

# Testing
npm run test:unit    # Run unit tests
```

## 🚀 Production Deployment

### Option 1: Integrated Server (Recommended)
```bash
# Build the frontend
npm run build

# Start integrated server (API + Frontend from same port)
cd ..
./dist/zupfmanager api --port 8080 --frontend frontend/dist

# Access everything from: http://localhost:8080
```

### Option 2: Separate Deployment
```bash
# Build the frontend
npm run build

# Deploy dist/ directory to static file server
# Configure reverse proxy to route /api/* to Go backend
```

### Option 3: Development with Proxy
```bash
# Terminal 1: Start API
./dist/zupfmanager api --port 8080

# Terminal 2: Start frontend dev server with proxy
cd frontend && npm run dev

# Access from: http://localhost:5173
```

## 💡 Development Tips

1. **API First**: Start the API server before the frontend
2. **Hot Reload**: Changes are reflected instantly during development
3. **Type Safety**: Full TypeScript support with Vue 3 Composition API
4. **Error Handling**: Centralized error handling with user-friendly messages
5. **Responsive**: Mobile-first design with Tailwind CSS

## 🔍 Troubleshooting

### Proxy Not Working
- Ensure API server is running on port 8080
- Check `vite.config.ts` proxy configuration
- Restart the dev server after config changes

### API Errors
- Verify API server is accessible at http://localhost:8080
- Check browser network tab for failed requests
- Review API server logs for errors

### Build Issues
- Run `npm run type-check` to find TypeScript errors
- Ensure all dependencies are installed
- Clear node_modules and reinstall if needed
