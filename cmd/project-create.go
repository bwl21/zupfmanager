/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"log/slog"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/spf13/cobra"
)

// projectCreateCmd represents the create-project command
var projectCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Long:  `Create a new project in the database with specified title and configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		client, err := database.New()
		if err != nil {
			return err
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
			if defaultCfg, _ := cmd.Flags().GetBool("default_config"); defaultCfg {
				// Load default config from file
				configFile, err := os.ReadFile("default-project-config.json")
				if err != nil {
					return fmt.Errorf("failed to read default config file: %w", err)
				}
				if err := json.Unmarshal(configFile, &config); err != nil {
					return fmt.Errorf("failed to parse default config JSON: %w", err)
				}
			}
		}

		_, err = client.CreateOrUpdateProject(context.Background(), 0, title, shortName, config)
		if err != nil {
			return err
		}

		// Create directory with shortName and tpl subdirectory
		projectDir := shortName
		tplDir := projectDir + "/tpl"
		err = os.MkdirAll(tplDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		slog.Info("Created project directory", "path", projectDir)

		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectCreateCmd)

	// Required flags
	projectCreateCmd.Flags().StringP("title", "t", "", "Title of the project (required)")
	projectCreateCmd.MarkFlagRequired("title")
	projectCreateCmd.Flags().StringP("short_name", "s", "", "Short name of the project (required)")
	projectCreateCmd.MarkFlagRequired("short_name")

	// Optional flags
	projectCreateCmd.Flags().StringP("config", "c", "", "Project configuration as JSON string")
	projectCreateCmd.Flags().BoolP("default_config", "d", false, "Load default project configuration from file")
}
