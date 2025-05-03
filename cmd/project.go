/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:     "project",
	Short:   "Interact with projects",
	Aliases: []string{"p", "projects"},
	Args:    cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
