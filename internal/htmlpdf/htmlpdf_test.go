package htmlpdf

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bwl21/zupfmanager/internal/ent"
)

func createTestHTML(content string) (string, error) {
	tempFile, err := os.CreateTemp("", "test-*.html")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	_, err = tempFile.WriteString(content)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

func createTestSong(title string) *ent.ProjectSong {
	return &ent.ProjectSong{
		Edges: ent.ProjectSongEdges{
			Song: &ent.Song{
				Title:    title,
				Filename: "test_song.abc",
			},
		},
	}
}

func TestPageNumberInjector(t *testing.T) {
	injector := NewPageNumberInjector("bottom-right")

	request := &ConversionRequest{
		SongIndex:  42,
		Song:       createTestSong("Test Song"),
		DOMScripts: make([]string, 0),
	}

	err := injector.InjectIntoDOM(context.Background(), request)
	if err != nil {
		t.Fatalf("InjectIntoDOM failed: %v", err)
	}

	if len(request.DOMScripts) != 2 {
		t.Errorf("Expected 2 DOM scripts (CSS + element), got %d", len(request.DOMScripts))
	}

	// Validate that page number is present
	foundPageNumber := false
	foundCSS := false
	for _, script := range request.DOMScripts {
		if strings.Contains(script, "42") {
			foundPageNumber = true
		}
		if strings.Contains(script, "pageStyle") || strings.Contains(script, "druckParagraph") {
			foundCSS = true
		}
	}

	if !foundPageNumber {
		t.Error("Page number should be injected into DOM scripts")
	}
	if !foundCSS {
		t.Error("CSS or page element should be present")
	}
}

func TestTextCleanupInjector(t *testing.T) {
	injector := NewTextCleanupInjector("#vb", "#debug")

	request := &ConversionRequest{
		DOMScripts: make([]string, 0),
	}

	err := injector.InjectIntoDOM(context.Background(), request)
	if err != nil {
		t.Fatalf("InjectIntoDOM failed: %v", err)
	}

	if len(request.DOMScripts) != 2 {
		t.Errorf("Expected 2 DOM scripts, got %d", len(request.DOMScripts))
	}

	// Validate that both patterns are in the scripts
	script := strings.Join(request.DOMScripts, " ")
	if !strings.Contains(script, "#vb") {
		t.Error("Script should contain #vb pattern")
	}
	if !strings.Contains(script, "#debug") {
		t.Error("Script should contain #debug pattern")
	}
	if !strings.Contains(script, "remove()") {
		t.Error("Script should contain remove() call")
	}
}

func TestCustomDOMInjector(t *testing.T) {
	injector := NewCustomDOMInjector()
	injector.AddCSS("body { background: red; }")
	injector.AddElement(HTMLElement{
		Tag:      "div",
		ID:       "test",
		Content:  "Song: ${SONG_TITLE}",
		Position: BodyStart,
	})
	injector.AddCleanupRule(CleanupRule{
		Selector: "span",
		Action:   "remove",
		Pattern:  "unwanted",
	})

	request := &ConversionRequest{
		SongIndex:  5,
		Song:       createTestSong("My Test Song"),
		DOMScripts: make([]string, 0),
	}

	err := injector.InjectIntoDOM(context.Background(), request)
	if err != nil {
		t.Fatalf("InjectIntoDOM failed: %v", err)
	}

	if len(request.DOMScripts) != 3 {
		t.Errorf("Expected 3 DOM scripts (cleanup + CSS + element), got %d", len(request.DOMScripts))
	}

	script := strings.Join(request.DOMScripts, " ")
	if !strings.Contains(script, "background: red") {
		t.Error("CSS should be injected")
	}
	if !strings.Contains(script, "My Test Song") {
		t.Error("Song title should be replaced in element content")
	}
	if !strings.Contains(script, "unwanted") {
		t.Error("Cleanup rule should be present")
	}
}

// Integration test - requires Chrome/Chromium to be installed
func TestChromeDPConverter_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create test HTML file (simulates Zupfnoter output)
	testHTML := `
	<!DOCTYPE html>
	<html>
	<head><title>Zupfnoter Generated Song</title></head>
	<body>
		<h1>Song Title</h1>
		<text>#vb</text>
		<div class="music-notation">
			<!-- Zupfnoter-generated music notation -->
			<p>Music content here</p>
		</div>
	</body>
	</html>
	`

	tempHTML, err := createTestHTML(testHTML)
	if err != nil {
		t.Fatalf("Failed to create test HTML: %v", err)
	}
	defer os.Remove(tempHTML)

	// Test conversion with DOM injection
	converter := NewChromeDPConverter(
		NewTextCleanupInjector("#vb"),
		NewPageNumberInjector("bottom-right"),
	)
	defer converter.Close()

	testSong := createTestSong("Test Song")
	outputPath := filepath.Join(os.TempDir(), "test_song_noten.pdf")
	defer os.Remove(outputPath)

	request := &ConversionRequest{
		HTMLFilePath: tempHTML,
		OutputPath:   outputPath,
		SongIndex:    5,
		Song:         testSong,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := converter.ConvertToPDF(ctx, request)
	if err != nil {
		t.Fatalf("ConvertToPDF failed: %v", err)
	}

	if result == nil {
		t.Fatal("Result should not be nil")
	}

	if result.FileSize <= 0 {
		t.Error("PDF file size should be greater than 0")
	}

	// Validate that output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output PDF file should exist")
	}

	// Validate that original HTML is unchanged
	originalContent, err := os.ReadFile(tempHTML)
	if err != nil {
		t.Fatalf("Failed to read original HTML: %v", err)
	}
	if !strings.Contains(string(originalContent), "#vb") {
		t.Error("Original HTML should still contain #vb (unchanged)")
	}
}

func TestChromeDPConverter_ValidateHTML(t *testing.T) {
	converter := NewChromeDPConverter()
	defer converter.Close()

	// Test with non-existent file
	err := converter.ValidateHTML("/non/existent/file.html")
	if err == nil {
		t.Error("ValidateHTML should return error for non-existent file")
	}

	// Test with existing file
	tempHTML, err := createTestHTML("<html><body>Test</body></html>")
	if err != nil {
		t.Fatalf("Failed to create test HTML: %v", err)
	}
	defer os.Remove(tempHTML)

	err = converter.ValidateHTML(tempHTML)
	if err != nil {
		t.Errorf("ValidateHTML should not return error for existing file: %v", err)
	}
}