package cmd

import (
	"fmt"
	"os"
	"log"

	"github.com/spf13/cobra"
	"github.com/bwl21/zupfmanager/internal/ent"
	"context"
	"strconv"
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

		fmt.Printf("Are you sure you want to delete project with ID %s? (y/N): ", projectID)
		var confirmation string
		fmt.Scanln(&confirmation)

		if confirmation != "y" {
			fmt.Println("Deletion cancelled.")
			return
		}

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
