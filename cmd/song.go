/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// songCmd represents the song command
var songCmd = &cobra.Command{
	Use:     "song <command>",
	Short:   "Interact with songs",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"s", "songs"},
}

func init() {
	rootCmd.AddCommand(songCmd)
}
