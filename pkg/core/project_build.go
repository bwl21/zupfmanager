package core

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
	"strings"

	"dario.cat/mergo"
	"github.com/bwl21/zupfmanager/internal/ent"
	entprojectsong "github.com/bwl21/zupfmanager/internal/ent/projectsong"
	"github.com/bwl21/zupfmanager/internal/zupfnoter"
	"golang.org/x/sync/errgroup"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

const (
	zupfnoterConfigString = "%%%%zupfnoter.config"
)

// ProgressCallback is a function type for progress updates
type ProgressCallback func(progress int, message string)

// ExecuteProjectBuild performs the actual project build logic
func (s *projectService) ExecuteProjectBuild(ctx context.Context, req BuildProjectRequest) error {
	return s.ExecuteProjectBuildWithProgress(ctx, req, nil)
}

// ExecuteProjectBuildWithProgress performs the actual project build logic with progress updates
func (s *projectService) ExecuteProjectBuildWithProgress(ctx context.Context, req BuildProjectRequest, progressCallback ProgressCallback) error {
	// Get the project first
	project, err := s.db.Project.Get(ctx, req.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Then query the project songs separately with the priority filter (like the working version)
	projectSongs, err := s.db.ProjectSong.Query().
		Where(
			entprojectsong.And(
				entprojectsong.ProjectID(req.ProjectID),
				entprojectsong.PriorityLTE(req.PriorityThreshold),
			),
		).
		WithSong().
		WithProject().
		Order(ent.Asc("priority")).
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed to query project songs: %w", err)
	}

	// Attach the songs to the project for compatibility with the rest of the code
	project.Edges.ProjectSongs = projectSongs

	return s.buildProject(ctx, req.AbcFileDir, req.OutputDir, project, req.SampleID, progressCallback)
}

