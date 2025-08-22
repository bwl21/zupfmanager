package core

import (
	"context"
	"fmt"
	"io/fs"
	"sync"

	"github.com/bwl21/zupfmanager/internal/database"
)

// Services container holds all service instances with shared dependencies
type Services struct {
	db         *database.Client
	Project    ProjectService
	Song       SongService
	Import     ImportService
	Config     ConfigService
	FileSystem FileSystemService
	
	// Resource management
	ctx    context.Context
	cancel context.CancelFunc
	closed bool
	mu     sync.RWMutex
}

// NewServices creates a new services container with shared database connection
func NewServices() (*Services, error) {
	return NewServicesWithContext(context.Background())
}

// NewServicesWithContext creates a new services container with a custom context
func NewServicesWithContext(ctx context.Context) (*Services, error) {
	return NewServicesWithEmbedded(ctx, nil)
}

// NewServicesWithEmbedded creates a new services container with embedded filesystem support
func NewServicesWithEmbedded(ctx context.Context, embeddedConfigFS fs.FS) (*Services, error) {
	db, err := database.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	var config ConfigService
	if embeddedConfigFS != nil {
		config = NewConfigServiceWithEmbedded("default-project-config.json", embeddedConfigFS)
	} else {
		config = NewConfigService()
	}
	
	fileSystem := NewFileSystemService()
	
	// Create cancellable context for resource management
	serviceCtx, cancel := context.WithCancel(ctx)

	return &Services{
		db:         db,
		Project:    NewProjectServiceWithDeps(db, config, fileSystem),
		Song:       NewSongServiceWithDeps(db),
		Import:     NewImportServiceWithDeps(db),
		Config:     config,
		FileSystem: fileSystem,
		ctx:        serviceCtx,
		cancel:     cancel,
		closed:     false,
	}, nil
}

// Context returns the services context
func (s *Services) Context() context.Context {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ctx
}

// Close closes all resources safely
func (s *Services) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.closed {
		return nil
	}
	
	s.closed = true
	
	// Cancel context to signal shutdown
	if s.cancel != nil {
		s.cancel()
	}
	
	// Close database connection
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
	}
	
	return nil
}

// IsClosed returns true if services have been closed
func (s *Services) IsClosed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.closed
}

// DB returns the database client for testing purposes
func (s *Services) DB() *database.Client {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db
}

// ensureNotClosed checks if services are still open
func (s *Services) ensureNotClosed() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return fmt.Errorf("services have been closed")
	}
	return nil
}
