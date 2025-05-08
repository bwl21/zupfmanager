package cmd

import (
	"fmt"
	"log"
	"os"

	"context"
	"github.com/bwl21/zupfmanager/internal/ent"
	"github.com/bwl21/zupfmanager/internal/ent/project"
	"github.com/spf13/cobra"
	"strconv"
	"text/tabwriter"
)

var projectDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a project",
	Long:  `Delete a project from the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide the project ID to delete.")
			os.Exit(1)
		}
		projectID := args[0]

		// Convert projectID to int
		id, err := strconv.Atoi(projectID)
		if err != nil {
			fmt.Println("Invalid project ID. Must be an integer.")
			os.Exit(1)
		}

		// Database connection
		client, err := ent.Open("sqlite3", "file:zupfmanager.db?_fk=1")
		if err != nil {
			log.Fatalf("failed to open sqlite3 connection: %v", err)
		}
		defer client.Close()

		// Query project with associated songs
		proj, err := client.Project.Query().
			Where(project.ID(id)).
			WithProjectSongs(func(q *ent.ProjectSongQuery) {
				q.WithSong()
			}).
			First(context.Background())

		if err != nil {
			fmt.Printf("Failed to find project with ID %s: %v\n", projectID, err)
			os.Exit(1)
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
			fmt.Fprintln(w, "ID\tTITLE\tPRIORITY\tDIFFICULTY\tCOMMENT")
			fmt.Fprintln(w, "--\t-----\t--------\t----------\t-------")

			for _, ps := range proj.Edges.ProjectSongs {
				comment := ps.Comment
				if comment == "" {
					comment = "-"
				}
				fmt.Fprintf(w, "%d\t%s\t%d\t%s\t%s\n",
					ps.Edges.Song.ID,
					ps.Edges.Song.Title,
					ps.Priority,
					ps.Difficulty,
					comment)
			}
			w.Flush()
		} else {
			fmt.Println("\nNo songs associated with this project.")
		}

		fmt.Printf("Are you sure you want to delete project with ID %s? (y/N): ", projectID)
		var confirmation string
		fmt.Scanln(&confirmation)

		if confirmation != "y" {
			fmt.Println("Deletion cancelled.")
			return
		}

		// Delete the project
		err = client.Project.DeleteOneID(id).Exec(context.Background())
		if err != nil {
			fmt.Printf("Failed to delete project with ID %s: %v\n", projectID, err)
			os.Exit(1)
		}

		fmt.Printf("Project with ID %s deleted successfully.\n", projectID)
	},
}

func init() {
	projectCmd.AddCommand(projectDeleteCmd)
}
