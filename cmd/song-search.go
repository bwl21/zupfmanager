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
	"github.com/bwl21/zupfmanager/internal/ent/predicate"
	"github.com/bwl21/zupfmanager/internal/ent/song"
	"github.com/spf13/cobra"
)

// songSearchCmd represents the search-songs command
var songSearchCmd = &cobra.Command{
	Use:     "search <query>",
	Short:   "Search for songs by title or filename",
	Long:    `Search for songs in the database by title or filename using substring matching.`,
	Aliases: []string{"s", "find"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		client, err := database.New()
		if err != nil {
			return err
		}

		// Get the search query
		searchQuery := args[0]

		// Option to specify search fields
		searchTitle, _ := cmd.Flags().GetBool("title")
		searchFilename, _ := cmd.Flags().GetBool("filename")
		searchGenre, _ := cmd.Flags().GetBool("genre")

		// If no specific fields are selected, search all fields
		if !searchTitle && !searchFilename && !searchGenre {
			searchTitle = true
			searchFilename = true
			searchGenre = true
		}

		// Build predicates for search
		var predicates []predicate.Song
		if searchTitle {
			predicates = append(predicates,
				song.TitleContainsFold(searchQuery))
		}
		if searchFilename {
			predicates = append(predicates,
				song.FilenameContainsFold(searchQuery))
		}
		if searchGenre && searchQuery != "" {
			predicates = append(predicates,
				song.GenreContainsFold(searchQuery))
		}

		// Query songs with the search term
		query := client.Song.Query()
		if len(predicates) > 0 {
			query = query.Where(song.Or(predicates...))
		}

		// Query matching songs
		songs, err := query.All(context.Background())
		if err != nil {
			return err
		}

		// Check if no results found
		if len(songs) == 0 {
			fmt.Println("No songs found matching the search query.")
			return nil
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
		fmt.Fprintln(w, "ID\tTITLE\tFILENAME\tGENRE")
		fmt.Fprintln(w, "--\t-----\t--------\t-----")

		for _, s := range songs {
			genre := s.Genre
			if genre == "" {
				genre = "-"
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", s.ID, s.Title, s.Filename, genre)
		}

		return w.Flush()
	},
}

func init() {
	songCmd.AddCommand(songSearchCmd)

	// Flags
	songSearchCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	songSearchCmd.Flags().BoolP("title", "t", false, "Search only in song titles")
	songSearchCmd.Flags().BoolP("filename", "f", false, "Search only in filenames")
	songSearchCmd.Flags().BoolP("genre", "g", false, "Search only in genres")
}
