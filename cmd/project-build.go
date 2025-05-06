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
	"strings"
	"io"

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

		projectBuildOutputDir = project.ShortName
		cmd.Flags().StringVarP(&projectBuildOutputDir, "output-dir", "o", projectBuildOutputDir, "The directory to output the build results")

		return buildProject(projectBuildAbcFileDir, projectBuildOutputDir, project)
	},
}

func buildProject(abcFileDir, outputDir string, project *ent.Project) error {
	os.RemoveAll(filepath.Join(outputDir, "pdf"))
	os.RemoveAll(filepath.Join(outputDir, "abc"))
	os.RemoveAll(filepath.Join(outputDir, "log"))
	os.RemoveAll(filepath.Join(outputDir, "druckdateien"))

	_ = os.MkdirAll(outputDir, 0755)
	_ = os.MkdirAll(filepath.Join(outputDir, "pdf"), 0755)
	_ = os.MkdirAll(filepath.Join(outputDir, "abc"), 0755)
	_ = os.MkdirAll(filepath.Join(outputDir, "log"), 0755)
	_ = os.MkdirAll(filepath.Join(outputDir, "druckdateien"), 0755)

	eg, ctx := errgroup.WithContext(context.Background())
	eg.SetLimit(5)

	projectSongs := project.Edges.ProjectSongs
	sort.Slice(projectSongs, func(i, j int) bool {
		return projectSongs[i].Edges.Song.Title < projectSongs[j].Edges.Song.Title
	})
	for id, song := range projectSongs {
		song := song
		eg.Go(func() error {
			return buildSong(ctx, abcFileDir, outputDir, id+1, song)
		})
	}
	err := eg.Wait()
	if err != nil {
		return fmt.Errorf("failed to build songs: %w", err)
	}

	err2 := createToc(projectSongs, err, outputDir)
	if err2 != nil {
		return err2
	}
	//os.Remove(tempFile.Name())

	return nil
}

