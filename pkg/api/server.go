package api

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/bwl21/zupfmanager/docs"
	"github.com/bwl21/zupfmanager/pkg/api/handlers"
	"github.com/bwl21/zupfmanager/pkg/core"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server represents the API server
type Server struct {
	router   *gin.Engine
	services *core.Services
	server   *http.Server
	
	// Handlers
	importHandler      *handlers.ImportHandler
	projectHandler     *handlers.ProjectHandler
	songHandler        *handlers.SongHandler
	projectSongHandler *handlers.ProjectSongHandler

	// Frontend serving
	frontendPath string
	frontendFS   fs.FS
	useEmbedded  bool
	
	// Version info
	version   string
	gitCommit string
}

// ServerOptions configures the server
type ServerOptions struct {
	FrontendPath string // Path to frontend dist directory (fallback)
	UseEmbedded  bool   // Use embedded frontend files
	Version      string // Version string
	GitCommit    string // Git commit hash
}

// NewServer creates a new API server
func NewServer(services *core.Services, opts ...ServerOptions) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	
	var frontendPath, version, gitCommit string
	var useEmbedded bool
	if len(opts) > 0 {
		frontendPath = opts[0].FrontendPath
		useEmbedded = opts[0].UseEmbedded
		version = opts[0].Version
		gitCommit = opts[0].GitCommit
	}
	
	s := &Server{
		router:             router,
		services:           services,
		importHandler:      handlers.NewImportHandler(services),
		projectHandler:     handlers.NewProjectHandler(services),
		songHandler:        handlers.NewSongHandler(services),
		projectSongHandler: handlers.NewProjectSongHandler(services),
		frontendPath:       frontendPath,
		useEmbedded:        useEmbedded,
		version:            version,
		gitCommit:          gitCommit,
	}

	s.setupRoutes()
	return s
}

// Router returns the gin router for testing
func (s *Server) Router() *gin.Engine {
	return s.router
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.healthCheck)
	
	// Version info
	s.router.GET("/api/version", s.versionInfo)
	
	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Import endpoints
		v1.POST("/import/file", s.importHandler.ImportFile)
		v1.POST("/import/directory", s.importHandler.ImportDirectory)
		v1.GET("/import/last-path", s.importHandler.GetLastImportPath)

		// Debug: Check if handlers are initialized
		if s.projectHandler == nil {
			slog.Error("Project handler is nil!")
		}
		if s.songHandler == nil {
			slog.Error("Song handler is nil!")
		}
		
		// Project endpoints
		projects := v1.Group("/projects")
		{
			projects.GET("", s.projectHandler.ListProjects)
			projects.POST("", s.projectHandler.CreateProject)
			projects.GET("/default-config", s.projectHandler.GetDefaultConfig)
			projects.GET("/:id", s.projectHandler.GetProject)
			projects.PUT("/:id", s.projectHandler.UpdateProject)
			projects.DELETE("/:id", s.projectHandler.DeleteProject)
			projects.PUT("/:id/abc-file-dir", s.projectHandler.UpdateAbcFileDirPreference)

			// Project-Song relationship endpoints
			projects.GET("/:id/songs", s.projectSongHandler.ListProjectSongs)
			projects.POST("/:id/songs/:songId", s.projectSongHandler.AddSongToProject)
			projects.PUT("/:id/songs/:songId", s.projectSongHandler.UpdateProjectSong)
			projects.DELETE("/:id/songs/:songId", s.projectSongHandler.RemoveSongFromProject)
			
			// Project build endpoints
			projects.POST("/:id/build", s.projectHandler.BuildProject)
			projects.GET("/:id/build/defaults", s.projectHandler.GetBuildDefaults)
			projects.GET("/:id/builds", s.projectHandler.ListBuilds)
			projects.DELETE("/:id/builds", s.projectHandler.ClearBuildHistory)
			projects.GET("/:id/builds/:buildId/status", s.projectHandler.GetBuildStatus)
		}

		// Song endpoints
		songs := v1.Group("/songs")
		{
			songs.GET("", s.songHandler.ListSongs)
			songs.GET("/:id", s.songHandler.GetSong)
			songs.GET("/search", s.songHandler.SearchSongs)
			songs.DELETE("/:id", s.songHandler.DeleteSong)
			
			// Preview endpoints
			songs.POST("/:id/generate-preview", s.songHandler.GeneratePreview)
			songs.GET("/:id/preview-pdfs", s.songHandler.ListPreviewPDFs)
			songs.GET("/:id/preview-pdf/:filename", s.songHandler.GetPreviewPDF)
			songs.DELETE("/:id/preview-pdfs", s.songHandler.CleanupPreviewPDFs)
		}
	}
	
	// Dynamic swagger config that adapts to current host/scheme
	s.router.GET("/api/swagger.json", s.swaggerConfig)
	
	// Swagger documentation with dynamic URL
	url := ginSwagger.URL("/api/swagger.json")
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	
	// Serve frontend - embedded takes priority over external path
	if s.useEmbedded {
		if err := s.setupEmbeddedFrontend(); err != nil {
			slog.Info("Embedded frontend not available, using external files", "reason", err.Error())
			if s.frontendPath != "" {
				s.setupFrontendServing()
			}
		}
	} else if s.frontendPath != "" {
		s.setupFrontendServing()
	}
}

