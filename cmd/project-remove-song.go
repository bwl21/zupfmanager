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

// projectRemoveSongCmd represents the remove-song-from-project command
var projectRemoveSongCmd = &cobra.Command{
	Use:     "remove <project-id> <song-id>",
	Short:   "Remove a song from a project",
	Long:    `Remove a song from a project while keeping the song in the database.`,
	Aliases: []string{"remove-song", "rm", "rs"},
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

		// Verify song is associated with project
		projectSong, err := client.ProjectSong.Query().
			Where(
				projectsong.ProjectID(projectID),
				projectsong.SongID(songID),
			).
			Only(context.Background())
		if err != nil {
			return fmt.Errorf("song ID %d is not associated with project ID %d", songID, projectID)
		}

		// Delete the project-song relationship
		err = client.ProjectSong.DeleteOne(projectSong).Exec(context.Background())
		if err != nil {
			return err
		}

		slog.Info("Removed song from project",
			"project_id", projectID,
			"song_id", songID)

		fmt.Printf("Song ID %d removed from project ID %d\n", songID, projectID)
		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectRemoveSongCmd)
}
