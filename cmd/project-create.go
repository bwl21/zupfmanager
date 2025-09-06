/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log/slog"

	"github.com/bwl21/zupfmanager/pkg/core"
	"github.com/spf13/cobra"
)

// projectCreateCmd represents the create-project command
var projectCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Long:  `Create a new project in the database with specified title and configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		services, err := core.NewServices()
		if err != nil {
			return err
		}
		defer services.Close()

		title, _ := cmd.Flags().GetString("title")
		shortName, _ := cmd.Flags().GetString("short_name")
		configFile, _ := cmd.Flags().GetString("config")
		defaultConfig, _ := cmd.Flags().GetBool("default_config")

		req := core.CreateProjectRequest{
			Title:         title,
			ShortName:     shortName,
			ConfigFile:    configFile,
			DefaultConfig: defaultConfig,
		}

		_, err = services.Project.Create(context.Background(), req)
		if err != nil {
			return err
		}

		slog.Info("Created project directory", "path", shortName)

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
