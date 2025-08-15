package core

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestImportService_parseABCMetadata(t *testing.T) {
	service := &ImportService{}

	tests := []struct {
		name     string
		content  string
		expected ABCMetadata
	}{
		{
			name: "basic metadata",
			content: `T:Test Song
Z:genre Folk
Z:copyright Public Domain
C:M: Test Info`,
			expected: ABCMetadata{
				Title:     "Test Song",
				Genre:     "Folk",
				Copyright: "Public Domain",
				Tocinfo:   "Test Info",
			},
		},
		{
			name: "minimal metadata",
			content: `T:Simple Song
K:C`,
			expected: ABCMetadata{
				Title: "Simple Song",
			},
		},
		{
			name: "various tocinfo patterns",
			content: `T:Pattern Test
C:M+T: Mixed Info`,
			expected: ABCMetadata{
				Title:   "Pattern Test",
				Tocinfo: "Mixed Info",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.parseABCMetadata([]byte(tt.content))
			if result.Title != tt.expected.Title {
				t.Errorf("Expected title %s, got %s", tt.expected.Title, result.Title)
			}
			if result.Genre != tt.expected.Genre {
				t.Errorf("Expected genre %s, got %s", tt.expected.Genre, result.Genre)
			}
			if result.Copyright != tt.expected.Copyright {
				t.Errorf("Expected copyright %s, got %s", tt.expected.Copyright, result.Copyright)
			}
			if result.Tocinfo != tt.expected.Tocinfo {
				t.Errorf("Expected tocinfo %s, got %s", tt.expected.Tocinfo, result.Tocinfo)
			}
		})
	}
}

func TestImportService_ImportFile(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	service, err := NewImportService()
	if err != nil {
		t.Fatalf("Failed to create import service: %v", err)
	}
	defer service.Close()

	// Create test ABC file
	testContent := `T:Test Import Song
Z:genre Test Genre
Z:copyright Test Copyright
C:M: Test Tocinfo
K:C
CDEF|`

	testFile := filepath.Join(tempDir, "test.abc")
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test first import (create)
	result := service.ImportFile(context.Background(), testFile)
	if result.Error != nil {
		t.Fatalf("ImportFile() error = %v", result.Error)
	}
	if result.Action != "created" {
		t.Errorf("Expected action 'created', got %s", result.Action)
	}
	if result.Title != "Test Import Song" {
		t.Errorf("Expected title 'Test Import Song', got %s", result.Title)
	}

	// Test second import (unchanged)
	result = service.ImportFile(context.Background(), testFile)
	if result.Error != nil {
		t.Fatalf("ImportFile() error = %v", result.Error)
	}
	if result.Action != "unchanged" {
		t.Errorf("Expected action 'unchanged', got %s", result.Action)
	}

	// Test import with changes (update)
	updatedContent := `T:Updated Import Song
Z:genre Updated Genre
Z:copyright Test Copyright
C:M: Test Tocinfo
K:C
CDEF|`

	err = os.WriteFile(testFile, []byte(updatedContent), 0644)
	if err != nil {
		t.Fatalf("Failed to update test file: %v", err)
	}

	result = service.ImportFile(context.Background(), testFile)
	if result.Error != nil {
		t.Fatalf("ImportFile() error = %v", result.Error)
	}
	if result.Action != "updated" {
		t.Errorf("Expected action 'updated', got %s", result.Action)
	}
	if len(result.Changes) == 0 {
		t.Error("Expected changes to be recorded")
	}
}

func TestImportService_ImportDirectory(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	service, err := NewImportService()
	if err != nil {
		t.Fatalf("Failed to create import service: %v", err)
	}
	defer service.Close()

	// Create test directory with ABC files
	testDir := filepath.Join(tempDir, "testdir")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create multiple test files
	files := []struct {
		name    string
		content string
	}{
		{
			name: "song1.abc",
			content: `T:Song One
K:C`,
		},
		{
			name: "song2.abc",
			content: `T:Song Two
Z:genre Folk
K:G`,
		},
	}

	for _, file := range files {
		err = os.WriteFile(filepath.Join(testDir, file.name), []byte(file.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file.name, err)
		}
	}

	// Import directory
	results, err := service.ImportDirectory(context.Background(), testDir)
	if err != nil {
		t.Fatalf("ImportDirectory() error = %v", err)
	}

	if len(results) != len(files) {
		t.Errorf("Expected %d results, got %d", len(files), len(results))
	}

	for _, result := range results {
		if result.Error != nil {
			t.Errorf("Import result for %s has error: %v", result.Filename, result.Error)
		}
		if result.Action != "created" {
			t.Errorf("Expected action 'created' for %s, got %s", result.Filename, result.Action)
		}
	}
}
