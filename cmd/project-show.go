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
	"github.com/bwl21/zupfmanager/internal/ent/project"
	"github.com/bwl21/zupfmanager/internal/ent/projectsong"
	"github.com/spf13/cobra"
)

// projectShowCmd represents the show-project command
var projectShowCmd = &cobra.Command{
	Use:     "show <project-id>",
	Short:   "Show project details",
	Long:    `Show detailed information about a project including its associated songs.`,
	Aliases: []string{"s", "show"},
	Args:    cobra.ExactArgs(1),
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

		// First, get the project to ensure it exists
		proj, err := client.Project.Get(context.Background(), projectID)
		if ent.IsNotFound(err) {
			return fmt.Errorf("project with ID %d not found", projectID)
		}
		if err != nil {
			return fmt.Errorf("error fetching project: %w", err)
		}

		// Then query the project songs separately
		projectSongs, err := client.ProjectSong.Query().
			Where(projectsong.HasProjectWith(project.ID(projectID))).
			WithSong().
			Order(ent.Asc("priority")).
			All(context.Background())
		if err != nil {
			return fmt.Errorf("error fetching project songs: %w", err)
		}

		// Attach the songs to the project for compatibility with the rest of the code
		proj.Edges.ProjectSongs = projectSongs

		// Check if json output is requested
		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			jsonData, err := json.MarshalIndent(proj, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(jsonData))
			return nil
		}

		// Display project details
		fmt.Printf("Project: %s (ID: %d)\n", proj.Title, proj.ID)
		fmt.Printf("Number of songs: %d\n", len(proj.Edges.ProjectSongs))

		if len(proj.Config) > 0 {
			fmt.Println("\nConfiguration:")
			for k, v := range proj.Config {
				fmt.Printf("  %s: %v\n", k, v)
			}
		}

		// Display associated songs if any
		if len(proj.Edges.ProjectSongs) > 0 {
			fmt.Println("\nSongs:")
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
			fmt.Fprintln(w, "ID\tFILENAME\tPRIORITY\tDIFFICULTY\tCOPYRIGHT\tGENRE\tTOCINFO")
			fmt.Fprintln(w, "--\t--------\t--------\t----------\t---------\t-----\t-------")
			for _, ps := range proj.Edges.ProjectSongs {
				comment := ps.Comment
				if comment == "" {
					comment = "-"
				}
				fmt.Fprintf(w, "%d\t%s\t%d\t%s\t%s\t%s\t%s\n",
					ps.Edges.Song.ID,
					ps.Edges.Song.Filename,
					ps.Priority,
					ps.Difficulty,
					ps.Edges.Song.Copyright,
					ps.Edges.Song.Genre,
					ps.Edges.Song.Tocinfo,
				)

			}
			w.Flush()
		} else {
			fmt.Println("\nNo songs associated with this project.")
		}

		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectShowCmd)

	// Flags
	projectShowCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
