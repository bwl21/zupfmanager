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

		// Customize the slog output format to remove the timestamp
		replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		}

		handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{ReplaceAttr: replaceAttr})
		slog.SetDefault(slog.New(handler))

		client, err := database.New()
		if err != nil {
			slog.Error("Failed to create database client", "error", err)
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
		title     string
		genre     string
		copyright string
		tocinfo   string
	)
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "T:") {
			title = strings.TrimPrefix(line, "T:")
			title = strings.TrimSpace(title)
		} else if strings.HasPrefix(line, "Z:genre") {
			genre = strings.TrimPrefix(line, "Z:genre")
			genre = strings.TrimSpace(genre)
		} else if strings.HasPrefix(line, "Z:copyright") {
			copyright = strings.TrimPrefix(line, "Z:copyright")
			copyright = strings.TrimSpace(copyright)
		} else if tocinfo == "" {
			// Überprüfe auf C:M: Muster
			if strings.HasPrefix(line, "C:M: ") {
				tocinfo = strings.TrimSpace(strings.TrimPrefix(line, "C:M: "))
			} else if strings.HasPrefix(line, "C:M:") {
				tocinfo = strings.TrimSpace(strings.TrimPrefix(line, "C:M:"))
			} else if strings.HasPrefix(line, "C:M+T: ") {
				// Überprüfe auf C:M+T: Muster
				tocinfo = strings.TrimSpace(strings.TrimPrefix(line, "C:M+T: "))
			} else if strings.HasPrefix(line, "C:M+T:") {
				tocinfo = strings.TrimSpace(strings.TrimPrefix(line, "C:M+T:"))
			} else if strings.HasPrefix(line, "C:T+M: ") {
				// Überprüfe auf C:T+M: Muster
				tocinfo = strings.TrimSpace(strings.TrimPrefix(line, "C:T+M: "))
			} else if strings.HasPrefix(line, "C:T+M:") {
				tocinfo = strings.TrimSpace(strings.TrimPrefix(line, "C:T+M:"))
			}
		}
	}
	if title == "" {
		return fmt.Errorf("no title found in file")
	}

	filename := filepath.Base(file)

	sng, err := client.Song.Query().Where(song.Filename(filename)).First(context.Background())
	if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("failed to query song: %w", err)
	}

	if sng == nil {
		// Create a new song
		slog.Info("Creating new song", "filename", filename, "title", title, "genre", genre, "copyright", copyright)
		_, err = client.Song.Create().
			SetTitle(title).
			SetFilename(filename).
			SetGenre(genre).
			SetCopyright(copyright).
			SetTocinfo(tocinfo).
			Save(context.Background())
		if err != nil {
			slog.Error("Failed to create song", "filename", filename, "error", err)
			return fmt.Errorf("failed to create song: %w", err)
		}
		slog.Info("Successfully created song", "filename", filename, "title", title, "genre", genre, "copyright", copyright)
	} else {
		// Update an existing song
		changes := make([]string, 0)
		if sng.Title != title {
			changes = append(changes, fmt.Sprintf("title: %s -> %s", sng.Title, title))
			sng.Title = title
		}
		if sng.Genre != genre {
			changes = append(changes, fmt.Sprintf("genre: %s -> %s", sng.Genre, genre))
			sng.Genre = genre
		}
		if sng.Copyright != copyright {
			changes = append(changes, fmt.Sprintf("copyright: %s -> %s", sng.Copyright, copyright))
			sng.Copyright = copyright
		}
		if sng.Tocinfo != tocinfo {
			changes = append(changes, fmt.Sprintf("tocinfo: %s -> %s", sng.Tocinfo, tocinfo))
			sng.Tocinfo = tocinfo
		}

		if len(changes) > 0 {
			slog.Info("Updating existing song", "filename", filename, "changes", strings.Join(changes, ", "))
			_, err = sng.Update().
				SetTitle(title).
				SetGenre(genre).
				SetCopyright(copyright).
				SetTocinfo(tocinfo).
				Save(context.Background())
			if err != nil {
				slog.Error("Failed to update song", "filename", filename, "error", err)
				return fmt.Errorf("failed to update song: %w", err)
			}
			slog.Info("Successfully updated song", "filename", filename, "changes", strings.Join(changes, ", "))
		}
	}

	return nil
}
