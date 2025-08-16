package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
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
	importHandler  *handlers.ImportHandler
	projectHandler *handlers.ProjectHandler
	songHandler    *handlers.SongHandler
}

// NewServer creates a new API server
func NewServer(services *core.Services) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	
	s := &Server{
		router:         router,
		services:       services,
		importHandler:  handlers.NewImportHandler(services),
		projectHandler: handlers.NewProjectHandler(services),
		songHandler:    handlers.NewSongHandler(services),
	}
	
	s.setupRoutes()
	return s
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.healthCheck)
	
	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Import endpoints
		v1.POST("/import/file", s.importHandler.ImportFile)
		v1.POST("/import/directory", s.importHandler.ImportDirectory)
		
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
			projects.GET("/:id", s.projectHandler.GetProject)
			projects.PUT("/:id", s.projectHandler.UpdateProject)
			projects.DELETE("/:id", s.projectHandler.DeleteProject)
		}
		
		// Song endpoints
		songs := v1.Group("/songs")
		{
			songs.GET("", s.songHandler.ListSongs)
			songs.GET("/:id", s.songHandler.GetSong)
			songs.GET("/search", s.songHandler.SearchSongs)
		}
	}
	
	// Dynamic swagger config that adapts to current host/scheme
	s.router.GET("/api/swagger.json", s.swaggerConfig)
	
	// Swagger documentation with dynamic URL
	url := ginSwagger.URL("/api/swagger.json")
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
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
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}

