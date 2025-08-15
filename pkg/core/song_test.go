package core

import (
	"context"
	"os"
	"testing"
)

func TestSongService_ListSongs(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	service, err := NewSongService()
	if err != nil {
		t.Fatalf("Failed to create song service: %v", err)
	}
	defer service.Close()

	// Initially should have no songs
	songs, err := service.ListSongs(context.Background())
	if err != nil {
		t.Fatalf("ListSongs() error = %v", err)
	}

	initialCount := len(songs)

	// Create a song via import service to test listing
	importService, err := NewImportService()
	if err != nil {
		t.Fatalf("Failed to create import service: %v", err)
	}
	defer importService.Close()

	// Create test ABC file
	testContent := `T:Test List Song
Z:genre Test
K:C`

	testFile := "test-list.abc"
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Import the file to create a song
	result := importService.ImportFile(context.Background(), testFile)
	if result.Error != nil {
		t.Fatalf("Failed to import test song: %v", result.Error)
	}

	// Now list songs again
	songs, err = service.ListSongs(context.Background())
	if err != nil {
		t.Fatalf("ListSongs() error = %v", err)
	}

	if len(songs) != initialCount+1 {
		t.Errorf("Expected %d songs, got %d", initialCount+1, len(songs))
	}

	// Check if our test song is in the list
	found := false
	for _, song := range songs {
		if song.Title == "Test List Song" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Test song not found in list")
	}
}

func TestSongService_GetSong(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	service, err := NewSongService()
	if err != nil {
		t.Fatalf("Failed to create song service: %v", err)
	}
	defer service.Close()

	// Create a song via import service
	importService, err := NewImportService()
	if err != nil {
		t.Fatalf("Failed to create import service: %v", err)
	}
	defer importService.Close()

	testContent := `T:Test Get Song
Z:genre Rock
Z:copyright Test Copyright
K:C`

	testFile := "test-get.abc"
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result := importService.ImportFile(context.Background(), testFile)
	if result.Error != nil {
		t.Fatalf("Failed to import test song: %v", result.Error)
	}

	// Get all songs to find the ID of our test song
	songs, err := service.ListSongs(context.Background())
	if err != nil {
		t.Fatalf("Failed to list songs: %v", err)
	}

	var testSongID int
	for _, song := range songs {
		if song.Title == "Test Get Song" {
			testSongID = song.ID
			break
		}
	}

	if testSongID == 0 {
		t.Fatal("Test song not found")
	}

	// Test GetSong
	song, err := service.GetSong(context.Background(), testSongID)
	if err != nil {
		t.Fatalf("GetSong() error = %v", err)
	}

	if song.Title != "Test Get Song" {
		t.Errorf("Expected title 'Test Get Song', got %s", song.Title)
	}
	if song.Genre != "Rock" {
		t.Errorf("Expected genre 'Rock', got %s", song.Genre)
	}
	if song.Copyright != "Test Copyright" {
		t.Errorf("Expected copyright 'Test Copyright', got %s", song.Copyright)
	}
}

func TestSongService_SearchSongs(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	service, err := NewSongService()
	if err != nil {
		t.Fatalf("Failed to create song service: %v", err)
	}
	defer service.Close()

	// Create multiple songs via import service
	importService, err := NewImportService()
	if err != nil {
		t.Fatalf("Failed to create import service: %v", err)
	}
	defer importService.Close()

	testSongs := []struct {
		filename string
		title    string
	}{
		{"search1.abc", "Amazing Grace"},
		{"search2.abc", "Grace Notes"},
		{"search3.abc", "Simple Song"},
	}

	for _, song := range testSongs {
		content := "T:" + song.title + "\nK:C"
		err = os.WriteFile(song.filename, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", song.filename, err)
		}

		result := importService.ImportFile(context.Background(), song.filename)
		if result.Error != nil {
			t.Fatalf("Failed to import test song %s: %v", song.filename, result.Error)
		}
	}

	// Test search
	results, err := service.SearchSongs(context.Background(), "Grace")
	if err != nil {
		t.Fatalf("SearchSongs() error = %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'Grace', got %d", len(results))
	}

	// Check that both Grace songs are found
	foundTitles := make(map[string]bool)
	for _, song := range results {
		foundTitles[song.Title] = true
	}

	if !foundTitles["Amazing Grace"] {
		t.Error("'Amazing Grace' not found in search results")
	}
	if !foundTitles["Grace Notes"] {
		t.Error("'Grace Notes' not found in search results")
	}
	if foundTitles["Simple Song"] {
		t.Error("'Simple Song' should not be in search results")
	}
}
