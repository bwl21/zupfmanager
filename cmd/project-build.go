package cmd

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"dario.cat/mergo"
	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/bwl21/zupfmanager/internal/ent"
	"github.com/bwl21/zupfmanager/internal/ent/project"
	"github.com/bwl21/zupfmanager/internal/zupfnoter"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

const (
	zupfnoterConfigString = "%%%%zupfnoter.config"
)

var (
	projectBuildOutputDir  string
	projectBuildAbcFileDir string
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

		project, err := client.Project.Query().Where(project.ID(projectID)).WithProjectSongs(
			func(psq *ent.ProjectSongQuery) {
				psq.WithProject().WithSong()
			},
		).First(context.Background())
		if err != nil {
			return fmt.Errorf("failed to find project with ID %d: %w", projectID, err)
		}

		return buildProject(projectBuildAbcFileDir, projectBuildOutputDir, project)
	},
}

func buildProject(abcFileDir, outputDir string, project *ent.Project) error {
	_ = os.MkdirAll(outputDir, 0755)
	_ = os.MkdirAll(filepath.Join(outputDir, "pdf"), 0755)
	_ = os.MkdirAll(filepath.Join(outputDir, "abc"), 0755)
	_ = os.MkdirAll(filepath.Join(outputDir, "log"), 0755)

	eg, ctx := errgroup.WithContext(context.Background())
	eg.SetLimit(5)

	projectSongs := project.Edges.ProjectSongs
	sort.Slice(projectSongs, func(i, j int) bool {
		return projectSongs[i].Edges.Song.Title < projectSongs[j].Edges.Song.Title
	})
	for id, song := range projectSongs {
		song := song
		eg.Go(func() error {
			return buildSong(ctx, abcFileDir, outputDir, id, song)
		})
	}
	err := eg.Wait()
	if err != nil {
		return fmt.Errorf("failed to build songs: %w", err)
	}

	// copy PDF files to output dir

	return nil
}

func buildSong(ctx context.Context, abcFileDir, outputDir string, songIndex int, song *ent.ProjectSong) error {
	slog.Info("building song", "song", song.Edges.Song.Title)

	abcFile, err := os.ReadFile(filepath.Join(abcFileDir, song.Edges.Song.Filename))
	if err != nil {
		return fmt.Errorf("failed to read ABC file: %w", err)
	}
	fileConfig, err := extractConfigFromABCFile(abcFile)
	if err != nil {
		return fmt.Errorf("failed to extract config from ABC file: %w", err)
	}

	fc, err := json.Marshal(song.Edges.Project.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal project config: %w", err)
	}

	fc = bytes.ReplaceAll(fc, []byte("#{PREFIX}"), []byte(song.Edges.Project.ShortName))
	fc = bytes.ReplaceAll(fc, []byte("#{the_index}"), []byte(fmt.Sprintf("%03d", songIndex)))

	var finalConfig map[string]any
	err = json.Unmarshal(fc, &finalConfig)
	if err != nil {
		return fmt.Errorf("failed to unmarshal project config: %w", err)
	}
	err = mergo.Merge(&finalConfig, fileConfig)
	if err != nil {
		return fmt.Errorf("failed to merge config: %w", err)
	}

	tempFile, err := os.CreateTemp("", "zupfnoter-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	// defer os.Remove(tempFile.Name())
	json.NewEncoder(tempFile).Encode(finalConfig)
	tempFile.Close()

	// node "#{ZUPFNOTER}" "#{srcfile}" "#{workdir}/pdf" x.json
	err = zupfnoter.Run(ctx, filepath.Join(abcFileDir, song.Edges.Song.Filename), filepath.Join(outputDir, "pdf"), tempFile.Name())
	if err != nil {
		return fmt.Errorf("failed to run zupfnoter: %w", err)
	}
	os.Remove(tempFile.Name())

	// copy ABC file to the output dir
	err = os.WriteFile(filepath.Join(outputDir, "abc", song.Edges.Song.Filename), abcFile, 0644)
	if err != nil {
		return fmt.Errorf("failed to copy ABC file to output dir: %w", err)
	}

	logFN := fmt.Sprintf("%s.err.log", song.Edges.Song.Filename)
	err = os.Rename(filepath.Join(outputDir, "pdf", logFN), filepath.Join(outputDir, "log", logFN))
	if err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	return nil
}

func extractConfigFromABCFile(abcFile []byte) (map[string]any, error) {
	// search until the first line that starts with %%%%zupfnoter.config and everything after that is JSON
	configLine := bytes.Index(abcFile, []byte(zupfnoterConfigString))
	if configLine == -1 {
		return nil, fmt.Errorf("no config found in ABC file")
	}
	config := bytes.TrimSpace(abcFile[configLine+len(zupfnoterConfigString):])

	// parse the JSON
	var configMap map[string]any
	err := json.Unmarshal(config, &configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return configMap, nil
}

func init() {
	projectCmd.AddCommand(projectBuildCmd)

	projectBuildCmd.Flags().StringVarP(&projectBuildOutputDir, "output-dir", "o", "output/", "The directory to output the build results")
	projectBuildCmd.Flags().StringVarP(&projectBuildAbcFileDir, "abc-file-dir", "a", "", "The directory to find the ABC files")
}
