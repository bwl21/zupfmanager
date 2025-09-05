package core

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bwl21/zupfmanager/internal/ent"
)

func TestCreateHTMLToc(t *testing.T) {
	// Create test project
	project := &ent.Project{
		ShortName: "TEST",
		Config:    make(map[string]interface{}),
	}

	// Create test songs
	projectSongs := []*ent.ProjectSong{
		{
			Edges: ent.ProjectSongEdges{
				Song: &ent.Song{
					Title:    "Test Song 1",
					Filename: "test1.abc",
					Tocinfo:  "Komponist A",
				},
			},
		},
		{
			Edges: ent.ProjectSongEdges{
				Song: &ent.Song{
					Title:    "Test Song 2",
					Filename: "test2.abc",
					Tocinfo:  "Komponist B",
				},
			},
		},
		{
			Edges: ent.ProjectSongEdges{
				Song: &ent.Song{
					Title:    "Test Song 3",
					Filename: "test3.abc",
					// No Tocinfo to test empty case
				},
			},
		},
	}

	// Create temporary output directory
	outputDir, err := os.MkdirTemp("", "html_toc_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// Create service
	service := &projectService{}

	// Test HTML TOC creation
	err = service.createHTMLToc(context.Background(), project, projectSongs, outputDir)
	if err != nil {
		t.Fatalf("Failed to create HTML TOC: %v", err)
	}

	// Check if HTML file was created
	htmlFile := filepath.Join(outputDir, "html", "00_inhaltsverzeichnis.html")
	if _, err := os.Stat(htmlFile); os.IsNotExist(err) {
		t.Fatalf("HTML TOC file was not created: %s", htmlFile)
	}

	// Read and validate HTML content
	content, err := os.ReadFile(htmlFile)
	if err != nil {
		t.Fatalf("Failed to read HTML file: %v", err)
	}

	htmlContent := string(content)

	// Validate HTML structure
	if !strings.Contains(htmlContent, "<!DOCTYPE html>") {
		t.Error("HTML should contain DOCTYPE declaration")
	}

	if !strings.Contains(htmlContent, "<title>Inhaltsverzeichnis</title>") {
		t.Error("HTML should contain correct title")
	}

	if !strings.Contains(htmlContent, "<h1>Inhaltsverzeichnis</h1>") {
		t.Error("HTML should contain main heading")
	}

	// Validate song entries
	expectedEntries := []string{
		"01</span>",
		"Test Song 1",
		"Komponist A",
		"02</span>",
		"Test Song 2", 
		"Komponist B",
		"03</span>",
		"Test Song 3",
	}

	for _, expected := range expectedEntries {
		if !strings.Contains(htmlContent, expected) {
			t.Errorf("HTML should contain '%s'", expected)
		}
	}

	// Validate CSS classes
	expectedClasses := []string{
		"toc-entry",
		"toc-number",
		"toc-title",
		"toc-info",
	}

	for _, class := range expectedClasses {
		if !strings.Contains(htmlContent, class) {
			t.Errorf("HTML should contain CSS class '%s'", class)
		}
	}

	t.Logf("✅ HTML TOC created successfully at: %s", htmlFile)
}