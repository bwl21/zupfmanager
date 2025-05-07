package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/spf13/cobra"
)

// projectUpdateCmd represents the project edit config command
var projectUpdateCmd = &cobra.Command{
	Use:   "update PROJECT_ID <filename|->",
	Short: "Update the config JSON for a project from a file or STDIN",
	Long: `Updates the 'config' field of a project specified by its ID.
The configuration should be provided as a JSON object.
Input can be read from a file using the --input flag or from STDIN if the flag is omitted.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid project ID: %w", err)
		}

		projectEditConfigInputFile := args[1]

		var reader io.Reader
		if projectEditConfigInputFile == "-" {
			fmt.Println("Reading project config JSON from STDIN. Press Ctrl+D when finished.")
			reader = os.Stdin
		} else {
			file, err := os.Open(projectEditConfigInputFile)
			if err != nil {
				return fmt.Errorf("failed to open input file %s: %w", projectEditConfigInputFile, err)
			}
			defer file.Close()
			reader = file
		}

		jsonData, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		var configData map[string]interface{}
		if err := json.Unmarshal(jsonData, &configData); err != nil {
			return fmt.Errorf("failed to parse JSON input: %w", err)
		}

		client, err := database.New()
		if err != nil {
			return err
		}
		defer client.Close()

		project, err := client.Project.Get(context.Background(), projectID)
		if err != nil {
			return fmt.Errorf("failed to find project with ID %d: %w", projectID, err)
		}

		_, err = project.Update().
			SetConfig(configData).
			Save(context.Background())
		if err != nil {
			return fmt.Errorf("failed to update project config for ID %d: %w", projectID, err)
		}

		fmt.Printf("Successfully updated config for project ID %d\n", projectID)
		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectUpdateCmd)
}
