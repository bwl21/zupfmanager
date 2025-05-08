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

		// Query project with associated songs
		proj, err := client.Project.Query().
			Where(project.ID(projectID)).
			WithProjectSongs(func(q *ent.ProjectSongQuery) {
				q.WithSong()
			}).
			First(context.Background())

		if err != nil {
			return err
		}

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
			fmt.Fprintln(w, "ID\tFILENAME\tPRIORITY\tDIFFICULTY\tCOPYRIGHT\tGENRE")
			fmt.Fprintln(w, "--\t--------\t--------\t----------\t---------\t-----")
			for _, ps := range proj.Edges.ProjectSongs {
				comment := ps.Comment
				if comment == "" {
					comment = "-"
				}
				fmt.Fprintf(w, "%d\t%s\t%d\t%s\t%s\t%s\n",
					ps.Edges.Song.ID,
					ps.Edges.Song.Filename,
					ps.Priority,
					ps.Difficulty,
					ps.Edges.Song.Copyright,
					ps.Edges.Song.Genre,
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
