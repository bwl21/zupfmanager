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

	if !strings.Contains(htmlContent, "<title>Inhaltsverzeichnis - TEST</title>") {
		t.Error("HTML should contain correct title with project name")
	}

	if !strings.Contains(htmlContent, "<h1>Inhaltsverzeichnis</h1>") {
		t.Error("HTML should contain main heading")
	}

	if !strings.Contains(htmlContent, "Notensammlung TEST") {
		t.Error("HTML should contain project subtitle")
	}

	// Validate table structure
	if !strings.Contains(htmlContent, "<table class=\"toc-table\">") {
		t.Error("HTML should contain TOC table")
	}

	if !strings.Contains(htmlContent, "<th class=\"toc-number\">Nr.</th>") {
		t.Error("HTML should contain table headers")
	}

	// Validate song entries
	expectedEntries := []string{
		"<td class=\"toc-number\">01</td>",
		"<td class=\"toc-title\">Test Song 1</td>",
		"<td class=\"toc-info\">Komponist A</td>",
		"<td class=\"toc-number\">02</td>",
		"<td class=\"toc-title\">Test Song 2</td>", 
		"<td class=\"toc-info\">Komponist B</td>",
		"<td class=\"toc-number\">03</td>",
		"<td class=\"toc-title\">Test Song 3</td>",
		"<td class=\"toc-info\"></td>", // Empty for song without tocinfo
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
		"toc-table",
	}

	for _, class := range expectedClasses {
		if !strings.Contains(htmlContent, class) {
			t.Errorf("HTML should contain CSS class '%s'", class)
		}
	}

	// Validate print styles
	if !strings.Contains(htmlContent, "@media print") {
		t.Error("HTML should contain print-specific CSS")
	}

	if !strings.Contains(htmlContent, "@page") {
		t.Error("HTML should contain page formatting CSS")
	}

	t.Logf("✅ HTML TOC created successfully at: %s", htmlFile)
}