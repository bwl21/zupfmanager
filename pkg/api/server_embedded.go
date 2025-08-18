package api

import (
	"log/slog"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// setupEmbeddedFrontend configures embedded frontend serving
func (s *Server) setupEmbeddedFrontend() error {
	frontendFS, err := GetFrontendFS()
	if err != nil {
		return err
	}
	
	s.frontendFS = frontendFS
	s.useEmbedded = true
	
	slog.Info("Using embedded frontend files")
	
	// Serve static assets
	s.router.GET("/assets/*filepath", s.serveEmbeddedStatic)
	
	// Serve favicon and other root files
	s.router.GET("/favicon.ico", func(c *gin.Context) {
		s.serveEmbeddedFile(c, "favicon.ico")
	})
	s.router.GET("/vite.svg", func(c *gin.Context) {
		s.serveEmbeddedFile(c, "vite.svg")
	})
	
	// SPA fallback for embedded files
	s.router.NoRoute(s.serveEmbeddedSPA)
	
	return nil
}

// serveEmbeddedStatic serves static files from embedded filesystem
func (s *Server) serveEmbeddedStatic(c *gin.Context) {
	filePath := c.Param("filepath")
	fullPath := path.Join("assets", filePath)
	
	s.serveEmbeddedFile(c, fullPath)
}

// serveEmbeddedFile serves a file from embedded filesystem
func (s *Server) serveEmbeddedFile(c *gin.Context, filePath string) {
	if s.frontendFS == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "frontend not available"})
		return
	}
	
	file, err := s.frontendFS.Open(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	defer file.Close()
	
	// Get file info for content type detection
	stat, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "file stat error"})
		return
	}
	
	// Set MIME type based on file extension
	ext := path.Ext(filePath)
	var mimeType string
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
	case ".html":
		mimeType = "text/html; charset=utf-8"
	}
	
	if mimeType != "" {
		c.Header("Content-Type", mimeType)
	}
	
	// Set cache headers
	if strings.HasPrefix(filePath, "assets/") {
		c.Header("Cache-Control", "public, max-age=31536000") // 1 year
	} else {
		c.Header("Cache-Control", "public, max-age=3600") // 1 hour
	}
	
	// Stream the file content
	c.DataFromReader(http.StatusOK, stat.Size(), mimeType, file, nil)
}

// serveEmbeddedSPA serves the SPA index.html from embedded filesystem
func (s *Server) serveEmbeddedSPA(c *gin.Context) {
	path := c.Request.URL.Path
	
	// Don't serve SPA for API routes, health check, or swagger
	if strings.HasPrefix(path, "/api/") || 
	   strings.HasPrefix(path, "/health") || 
	   strings.HasPrefix(path, "/swagger/") {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	
	if s.frontendFS == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "frontend not available"})
		return
	}
	
	// Serve index.html for SPA routing
	file, err := s.frontendFS.Open("index.html")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "frontend not found"})
		return
	}
	defer file.Close()
	
	// Set HTML headers
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	
	// Stream the HTML content
	stat, _ := file.Stat()
	c.DataFromReader(http.StatusOK, stat.Size(), "text/html; charset=utf-8", file, nil)
}
