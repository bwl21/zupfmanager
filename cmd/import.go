/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/bwl21/zupfmanager/internal/ent"
	"github.com/bwl21/zupfmanager/internal/ent/song"
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

		client, err := database.New()
		if err != nil {
			return err
		}

		files, err := filepath.Glob(filepath.Join(args[0], "*.abc"))
		if err != nil {
			return err
		}

		for _, file := range files {
			err = importFile(client, file)
			if err != nil {
				slog.Warn("Failed to import file", "file", file, "error", err)
				continue
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}

func importFile(client *database.Client, file string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	var (
		title string
	)
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "T:") {
			title = strings.TrimPrefix(line, "T:")
			title = strings.TrimSpace(title)
			break
		}
	}
	if title == "" {
		return fmt.Errorf("no title found in file")
	}

	filename := filepath.Base(file)
	_, err = client.Song.Create().SetFilename(filename).SetTitle(title).Save(context.Background())
	if ent.IsConstraintError(err) {
		sng, err := client.Song.Query().Where(song.Filename(filename)).First(context.Background())
		if err != nil {
			return err
		}
		if sng.Title == title {
			return nil
		}

		sng.Title = title
		_, err = sng.Update().Save(context.Background())
		if err != nil {
			return err
		}
		slog.Info("Updated from file", "file", file)

		return nil
	} else if err != nil {
		return err
	} else {
		slog.Info("Imported from file", "file", file)
	}

	return nil
}
