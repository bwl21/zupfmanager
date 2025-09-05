package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/bwl21/zupfmanager/internal/ent"
	"github.com/bwl21/zupfmanager/internal/ent/projectsong"
	"github.com/bwl21/zupfmanager/internal/htmlpdf"
	"github.com/bwl21/zupfmanager/pkg/core"
	"github.com/spf13/cobra"
)

const (
	zupfnoterConfigString = "%%%%zupfnoter.config"
)

var (
	projectBuildOutputDir         string
	projectBuildAbcFileDir        string
	projectBuildPriorityThreshold int
	projectSampleId               string
)

var projectBuildCmd = &cobra.Command{
	Use:   "build PROJECT_ID",
	Short: "Build a project",
	Long:  `Builds a project by running the build command specified in the project's configuration.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid project ID: %w", err)
		}

		client, err := database.New()
		if err != nil {
			return err
		}
		defer client.Close()

		// First, get the project to ensure it exists
		project, err := client.Project.Get(context.Background(), projectID)
		if err != nil {
			return fmt.Errorf("failed to find project with ID %d: %w", projectID, err)
		}

		// Then query the project songs separately with the priority filter
		projectSongs, err := client.ProjectSong.Query().
			Where(
				projectsong.And(
					projectsong.ProjectID(projectID),
					projectsong.PriorityLTE(projectBuildPriorityThreshold),
				),
			).
			WithSong().
			WithProject().
			Order(ent.Asc("priority")).
			All(context.Background())
		if err != nil {
			return fmt.Errorf("failed to query project songs: %w", err)
		}

		// Attach the songs to the project for compatibility with the rest of the code
		project.Edges.ProjectSongs = projectSongs

		if projectBuildOutputDir == "" {
			projectBuildOutputDir = project.ShortName
		}

		if projectBuildAbcFileDir == "" {
			// Priority order: abc_file_dir_preference > abc_file_dir (from config) > last import directory
			if project.AbcFileDirPreference != "" {
				projectBuildAbcFileDir = project.AbcFileDirPreference
			} else if abcFileDir, ok := project.Config["abc_file_dir"].(string); ok && abcFileDir != "" {
				projectBuildAbcFileDir = abcFileDir
			} else {
				// Try to use the most recent import directory as default
				lastImportDir, err := core.GetLastImportDir()
				if err == nil && lastImportDir != "" {
					projectBuildAbcFileDir = lastImportDir
				} else {
					// Provide a default value or handle the error appropriately
					projectBuildAbcFileDir = ""
				}
			}
		}

		// Use core service for build logic
		services, err := core.NewServices()
		if err != nil {
			return err
		}
		defer services.Close()

		buildReq := core.BuildProjectRequest{
			ProjectID:         projectID,
			OutputDir:         projectBuildOutputDir,
			AbcFileDir:        projectBuildAbcFileDir,
			PriorityThreshold: projectBuildPriorityThreshold,
			SampleID:          projectSampleId,
		}

		return services.Project.ExecuteProjectBuild(context.Background(), buildReq)
	},
}

func init() {
	projectCmd.AddCommand(projectBuildCmd)

	projectBuildCmd.Flags().StringVarP(&projectBuildOutputDir, "output-dir", "o", "", "The directory to output the build results")
	projectBuildCmd.Flags().StringVarP(&projectBuildAbcFileDir, "abc-file-dir", "a", "", "The directory to find the ABC files")
	projectBuildCmd.Flags().IntVarP(&projectBuildPriorityThreshold, "priority-threshold", "p", 1, "The maximum priority of songs to include in the build")
	projectBuildCmd.Flags().StringVarP(&projectSampleId, "sampleId", "s", projectSampleId, "A string to indentify the sample stage. Will be injected to the project config")

}