// setupFrontendServing configures static file serving for the frontend
func (s *Server) setupFrontendServing() {
	// Check if frontend dist directory exists
	if _, err := os.Stat(s.frontendPath); os.IsNotExist(err) {
		slog.Warn("Frontend path does not exist", "path", s.frontendPath)
		return
	}
	
	slog.Info("Serving frontend static files", "path", s.frontendPath)
	
	// Serve static assets with proper MIME types
	s.router.GET("/assets/*filepath", s.serveStaticWithMimeType)
	
	// Serve favicon and other root files with proper MIME types
	s.router.GET("/favicon.ico", func(c *gin.Context) {
		s.serveFileWithMimeType(c, filepath.Join(s.frontendPath, "favicon.ico"))
	})
	s.router.GET("/vite.svg", func(c *gin.Context) {
		s.serveFileWithMimeType(c, filepath.Join(s.frontendPath, "vite.svg"))
	})
	
	// SPA fallback: serve index.html for all non-API routes
	s.router.NoRoute(s.serveSPA)
}

// serveStaticWithMimeType serves static files with correct MIME types
func (s *Server) serveStaticWithMimeType(c *gin.Context) {
	// Get the file path from the URL
	filePath := c.Param("filepath")
	fullPath := filepath.Join(s.frontendPath, "assets", filePath)
	
	s.serveFileWithMimeType(c, fullPath)
}

// serveFileWithMimeType serves a file with the correct MIME type
func (s *Server) serveFileWithMimeType(c *gin.Context, filePath string) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	
	// Get MIME type based on file extension
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	
	// Set specific MIME types for common web assets
	switch ext {
	case ".css":
		mimeType = "text/css; charset=utf-8"
	case ".js":
		mimeType = "application/javascript; charset=utf-8"
	case ".json":
		mimeType = "application/json; charset=utf-8"
	case ".svg":
		mimeType = "image/svg+xml"
	case ".ico":
		mimeType = "image/x-icon"
	case ".png":
		mimeType = "image/png"
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".gif":
		mimeType = "image/gif"
	case ".woff":
		mimeType = "font/woff"
	case ".woff2":
		mimeType = "font/woff2"
	case ".ttf":
		mimeType = "font/ttf"
	case ".eot":
		mimeType = "application/vnd.ms-fontobject"
	}
	
	// Set the Content-Type header
	if mimeType != "" {
		c.Header("Content-Type", mimeType)
	}
	
	// Set cache headers for static assets
	if strings.HasPrefix(c.Request.URL.Path, "/assets/") {
		c.Header("Cache-Control", "public, max-age=31536000") // 1 year
	} else {
		c.Header("Cache-Control", "public, max-age=3600") // 1 hour
	}
	
	// Serve the file
	c.File(filePath)
}

// serveSPA serves the Single Page Application for client-side routing
func (s *Server) serveSPA(c *gin.Context) {
	path := c.Request.URL.Path
	
	// Don't serve SPA for API routes, health check, or swagger
	if strings.HasPrefix(path, "/api/") || 
	   strings.HasPrefix(path, "/health") || 
	   strings.HasPrefix(path, "/swagger/") {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	
	// Serve index.html for all other routes (SPA routing)
	indexPath := filepath.Join(s.frontendPath, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "frontend not found"})
		return
	}
	
	// Set proper HTML MIME type and cache headers
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate") // Don't cache HTML
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	
	c.File(indexPath)
}

// Start starts the API server
func (s *Server) Start(port int) error {
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}
	
	slog.Info("Starting API server", "addr", addr)
	return s.server.ListenAndServe()
}

// Stop gracefully stops the API server
func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	
	slog.Info("Stopping API server")
	return s.server.Shutdown(ctx)
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

// healthCheck returns server health status
// @Summary Health check
// @Description Check if the API server is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (s *Server) healthCheck(c *gin.Context) {
	version := s.version
	if version == "" {
		version = "dev"
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   version,
	})
}

// versionInfo returns detailed version information
// @Summary Version information
// @Description Get detailed version information including git commit
// @Tags version
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/version [get]
func (s *Server) versionInfo(c *gin.Context) {
	version := s.version
	if version == "" {
		version = "dev"
	}
	
	gitCommit := s.gitCommit
	if gitCommit == "" {
		gitCommit = "dirty"
	}
	
	c.JSON(http.StatusOK, gin.H{
		"version":    version,
		"git_commit": gitCommit,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	})
}

