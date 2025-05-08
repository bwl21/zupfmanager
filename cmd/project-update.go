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
	Use:   "update [projectID]",
	Short: "Update an existing project",
	Long:  `Update an existing project in the database with specified attributes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		client, err := database.New()
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return fmt.Errorf("project ID is required")
		}

		projectIDStr := args[0]
		projectID, err := strconv.Atoi(projectIDStr)
		if err != nil {
			return fmt.Errorf("invalid project ID: %w", err)
		}

		existingProject, err := client.GetProject(context.Background(), projectID)
		if err != nil {
			return fmt.Errorf("failed to get project: %w", err)
		}

		titleFlag, _ := cmd.Flags().GetString("title")
		shortNameFlag, _ := cmd.Flags().GetString("short_name")
		configStr, _ := cmd.Flags().GetString("config")

		// Use existing values if new values are not provided
		title := existingProject.Title
		if titleFlag != "" {
			title = titleFlag
		}

		shortName := existingProject.ShortName
		if shortNameFlag != "" {
			shortName = shortNameFlag
		}

		config := existingProject.Config
		if config == nil {
			config = map[string]interface{}{}
		}

		// Parse config JSON if provided
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
			defaultCfg, _ := cmd.Flags().GetString("default_config")
			if defaultCfg != "" {
				// Load default config from file
				configFile, err := os.ReadFile(defaultCfg)
				if err != nil {
					return fmt.Errorf("failed to read default config file: %w", err)
				}
				if err := json.Unmarshal(configFile, &config); err != nil {
					return fmt.Errorf("failed to parse default config JSON: %w", err)
				}
			}
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

	// Optional flags
	projectUpdateCmd.Flags().StringP("title", "t", "", "New title of the project")
	projectUpdateCmd.Flags().StringP("short_name", "s", "", "New short name of the project")
	projectUpdateCmd.Flags().StringP("config", "c", "", "Project configuration as JSON string")
    projectUpdateCmd.Flags().StringP("default_config", "d", "", "Load default project configuration from file")
}
