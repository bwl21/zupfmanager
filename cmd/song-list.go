/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/spf13/cobra"
)

// songListCmd represents the list-songs command
var songListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all songs",
	Long:    `List all songs in the database.`,
	Aliases: []string{"l", "ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		client, err := database.New()
		if err != nil {
			return err
		}

		// Query all songs
		songs, err := client.Song.Query().All(context.Background())
		if err != nil {
			return err
		}

		// Check if json output is requested
		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			jsonData, err := json.MarshalIndent(songs, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(jsonData))
			return nil
		}

		// Setup tabwriter for aligned output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, "ID\tTITLE\tFILENAME\tCOPYRIGHT\tGENRE\tTOCINFO")
		fmt.Fprintln(w, "--\t-----\t--------\t---------\t-----\t------")

		for _, s := range songs {
			genre := s.Genre
			if genre == "" {
				genre = "-"
			}
			copyright := s.Copyright
			if copyright == "" {
				copyright = "-"
			}
			tocinfo := s.Tocinfo
			if tocinfo == "" {
				tocinfo = "-"
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n", s.ID, s.Title, s.Filename, copyright, genre, tocinfo)
		}

		return w.Flush()
	},
}

func init() {
	songCmd.AddCommand(songListCmd)

	// Flags
	songListCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
