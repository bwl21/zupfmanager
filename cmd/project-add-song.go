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

		var songID int
		var songExists bool

		// Check if the --id flag is set
		idFlag := cmd.Flags().Lookup("id").Value.String() == "true"

		if idFlag {
			// If the --id flag is set, treat arg as song ID
			songID, err = strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid song ID: %v", err)
			}

			songExists, err = client.Song.Query().
				Where(song.ID(songID)).
				Exist(context.Background())
			if err != nil {
				return err
			}
			if !songExists {
				return fmt.Errorf("song with ID %d not found", songID)
			}
		} else {
			// If the --id flag is not set, treat arg as filename
			songs, err := client.Song.Query().
				Where(song.FilenameContains(args[1])).
				All(context.Background())
			if err != nil {
				return fmt.Errorf("error searching for songs with filename containing %s: %v", args[1], err)
			}

			if len(songs) == 0 {
				return fmt.Errorf("no songs found with filename containing %s", args[1])
			}

			if len(songs) > 1 {
				fmt.Println("Multiple songs found:")
				for _, s := range songs {
					fmt.Printf("ID: %d, Filename: %s\n", s.ID, s.Filename)
				}
				return fmt.Errorf("multiple songs found with filename containing %s, please specify the ID", args[1])
			}

			songID = songs[0].ID
			songExists = true
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

		// Get song filename
		songResult, err := client.Song.Get(context.Background(), songID)
		if err != nil {
			return fmt.Errorf("failed to get song: %w", err)
		}

		slog.Info("Added song to project",
			"project_id", projectID,
			"song_id", songID,
			"song_filename", songResult.Filename,
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

	// Flag to specify searching by ID
	projectAddSongCmd.Flags().Bool("id", false, "Search for song by ID")
}
