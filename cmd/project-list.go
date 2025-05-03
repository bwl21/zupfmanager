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

// listProjectsCmd represents the list-projects command
var listProjectsCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all projects",
	Long:    `List all projects in the database.`,
	Aliases: []string{"l", "ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		client, err := database.New()
		if err != nil {
			return err
		}

		// Query all projects
		projects, err := client.Project.Query().All(context.Background())
		if err != nil {
			return err
		}

		// Check if json output is requested
		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			jsonData, err := json.MarshalIndent(projects, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(jsonData))
			return nil
		}

		// Setup tabwriter for aligned output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, "ID\tTITLE")
		fmt.Fprintln(w, "--\t-----")

		for _, p := range projects {
			fmt.Fprintf(w, "%d\t%s\n", p.ID, p.Title)
		}

		return w.Flush()
	},
}

func init() {
	projectCmd.AddCommand(listProjectsCmd)

	// Flags
	listProjectsCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}
