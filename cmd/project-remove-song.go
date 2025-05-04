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

		// Get song filename
		songResult, err := client.Song.Get(context.Background(), songID)
		if err != nil {
			return fmt.Errorf("failed to get song: %w", err)
		}

		slog.Info("Removed song from project",
			"project_id", projectID,
			"song_id", songID,
			"song_filename", songResult.Filename)

		fmt.Printf("Song ID %d (%s) removed from project ID %d\n", songID, songResult.Filename, projectID)
		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectRemoveSongCmd)

	// Flag to specify searching by ID
	projectRemoveSongCmd.Flags().Bool("id", false, "Search for song by ID")
}
