package cmd

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"dario.cat/mergo"
	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/bwl21/zupfmanager/internal/ent"
	"github.com/bwl21/zupfmanager/internal/ent/project"
	"github.com/bwl21/zupfmanager/internal/ent/projectsong"
	"github.com/bwl21/zupfmanager/internal/zupfnoter"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

const (
	zupfnoterConfigString = "%%%%zupfnoter.config"
)

var (
	projectBuildOutputDir         string
	projectBuildAbcFileDir        string
	projectBuildPriorityThreshold int
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

		project, err := client.Project.Query().Where(project.ID(projectID)).
			WithProjectSongs(func(psq *ent.ProjectSongQuery) {
				psq.Where(projectsong.PriorityLTE(projectBuildPriorityThreshold))
				psq.WithProject().WithSong()
			}).
			First(context.Background())
		if err != nil {
			return fmt.Errorf("failed to find project with ID %d: %w", projectID, err)
		}

		projectBuildOutputDir = project.ShortName
		cmd.Flags().StringVarP(&projectBuildOutputDir, "output-dir", "o", projectBuildOutputDir, "The directory to output the build results")

		if projectBuildAbcFileDir == "" {
			// Check if abc_file_dir exists and is a string
			abcFileDir, ok := project.Config["abc_file_dir"].(string)
			if ok {
				projectBuildAbcFileDir = abcFileDir
			} else {
				// Provide a default value or handle the error appropriately
				projectBuildAbcFileDir = ""
			}
		}

		return buildProject(projectBuildAbcFileDir, projectBuildOutputDir, project)
	},
}

