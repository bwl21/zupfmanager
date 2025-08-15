/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/bwl21/zupfmanager/pkg/core"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import <directory>",
	Short: "Import a directory of ABC files",
	Long:  `Import a directory of ABC files into the database.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		// Customize the slog output format to remove the timestamp
		replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		}

		handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{ReplaceAttr: replaceAttr})
		slog.SetDefault(slog.New(handler))

		service, err := core.NewImportService()
		if err != nil {
			slog.Error("Failed to create import service", "error", err)
			return err
		}
		defer service.Close()

		results, err := service.ImportDirectory(context.Background(), args[0])
		if err != nil {
			return err
		}

		for _, result := range results {
			if result.Error != nil {
				slog.Warn("Failed to import file", "file", result.Filename, "error", result.Error)
				continue
			}

			switch result.Action {
			case "created":
				slog.Info("Creating new song", "filename", result.Filename, "title", result.Title)
				slog.Info("Successfully created song", "filename", result.Filename, "title", result.Title)
			case "updated":
				slog.Info("Updating existing song", "filename", result.Filename, "changes", strings.Join(result.Changes, ", "))
				slog.Info("Successfully updated song", "filename", result.Filename, "changes", strings.Join(result.Changes, ", "))
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}

