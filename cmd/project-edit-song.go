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

// projectEditSongCmd represents the edit-song-in-project command
var projectEditSongCmd = &cobra.Command{
	Use:   "edit-song <project-id> <song-id>",
	Short: "Edit a song entry in a project",
	Long:  `Edit priority, difficulty, and comment for a song that is part of a project.`,
	Args:  cobra.ExactArgs(2),
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

		// Get flags
		priorityFlag := cmd.Flags().Changed("priority")
		difficultyFlag := cmd.Flags().Changed("difficulty")
		commentFlag := cmd.Flags().Changed("comment")

		// Return error if no flags are provided
		if !priorityFlag && !difficultyFlag && !commentFlag {
			return fmt.Errorf("at least one of --priority, --difficulty, or --comment must be provided")
		}

		// Initialize the update
		update := client.ProjectSong.UpdateOne(projectSong)

		// Update priority if changed
		if priorityFlag {
			priority, _ := cmd.Flags().GetInt("priority")
			update = update.SetPriority(priority)
		}

		// Update difficulty if changed
		if difficultyFlag {
			difficulty, _ := cmd.Flags().GetString("difficulty")
			update = update.SetDifficulty(projectsong.Difficulty(difficulty))
		}

		// Update comment if changed
		if commentFlag {
			comment, _ := cmd.Flags().GetString("comment")
			update = update.SetComment(comment)
		}

		// Save the updates
		updatedProjectSong, err := update.Save(context.Background())
		if err != nil {
			return err
		}

		slog.Info("Updated song in project",
			"project_id", projectID,
			"song_id", songID,
			"priority", updatedProjectSong.Priority,
			"difficulty", updatedProjectSong.Difficulty)

		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectEditSongCmd)

	// Flags
	projectEditSongCmd.Flags().IntP("priority", "p", 1, "Priority of the song (1-4)")
	projectEditSongCmd.Flags().StringP("difficulty", "d", "medium", "Difficulty of the song: must be one of easy, medium, hard, expert")
	projectEditSongCmd.Flags().StringP("comment", "c", "", "Comment for the song in this project")
}
