/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/bwl21/zupfmanager/internal/ui"
	"github.com/spf13/cobra"
)

// uiCmd represents the ui command
var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Launch the interactive terminal UI",
	Long: `Launch an interactive terminal UI for managing projects and songs.
Navigate using keyboard shortcuts.

The UI provides the same functionality as the CLI commands in an interactive interface.`,
	Aliases: []string{"tui"},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return ui.RunUI()
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
}
