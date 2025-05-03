/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/bwl21/zupfmanager/internal/ent"
	"github.com/bwl21/zupfmanager/internal/ent/song"
	"github.com/spf13/cobra"
)

// songShowCmd represents the show-song command
var songShowCmd = &cobra.Command{
	Use:     "show <song-id>",
	Short:   "Show song details",
	Long:    `Show detailed information about a song including its associated projects.`,
	Aliases: []string{"s", "show"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		client, err := database.New()
		if err != nil {
			return err
		}

		// Parse song ID
		songID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid song ID: %v", err)
		}

		// Query song with associated projects
		s, err := client.Song.Query().
			Where(song.ID(songID)).
			WithProjectSongs(func(q *ent.ProjectSongQuery) {
				q.WithProject()
			}).
			First(context.Background())

		if err != nil {
			return err
		}

		// Check if json output is requested
		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			jsonData, err := json.MarshalIndent(s, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(jsonData))
			return nil
		}

		// Display song details
		fmt.Printf("Song: %s (ID: %d)\n", s.Title, s.ID)
		fmt.Printf("Filename: %s\n", s.Filename)

		if s.Genre != "" {
			fmt.Printf("Genre: %s\n", s.Genre)
		}

		// Display associated projects if any
		if len(s.Edges.ProjectSongs) > 0 {
			fmt.Println("\nUsed in Projects:")
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
			fmt.Fprintln(w, "ID\tPROJECT\tPRIORITY\tDIFFICULTY\tCOMMENT")
			fmt.Fprintln(w, "--\t-------\t--------\t----------\t-------")

			for _, ps := range s.Edges.ProjectSongs {
				comment := ps.Comment
				if comment == "" {
					comment = "-"
				}
				fmt.Fprintf(w, "%d\t%s\t%d\t%s\t%s\n",
					ps.Edges.Project.ID,
					ps.Edges.Project.Title,
					ps.Priority,
					ps.Difficulty,
					comment)
			}
			w.Flush()
		} else {
			fmt.Println("\nThis song is not used in any projects.")
		}

		return nil
	},
}

func init() {
	songCmd.AddCommand(songShowCmd)

	// Flags
	songShowCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
