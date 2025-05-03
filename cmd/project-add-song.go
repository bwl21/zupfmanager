/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/bwl21/zupfmanager/internal/ent/project"
	"github.com/bwl21/zupfmanager/internal/ent/projectsong"
	"github.com/bwl21/zupfmanager/internal/ent/song"
	"github.com/spf13/cobra"
)

// projectAddSongCmd represents the add-song-to-project command
var projectAddSongCmd = &cobra.Command{
	Use:     "add <project-id> <song-id>",
	Short:   "Add a song to a project",
	Long:    `Add an existing song to a project with specified difficulty and priority.`,
	Aliases: []string{"add-song", "as"},
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		client, err := database.New()
		if err != nil {
			return err
		}

		// Parse project ID
		projectID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid project ID: %v", err)
		}

		// Parse song ID
		songID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid song ID: %v", err)
		}

		// Verify project exists
		projectExists, err := client.Project.Query().
			Where(project.ID(projectID)).
			Exist(context.Background())
		if err != nil {
			return err
		}
		if !projectExists {
			return fmt.Errorf("project with ID %d not found", projectID)
		}

		// Verify song exists
		songExists, err := client.Song.Query().
			Where(song.ID(songID)).
			Exist(context.Background())
		if err != nil {
			return err
		}
		if !songExists {
			return fmt.Errorf("song with ID %d not found", songID)
		}

		// Get flags
		priority, _ := cmd.Flags().GetInt("priority")
		difficulty, _ := cmd.Flags().GetString("difficulty")
		comment, _ := cmd.Flags().GetString("comment")

		// Add song to project
		projectSong, err := client.ProjectSong.Create().
			SetProjectID(projectID).
			SetSongID(songID).
			SetPriority(priority).
			SetDifficulty(projectsong.Difficulty(difficulty)).
			SetComment(comment).
			Save(context.Background())
		if err != nil {
			return err
		}

		slog.Info("Added song to project",
			"project_id", projectID,
			"song_id", songID,
			"priority", projectSong.Priority,
			"difficulty", projectSong.Difficulty)

		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectAddSongCmd)

	// Required flags
	projectAddSongCmd.Flags().IntP("priority", "p", 1, "Priority of the song (1-4)")
	projectAddSongCmd.Flags().StringP("difficulty", "d", "medium", "Difficulty of the song: must be one of easy, medium, hard, expert")

	// Optional flags
	projectAddSongCmd.Flags().StringP("comment", "c", "", "Optional comment")
}