func buildProject(abcFileDir, outputDir string, project *ent.Project) error {
	os.RemoveAll(filepath.Join(outputDir, "pdf"))
	os.RemoveAll(filepath.Join(outputDir, "abc"))
	os.RemoveAll(filepath.Join(outputDir, "log"))
	os.RemoveAll(filepath.Join(outputDir, "druckdateien"))
	os.RemoveAll(filepath.Join(outputDir, "referenz"))

	_ = os.MkdirAll(outputDir, 0755)
	_ = os.MkdirAll(filepath.Join(outputDir, "pdf"), 0755)
	_ = os.MkdirAll(filepath.Join(outputDir, "abc"), 0755)
	_ = os.MkdirAll(filepath.Join(outputDir, "log"), 0755)
	druckdateienDir := filepath.Join(outputDir, "druckdateien")
	_ = os.MkdirAll(druckdateienDir, 0755)
	grossDir := filepath.Join(druckdateienDir, "gross")
	_ = os.MkdirAll(grossDir, 0755)
	kleinDir := filepath.Join(druckdateienDir, "klein")
	_ = os.MkdirAll(kleinDir, 0755)

	eg, ctx := errgroup.WithContext(context.Background())
	eg.SetLimit(5)

	projectSongs := project.Edges.ProjectSongs
	sort.Slice(projectSongs, func(i, j int) bool {
		return strings.ToLower(projectSongs[i].Edges.Song.Title) < strings.ToLower(projectSongs[j].Edges.Song.Title)
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

	copyrightNames := getCopyrightNames(project)
	fmt.Println("Copyright Names:", copyrightNames)

	err = createCopyrightDirectories(copyrightNames)
	if err != nil {
		return fmt.Errorf("failed to create copyright directories: %w", err)
	}

	err = copyPdfsToCopyrightDirectories(project, outputDir)
	if err != nil {
		return fmt.Errorf("failed to copy PDFs to copyright directories: %w", err)
	}

	err2 := createToc(project, projectSongs, err, outputDir)
	if err2 != nil {
		return err2
	}
	//os.Remove(tempFile.Name())

	grossDir = filepath.Join(outputDir, "druckdateien", "gross")
	kleinDir = filepath.Join(outputDir, "druckdateien", "klein")

	err = mergePDFs(grossDir, filepath.Join(outputDir, "druckdateien", "gross.pdf"))
	if err != nil {
		return fmt.Errorf("failed to merge PDFs in gross directory: %w", err)
	}

	err = mergePDFs(kleinDir, filepath.Join(outputDir, "druckdateien", "klein.pdf"))
	if err != nil {
		return fmt.Errorf("failed to merge PDFs in klein directory: %w", err)
	}

	return nil
}

// getCopyrightNames returns a slice of copyright names used in the project.
func getCopyrightNames(project *ent.Project) []string {
	copyrightNames := make([]string, 0)
	for _, ps := range project.Edges.ProjectSongs {
		if ps.Edges.Song.Copyright != "" {
			copyrightNames = append(copyrightNames, ps.Edges.Song.Copyright)
		}
	}
	return copyrightNames
}

func createToc(project *ent.Project, projectSongs []*ent.ProjectSong, err error, outputDir string) error {
	tocabc := ""
	for id, song := range projectSongs {
		tocabc += fmt.Sprintf("W:%d %s\n", id+1, song.Edges.Song.Title)
	}

	templateFile := filepath.Join(project.ShortName, "tpl", "999_inhaltsverzeichnis_template.abc")
	toctemplateBytes, err := os.ReadFile(templateFile)
	if err != nil {
		slog.Warn("failed to read template file, using default", "path", templateFile, "error", err)
		defaultTemplateFile := "x/MBT-2025/999_inhaltsverzeichnis_template.abc"
		toctemplateBytes, err = os.ReadFile(defaultTemplateFile)
		if err != nil {
			return fmt.Errorf("failed to read default template file: %w", err)
		}
		slog.Warn("using default template file", "path", defaultTemplateFile)
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

	_, _, err = zupfnoter.Run(ctxb, filepath.Join(outputDir, "abc", tocSongFilename), filepath.Join(outputDir, "pdf"))
	if err != nil {
		fmt.Println(filepath.Join(outputDir, "abc", tocSongFilename))
		return fmt.Errorf("failed to run zupfnoter: %w [%s]", err, tocSongFilename)
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
	projectConfigBytes, err := json.Marshal(song.Edges.Project.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal project config: %w", err)
	}
	fc := bytes.ReplaceAll(projectConfigBytes, []byte("#{PREFIX}"), []byte(song.Edges.Project.ShortName))
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
	stdOutBuf, _, err := zupfnoter.Run(
		ctx,
		filepath.Join(abcFileDir, song.Edges.Song.Filename), // Pfad zur ABC-Datei
		filepath.Join(outputDir, "pdf"),                     // Ausgabeverzeichnis für PDF
		tempConfigFile.Name(),                               // Pfad zur Konfigurationsdatei
	)
	if err != nil {
		fmt.Println(stdOutBuf)
		fmt.Println(filepath.Join(outputDir, "abc", song.Edges.Song.Filename))
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
	projectBuildCmd.Flags().IntVarP(&projectBuildPriorityThreshold, "priority-threshold", "p", 1, "The maximum priority of songs to include in the build")
}

func distributeZupfnoterOutput(baseFilename string, outputDir string, songIndex int) error {
	pdfDir := filepath.Join(outputDir, "pdf")
	baseFilenameWithoutExt := strings.TrimSuffix(baseFilename, ".abc")
	pattern := filepath.Join(pdfDir, filepath.Base(baseFilenameWithoutExt)+"*.pdf")
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
			slog.Error("no target folder found: ", "filename", filename)
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

func mergePDFs(dir, dest string) error {
	slog.Info("merging pdf files", "dir", dir, "dest", dest)

	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".pdf" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	err = api.MergeCreateFile(files, dest, false, nil)
	if err != nil {
		return fmt.Errorf("failed to merge pdf files: %w", err)
	}

	return nil
}

// createCopyrightDirectory creates a directory for a given copyright name under the "referenz" directory.
func createCopyrightDirectory(copyrightName string) error {
	dirPath := filepath.Join("referenz", copyrightName)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// createCopyrightDirectories creates directories for a given list of copyright names under the "referenz" directory.
func createCopyrightDirectories(copyrightNames []string) error {
	for _, copyrightName := range copyrightNames {
		err := createCopyrightDirectory(copyrightName)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyPdfsToCopyrightDirectories(project *ent.Project, outputDir string) error {
	// Sicherstellen, dass die Projekt-Songs geladen sind
	if project.Edges.ProjectSongs == nil {
		return fmt.Errorf("project songs not loaded")
	}

	// Map zur Gruppierung nach Copyright (vermeidet doppelte Verzeichniserstellung)
	copyrightMap := make(map[string][]*ent.ProjectSong)

	// Songs nach Copyright gruppieren
	for _, ps := range project.Edges.ProjectSongs {
		if ps.Edges.Song == nil {
			continue
		}
		copyright := ps.Edges.Song.Copyright
		if copyright == "" {
			continue
		}
		copyrightMap[copyright] = append(copyrightMap[copyright], ps)
	}

	// PDFs pro Copyright-Verzeichnis kopieren
	for copyright, songs := range copyrightMap {
		destDir := filepath.Join(outputDir, "referenz", copyright)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", destDir, err)
		}

		// Einmaliges Durchsuchen des PDF-Verzeichnisses
		pdfDir := filepath.Join(outputDir, "pdf")
		err := filepath.Walk(pdfDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() || !strings.HasSuffix(info.Name(), ".pdf") {
				return nil
			}

			// Prüfe auf Übereinstimmung mit einem der Songs
			for _, ps := range songs {
				if ps.Edges.Song == nil || ps.Edges.Song.Filename == "" {
					continue
				}

				baseName := strings.TrimSuffix(ps.Edges.Song.Filename, ".abc")
				if strings.Contains(info.Name(), baseName) {
					destPath := filepath.Join(destDir, info.Name())
					if err := copyFile(path, destPath); err != nil {
						return fmt.Errorf("failed to copy %s to %s: %w", path, destPath, err)
					}
				}
			}
			return nil
		})

		if err != nil {
			return err
		}
	}
	return nil
}
