/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bwl21/zupfmanager/pkg/core"
	"github.com/spf13/cobra"
)

var forceDelete bool

// songDeleteCmd represents the delete song command
var songDeleteCmd = &cobra.Command{
	Use:     "delete <song-id>",
	Short:   "Delete a song",
	Long:    `Delete a song from the database. The song must not be used in any projects.`,
	Aliases: []string{"del", "rm", "remove"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		songID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid song ID: %s", args[0])
		}

		services, err := core.NewServices()
		if err != nil {
			return err
		}
		defer services.Close()

		ctx := context.Background()

		// Get song details first for confirmation
		song, err := services.Song.Get(ctx, songID)
		if err != nil {
			return fmt.Errorf("failed to get song: %w", err)
		}

		// Confirmation prompt (unless --force is used)
		if !forceDelete {
			fmt.Printf("Are you sure you want to delete song '%s' (ID: %d, File: %s)? [y/N]: ", 
				song.Title, song.ID, song.Filename)
			
			var response string
			fmt.Scanln(&response)
			
			if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
				fmt.Println("Delete cancelled.")
				return nil
			}
		}

		// Delete the song
		err = services.Song.Delete(ctx, songID)
		if err != nil {
			return fmt.Errorf("failed to delete song: %w", err)
		}

		fmt.Printf("Successfully deleted song '%s' (ID: %d)\n", song.Title, song.ID)
		return nil
	},
}

func init() {
	songCmd.AddCommand(songDeleteCmd)
	songDeleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Force delete without confirmation")
}