func createToc(projectSongs []*ent.ProjectSong, err error, outputDir string) error {
	tocabc := ""
	for id, song := range projectSongs {
		tocabc += fmt.Sprintf("W:%d %s\n", id+1, song.Edges.Song.Title)
	}

	toctemplateBytes, err := os.ReadFile(filepath.Join(outputDir, "999_inhaltsverzeichnis_template.abc"))
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	toctemplate := strings.Replace(string(toctemplateBytes), "W:{{TOC}}", tocabc, 1)

	tocSongFilename := "00_inhaltsverzeichnis.abc"
	err = os.WriteFile(filepath.Join(outputDir, "abc", tocSongFilename), []byte(toctemplate), 0644)
	if err != nil {
		return fmt.Errorf("failed to write toc file: %w", err)
	}

	tempFile, err := os.CreateTemp("", "zupfnoter-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	// defer os.Remove(tempFile.Name())
	json.NewEncoder(tempFile).Encode("{}")
	tempFile.Close()

	ctxb := context.Background()
	err = zupfnoter.Run(ctxb, filepath.Join(outputDir, "abc", tocSongFilename), filepath.Join(outputDir, "pdf"))
	if err != nil {
		fmt.Println(filepath.Join(outputDir, "abc", tocSongFilename))
		return fmt.Errorf("failed to run zupfnoter: %w", err)
	}

	// Distribute the table of contents PDF to the print files directories.
	err = distributeZupfnoterOutput(tocSongFilename, outputDir, 0)
	if err != nil {
		return fmt.Errorf("failed to distribute Zupfnoter output: %w", err)
	}

	return nil
}

// buildSong verarbeitet einen Song: Liest die ABC-Datei, kombiniert Konfigurationen,
// ruft das externe Tool "zupfnoter" auf, kopiert die ABC-Datei ins Zielverzeichnis
// und verschiebt das Logfile.
func buildSong(ctx context.Context, abcFileDir, outputDir string, songIndex int, song *ent.ProjectSong) error {
	// 1. Logge den Start der Verarbeitung für diesen Song.
	slog.Info("building song", "song", song.Edges.Song.Title)

	// 2. Lese die ABC-Datei des Songs ein.
	abcFile, err := os.ReadFile(filepath.Join(abcFileDir, song.Edges.Song.Filename))
	if err != nil {
		return fmt.Errorf("failed to read ABC file: %w", err)
	}

	// 3. Extrahiere Konfiguration aus der ABC-Datei (z.B. Metadaten).
	fileConfig, err := extractConfigFromABCFile(abcFile)
	if err != nil {
		return fmt.Errorf("failed to extract config from ABC file: %w", err)
	}

	// 4. Serialisiere die Projekt-Konfiguration als JSON.
	fc, err := json.Marshal(song.Edges.Project.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal project config: %w", err)
	}

	// 5. Ersetze Platzhalter im JSON (z.B. #{PREFIX}, #{the_index}).
	fc = bytes.ReplaceAll(fc, []byte("#{PREFIX}"), []byte(song.Edges.Project.ShortName))
	fc = bytes.ReplaceAll(fc, []byte("#{the_index}"), []byte(fmt.Sprintf("%02d", songIndex)))

	// 6. Deserialisiere das JSON wieder in eine Map für weitere Bearbeitung.
	var finalConfig map[string]any
	err = json.Unmarshal(fc, &finalConfig)
	if err != nil {
		return fmt.Errorf("failed to unmarshal project config: %w", err)
	}

	// 7. Führe die Konfiguration aus der Datei und dem Projekt zusammen.
	err = mergo.Merge(&finalConfig, fileConfig)
	if err != nil {
		return fmt.Errorf("failed to merge config: %w", err)
	}

	// 8. Schreibe die finale Konfiguration in eine temporäre Datei.
	tempConfigFile, err := os.CreateTemp("", "zupfnoter-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	// Schreibe die finale Konfiguration als JSON in die Datei.
	json.NewEncoder(tempConfigFile).Encode(finalConfig)
	tempConfigFile.Close()

	// 9. Rufe das externe Tool "zupfnoter" auf und übergebe die notwendigen Dateien.
	err = zupfnoter.Run(
		ctx,
		filepath.Join(abcFileDir, song.Edges.Song.Filename), // Pfad zur ABC-Datei
		filepath.Join(outputDir, "pdf"),                     // Ausgabeverzeichnis für PDF
		tempConfigFile.Name(),                               // Pfad zur Konfigurationsdatei
	)
	if err != nil {
		return fmt.Errorf("failed to run zupfnoter: %w", err)
	}
	// 10. Lösche die temporäre Konfigurationsdatei.
	os.Remove(tempConfigFile.Name())

	// 11. Distribute the Zupfnoter output to the print files directories.
	err = distributeZupfnoterOutput(song.Edges.Song.Filename, outputDir, songIndex)
	if err != nil {
		return fmt.Errorf("failed to distribute Zupfnoter output: %w", err)
	}

	// 12. Kopiere die ABC-Datei ins Zielverzeichnis (z.B. für das Archiv).
	err = os.WriteFile(
		filepath.Join(outputDir, "abc", song.Edges.Song.Filename),
		abcFile,
		0644,
	)
	if err != nil {
		return fmt.Errorf("failed to copy ABC file to output dir: %w", err)
	}

	// 12. Verschiebe das Logfile ins Log-Verzeichnis.
	logFN := fmt.Sprintf("%s.err.log", song.Edges.Song.Filename)
	err = os.Rename(
		filepath.Join(outputDir, "pdf", logFN),
		filepath.Join(outputDir, "log", logFN),
	)
	if err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	// 13. Erfolgreich abgeschlossen.
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

	projectBuildCmd.Flags().StringVarP(&projectBuildAbcFileDir, "abc-file-dir", "a", "", "The directory to find the ABC files")
}

func distributeZupfnoterOutput(baseFilename string, outputDir string, songIndex int) error {
	pdfDir := filepath.Join(outputDir, "pdf")
	baseFilenameWithoutExt := strings.TrimSuffix(baseFilename, ".abc")
	pattern := filepath.Join(pdfDir, filepath.Base(baseFilenameWithoutExt) + "*.pdf")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to glob PDF files: %w", err)
	}

	folderPatterns := map[string]string{
		"*_-A*_a3.pdf": "klein",
		"*_-M*_a3.pdf": "klein",
		"*_-O*_a3.pdf": "klein",
		"*_-B*_a3.pdf": "gross",
		"*_-X*_a3.pdf": "gross",
	}

	for _, pdfFile := range files {
		filename := filepath.Base(pdfFile)
		newFilename := fmt.Sprintf("%02d_%s", songIndex, filename)
		var targetDir string

		for pattern, folder := range folderPatterns {
			matched, err := filepath.Match(pattern, filename)
			if err != nil {
				return fmt.Errorf("failed to match pattern: %w", err)
			}
			if matched {
				targetDir = filepath.Join(outputDir, "druckdateien", folder)
				break
			}
		}

		if targetDir == "" {
			slog.Info("skipping file", "filename", filename)
			continue
		}

		slog.Info("target directory", "targetDir", targetDir)
		err := os.MkdirAll(targetDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create target directory: %w", err)
		}

		targetFile := filepath.Join(targetDir, newFilename)
		slog.Info("copying file", "source", pdfFile, "target", targetFile)
		err = copyFile(pdfFile, targetFile)
		if err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
