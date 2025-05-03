/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
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
		configStr, _ := cmd.Flags().GetString("config")

		// Parse config JSON if provided
		var config map[string]interface{}
		if configStr != "" {
			if err := json.Unmarshal([]byte(configStr), &config); err != nil {
				return err
			}
		} else {
			config = map[string]interface{}{}
		}

		project, err := client.Project.Create().
			SetTitle(title).
			SetConfig(config).
			Save(context.Background())

		if err != nil {
			return err
		}

		slog.Info("Created new project", "id", project.ID, "title", project.Title)
		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectCreateCmd)

	// Required flags
	projectCreateCmd.Flags().StringP("title", "t", "", "Title of the project (required)")
	projectCreateCmd.MarkFlagRequired("title")

	// Optional flags
	projectCreateCmd.Flags().StringP("config", "c", "", "Project configuration as JSON string")
}
