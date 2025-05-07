package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/spf13/cobra"
)

// projectUpdateCmd represents the project edit config command
var projectUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing project",
	Long:  `Update an existing project in the database with specified attributes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		client, err := database.New()
		if err != nil {
			return err
		}

		projectIDStr, _ := cmd.Flags().GetString("id")
		projectID, err := strconv.Atoi(projectIDStr)
		if err != nil {
			return fmt.Errorf("invalid project ID: %w", err)
		}

		title, _ := cmd.Flags().GetString("title")
		shortName, _ := cmd.Flags().GetString("short_name")
		configStr, _ := cmd.Flags().GetString("config")

		// Parse config JSON if provided
		var config map[string]interface{}
		if configStr != "" {
			// Read config from file
			configFile, err := os.ReadFile(configStr)
			if err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}
			if err := json.Unmarshal(configFile, &config); err != nil {
				return fmt.Errorf("failed to parse config JSON: %w", err)
			}
		} else {
			config = map[string]interface{}{}
		}

		_, err = client.CreateOrUpdateProject(context.Background(), projectID, title, shortName, config)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectUpdateCmd)

	// Required flags
	projectUpdateCmd.Flags().StringP("id", "i", "", "ID of the project to update (required)")
	projectUpdateCmd.MarkFlagRequired("id")

	// Optional flags
	projectUpdateCmd.Flags().StringP("title", "t", "", "New title of the project")
	projectUpdateCmd.Flags().StringP("short_name", "s", "", "New short name of the project")
	projectUpdateCmd.Flags().StringP("config", "c", "", "Project configuration as JSON string")
}
