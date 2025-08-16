package cmd

import (
	"context"
	"log/slog"
	"os"
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

Options:
- API only: zupfmanager api --port 8080
- Integrated (API + Frontend): zupfmanager api --port 8080 --frontend frontend/dist

Access the API documentation at http://localhost:8080/swagger/index.html`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		port, _ := cmd.Flags().GetInt("port")
		frontendPath, _ := cmd.Flags().GetString("frontend")

		// Create services
		services, err := core.NewServices()
		if err != nil {
			return err
		}
		defer services.Close()

		// Create API server with optional frontend serving
		var server *api.Server
		if frontendPath != "" {
			server = api.NewServer(services, api.ServerOptions{
				FrontendPath: frontendPath,
			})
		} else {
			server = api.NewServer(services)
		}

		// Setup graceful shutdown
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

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
