package cmd

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwl21/zupfmanager/pkg/api"
	"github.com/bwl21/zupfmanager/pkg/core"
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the REST API server",
	Long: `Start the REST API server to provide HTTP endpoints for Zupfmanager functionality.

The API server provides endpoints for:
- Importing ABC notation files
- Managing projects and songs
- Searching and listing content

The server includes an embedded frontend by default. Use --frontend to override with external files.

Options:
- Embedded frontend: zupfmanager api --port 8080
- External frontend: zupfmanager api --port 8080 --frontend frontend/dist

Access the web interface at http://localhost:8080/
Access the API documentation at http://localhost:8080/swagger/index.html`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		port, _ := cmd.Flags().GetInt("port")
		frontendPath, _ := cmd.Flags().GetString("frontend")

		// Auto-detect frontend path if not provided
		if frontendPath == "" {
			// Check for frontend files in common locations
			possiblePaths := []string{
				"dist/frontend",
				"frontend/dist",
			}
			for _, path := range possiblePaths {
				if _, err := os.Stat(filepath.Join(path, "index.html")); err == nil {
					frontendPath = path
					break
				}
			}
		}

		// Try to get embedded config filesystem
		var embeddedConfigFS fs.FS
		if configFS, err := api.GetDefaultConfigFS(); err == nil {
			embeddedConfigFS = configFS
		}

		// Create services with embedded config support
		services, err := core.NewServicesWithEmbedded(context.Background(), embeddedConfigFS)
		if err != nil {
			return err
		}
		defer services.Close()

		// Create API server with embedded frontend (fallback to external if provided)
		server := api.NewServer(services, api.ServerOptions{
			FrontendPath: frontendPath, // Used as fallback if embedded fails
			UseEmbedded:  true,         // Try embedded frontend first
			Version:      Version,
			GitCommit:    GitCommit,
		})

		// Setup graceful shutdown
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Get and log the current working directory
		if wd, err := os.Getwd(); err == nil {
			slog.Info("Current working directory", "path", wd)
		} else {
			slog.Warn("Could not get working directory", "error", err)
		}

		// Handle shutdown signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Start server in goroutine
		go func() {
			if err := server.Start(port); err != nil {
				slog.Error("Server failed to start", "error", err)
				cancel()
			}
		}()

		// Wait for shutdown signal
		select {
		case <-sigChan:
			slog.Info("Received shutdown signal")
		case <-ctx.Done():
			slog.Info("Context cancelled")
		}

		// Graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := server.Stop(shutdownCtx); err != nil {
			slog.Error("Server shutdown failed", "error", err)
			return err
		}

		slog.Info("Server shutdown complete")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.Flags().IntP("port", "p", 8080, "Port to run the API server on")
	apiCmd.Flags().StringP("frontend", "f", "", "Path to frontend dist directory (optional)")
}