func (s *projectService) buildProject(ctx context.Context, abcFileDir, outputDir string, project *ent.Project, sampleId string, progressCallback ProgressCallback) error {
	updateProgress := func(progress int, message string) {
		if progressCallback != nil {
			progressCallback(progress, message)
		}
	}

	updateProgress(15, "Preparing directories")
	
	// Remove existing directories
	for _, dir := range []string{"pdf", "abc", "log", "druckdateien", "referenz"} {
		if err := os.RemoveAll(filepath.Join(outputDir, dir)); err != nil {
			slog.Error("Failed to remove directory", "directory", dir, "error", err)
			return fmt.Errorf("failed to remove directory %s: %w", dir, err)
		}
	}

	// Create base directories
	dirs := []string{
		outputDir,
		filepath.Join(outputDir, "pdf"),
		filepath.Join(outputDir, "abc"),
		filepath.Join(outputDir, "log"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			slog.Error("Failed to create directory", "directory", dir, "error", err)
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create druckdateien directory
	druckdateienDir := filepath.Join(outputDir, "druckdateien")
	_ = os.MkdirAll(druckdateienDir, 0755)

	// Create target directories based on folder patterns
	folderPatterns := s.getFolderPatterns(project)
	folderSet := make(map[string]bool)
	for _, folder := range folderPatterns {
		if !folderSet[folder] {
			targetDir := filepath.Join(druckdateienDir, folder)
			_ = os.MkdirAll(targetDir, 0755)
			folderSet[folder] = true
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(5)

	projectSongs := project.Edges.ProjectSongs
	slog.Info("Project songs loaded", "count", len(projectSongs))
	updateProgress(25, fmt.Sprintf("Building %d songs", len(projectSongs)))

	// Ensure all songs are loaded
	for _, ps := range projectSongs {
		if ps.Edges.Song == nil {
			return fmt.Errorf("song not loaded for project song %d", ps.ID)
		}
	}

	// Sort the songs by title
	sort.Slice(projectSongs, func(i, j int) bool {
		return strings.ToLower(projectSongs[i].Edges.Song.Title) < strings.ToLower(projectSongs[j].Edges.Song.Title)
	})

	// Track completed songs for progress
	completedSongs := 0
	totalSongs := len(projectSongs)

	for id, song := range projectSongs {
		song := song
		songIndex := id
		eg.Go(func() error {
			err := s.buildSong(ctx, abcFileDir, outputDir, songIndex+1, song, sampleId, project)
			if err == nil {
				completedSongs++
				// Progress from 25% to 70% based on song completion
				progress := 25 + (completedSongs * 45 / totalSongs)
				updateProgress(progress, fmt.Sprintf("Built song %d/%d: %s", completedSongs, totalSongs, song.Edges.Song.Title))
			}
			return err
		})
	}
	err := eg.Wait()
	if err != nil {
		return fmt.Errorf("failed to build songs: %w", err)
	}

	updateProgress(75, "Processing copyright information")

	copyrightNames := s.getCopyrightNames(project)
	slog.Info("Copyright Names", "names", copyrightNames)

	err = s.createCopyrightDirectories(outputDir, copyrightNames)
	if err != nil {
		return fmt.Errorf("failed to create copyright directories: %w", err)
	}

	err = s.copyPdfsToCopyrightDirectories(project, outputDir)
	if err != nil {
		return fmt.Errorf("failed to copy PDFs to copyright directories: %w", err)
	}

	updateProgress(80, "Creating table of contents")
	if err := s.createToc(context.Background(), project, projectSongs, outputDir); err != nil {
		return fmt.Errorf("failed to create table of contents: %w", err)
	}

	updateProgress(85, "Merging PDF files")
	// Get folder patterns from project config or use defaults
	folderPatterns = s.getFolderPatterns(project)

	// Get unique target folders
	folderSet = make(map[string]bool)
	for _, targetFolder := range folderPatterns {
		folderSet[targetFolder] = true
	}

	// Merge PDFs for each target folder
	for folder := range folderSet {
		sourceDir := filepath.Join(outputDir, "druckdateien", folder)
		// Add project short name to the output filename
		destFile := filepath.Join(outputDir, "druckdateien", fmt.Sprintf("%s_%s.pdf", project.ShortName, folder))

		slog.Info("Merging PDFs for folder", "folder", folder, "source", sourceDir, "dest", destFile)

		err = s.mergePDFs(sourceDir, destFile)
		if err != nil {
			return fmt.Errorf("failed to merge PDFs in %s directory: %w", folder, err)
		}
	}

	return nil
}

// Rest of the helper functions...
// (I'll add them in the next step to keep this manageable)

// getCopyrightNames returns a slice of copyright names used in the project.
func (s *projectService) getCopyrightNames(project *ent.Project) []string {
	copyrightNames := make([]string, 0)
	for _, ps := range project.Edges.ProjectSongs {
		if ps.Edges.Song.Copyright != "" {
			copyrightNames = append(copyrightNames, ps.Edges.Song.Copyright)
		}
	}
	return copyrightNames
}

func (s *projectService) createToc(ctx context.Context, project *ent.Project, projectSongs []*ent.ProjectSong, outputDir string) error {
	tocabc := ""
	for id, song := range projectSongs {
		tocinfo := ""
		if song.Edges.Song.Tocinfo != "" {
			tocinfo = " - " + song.Edges.Song.Tocinfo
		}

		tocabc += fmt.Sprintf("W:%02d %s%s\n", id+1, song.Edges.Song.Title, tocinfo)
	}

	templateFile := filepath.Join(project.ShortName, "tpl", "999_inhaltsverzeichnis_template.abc")
	toctemplateBytes, err := os.ReadFile(templateFile)
	if err != nil {
		slog.Warn("failed to read template file, using default", "path", templateFile, "error", err)
		defaultTemplateFile := "x/MBT-2025/999_inhaltsverzeichnis_template.abc"
		toctemplateBytes, err = os.ReadFile(defaultTemplateFile)
		if err != nil {
			slog.Warn("failed to read default template file, using built-in template", "path", defaultTemplateFile, "error", err)
			// Use built-in template as last resort
			toctemplateBytes = []byte(`X:1
 sollte T:Inhaltsverzeichnis
M:4/4
L:1/4
K:C
W:{{TOC}}
`)
		} else {
			slog.Warn("using default template file", "path", defaultTemplateFile)
		}
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
	json.NewEncoder(tempFile).Encode("{}")
	tempFile.Close()

	stdOutBuf, stdErrBuf, err := zupfnoter.Run(ctx, filepath.Join(outputDir, "abc", tocSongFilename), filepath.Join(outputDir, "pdf"))
	if err != nil {
		errorMsg := fmt.Sprintf("Zupfnoter failed for TOC %s", tocSongFilename)
		if stdOutBuf != "" {
			errorMsg += fmt.Sprintf("\nStdout: %s", stdOutBuf)
		}
		if stdErrBuf != "" {
			errorMsg += fmt.Sprintf("\nStderr: %s", stdErrBuf)
		}
		slog.Error("zupfnoter failed for TOC", "output", stdOutBuf, "stderr", stdErrBuf, "file", tocSongFilename)
		return fmt.Errorf("%s: %w", errorMsg, err)
	}

	// Distribute the table of contents PDF to the print files directories.
	err = s.distributeZupfnoterOutput(project, tocSongFilename, outputDir, 0)
	if err != nil {
		return fmt.Errorf("failed to distribute Zupfnoter output: %w", err)
	}

	return nil
}

func (s *projectService) buildSong(ctx context.Context, abcFileDir, outputDir string, songIndex int, song *ent.ProjectSong, projectSampleId string, project *ent.Project) error {
	slog.Info("building song", "song", song.Edges.Song.Title)

	abcFile, err := os.ReadFile(filepath.Join(abcFileDir, song.Edges.Song.Filename))
	if err != nil {
		return fmt.Errorf("failed to read ABC file: %w", err)
	}

	fileConfig, err := s.extractConfigFromABCFile(abcFile)
	if err != nil {
		return fmt.Errorf("failed to extract config from ABC file: %w", err)
	}

	projectConfigBytes, err := json.Marshal(song.Edges.Project.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal project config: %w", err)
	}
	fc := bytes.ReplaceAll(projectConfigBytes, []byte("#{PREFIX}"), []byte(song.Edges.Project.ShortName))
	fc = bytes.ReplaceAll(fc, []byte("#{the_index}"), []byte(fmt.Sprintf("%02d", songIndex)))
	fc = bytes.ReplaceAll(fc, []byte("#{sampleId}"), []byte(projectSampleId))

	var finalConfig map[string]any
	err = json.Unmarshal(fc, &finalConfig)
	if err != nil {
		return fmt.Errorf("failed to unmarshal project config: %w", err)
	}

	err = mergo.Merge(&finalConfig, fileConfig)
	if err != nil {
		return fmt.Errorf("failed to merge config: %w", err)
	}

	tempConfigFile, err := os.CreateTemp("", "zupfnoter-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	json.NewEncoder(tempConfigFile).Encode(finalConfig)
	tempConfigFile.Close()

	stdOutBuf, stdErrBuf, err := zupfnoter.Run(
		ctx,
		filepath.Join(abcFileDir, song.Edges.Song.Filename),
		filepath.Join(outputDir, "pdf"),
		tempConfigFile.Name(),
	)
	if err != nil {
		errorMsg := fmt.Sprintf("Zupfnoter failed for %s", song.Edges.Song.Filename)
		if stdOutBuf != "" {
			errorMsg += fmt.Sprintf("\nStdout: %s", stdOutBuf)
		}
		if stdErrBuf != "" {
			errorMsg += fmt.Sprintf("\nStderr: %s", stdErrBuf)
		}
		slog.Error("zupfnoter failed", "output", stdOutBuf, "stderr", stdErrBuf, "file", song.Edges.Song.Filename)
		return fmt.Errorf("%s: %w", errorMsg, err)
	}
	os.Remove(tempConfigFile.Name())

	err = s.distributeZupfnoterOutput(project, song.Edges.Song.Filename, outputDir, songIndex)
	if err != nil {
		return fmt.Errorf("failed to distribute Zupfnoter output: %w", err)
	}

	err = os.WriteFile(
		filepath.Join(outputDir, "abc", song.Edges.Song.Filename),
		abcFile,
		0644,
	)
	if err != nil {
		return fmt.Errorf("failed to copy ABC file to output dir: %w", err)
	}

	logFN := fmt.Sprintf("%s.err.log", song.Edges.Song.Filename)
	err = os.Rename(
		filepath.Join(outputDir, "pdf", logFN),
		filepath.Join(outputDir, "log", logFN),
	)
	if err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	return nil
}

func (s *projectService) extractConfigFromABCFile(abcFile []byte) (map[string]any, error) {
	configLine := bytes.Index(abcFile, []byte(zupfnoterConfigString))
	if configLine == -1 {
		return make(map[string]any), nil
	}
	config := bytes.TrimSpace(abcFile[configLine+len(zupfnoterConfigString):])

	var configMap map[string]any
	err := json.Unmarshal(config, &configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return configMap, nil
}

// Add remaining helper functions...

func (s *projectService) getFolderPatterns(project *ent.Project) map[string]string {
	folderPatterns := map[string]string{
		"*_-A*_a3.pdf": "klein",
		"*_-M*_a3.pdf": "klein",
		"*_-O*_a3.pdf": "klein",
		"*_-B*_a3.pdf": "gross",
		"*_-X*_a3.pdf": "gross",
	}

	if configPatterns, ok := project.Config["folderPatterns"].(map[string]interface{}); ok {
		folderPatterns = make(map[string]string)
		for pattern, folder := range configPatterns {
			if folderStr, ok := folder.(string); ok {
				folderPatterns[pattern] = folderStr
			}
		}
	}

	return folderPatterns
}

func (s *projectService) distributeZupfnoterOutput(project *ent.Project, baseFilename string, outputDir string, songIndex int) error {
	pdfDir := filepath.Join(outputDir, "pdf")
	baseFilenameWithoutExt := strings.TrimSuffix(baseFilename, ".abc")
	pattern := filepath.Join(pdfDir, filepath.Base(baseFilenameWithoutExt)+"*.pdf")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to glob PDF files: %w", err)
	}

	folderPatterns := s.getFolderPatterns(project)

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
			slog.Error("no target folder found", "filename", filename)
			continue
		}

		err := os.MkdirAll(targetDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create target directory: %w", err)
		}

		targetFile := filepath.Join(targetDir, newFilename)
		err = s.copyFile(pdfFile, targetFile)
		if err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}
	}

	return nil
}

func (s *projectService) copyFile(src, dst string) error {
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

func (s *projectService) mergePDFs(dir, dest string) error {
	slog.Info("merging pdf files", "dir", dir, "dest", dest)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		slog.Warn("directory does not exist, skipping merge", "dir", dir)
		return nil
	}

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

	if len(files) == 0 {
		slog.Warn("no PDF files found to merge", "dir", dir)
		return nil
	}

	slog.Info("merging PDF files", "count", len(files), "from", dir, "to", dest)
	err = api.MergeCreateFile(files, dest, false, nil)
	if err != nil {
		return fmt.Errorf("failed to merge pdf files: %w", err)
	}

	slog.Info("successfully merged PDF files", "dest", dest)
	return nil
}

func (s *projectService) createCopyrightDirectory(outputDir string, copyrightName string) error {
	dirPath := filepath.Join(outputDir, "referenz", copyrightName)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *projectService) createCopyrightDirectories(outputDir string, copyrightNames []string) error {
	for _, copyrightName := range copyrightNames {
		err := s.createCopyrightDirectory(outputDir, copyrightName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *projectService) copyPdfsToCopyrightDirectories(project *ent.Project, outputDir string) error {
	if project.Edges.ProjectSongs == nil {
		return fmt.Errorf("project songs not loaded")
	}

	copyrightGroups := make(map[string][]*ent.ProjectSong)
	for _, ps := range project.Edges.ProjectSongs {
		if ps.Edges.Song.Copyright != "" {
			copyrightGroups[ps.Edges.Song.Copyright] = append(copyrightGroups[ps.Edges.Song.Copyright], ps)
		}
	}

	for copyrightName, songs := range copyrightGroups {
		copyrightDir := filepath.Join(outputDir, "referenz", copyrightName)

		for _, song := range songs {
			sourcePattern := filepath.Join(outputDir, "pdf", strings.TrimSuffix(song.Edges.Song.Filename, ".abc")+"*.pdf")
			files, err := filepath.Glob(sourcePattern)
			if err != nil {
				return fmt.Errorf("failed to glob files for song %s: %w", song.Edges.Song.Title, err)
			}

			for _, file := range files {
				filename := filepath.Base(file)
				destFile := filepath.Join(copyrightDir, filename)
				err = s.copyFile(file, destFile)
				if err != nil {
					return fmt.Errorf("failed to copy file %s to %s: %w", file, destFile, err)
				}
			}
		}
	}

	return nil
}
