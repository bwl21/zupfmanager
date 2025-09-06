package core

import (
	"context"
	"os"
	"testing"
)

func setupSongTest(t *testing.T) (*Services, func()) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	cleanup := func() {
		os.Chdir(oldWd)
	}
	os.Chdir(tempDir)

	services, err := NewServices()
	if err != nil {
		t.Fatalf("Failed to create services: %v", err)
	}

	return services, func() {
		services.Close()
		cleanup()
	}
}

func TestSongService_List(t *testing.T) {
	services, cleanup := setupSongTest(t)
	defer cleanup()

	// Initially should have no songs
	songs, err := services.Song.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	initialCount := len(songs)

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
	result := services.Import.ImportFile(context.Background(), testFile)
	if result.Error != nil {
		t.Fatalf("Failed to import test song: %v", result.Error)
	}

	// Now list songs again
	songs, err = services.Song.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
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

func TestSongService_Get(t *testing.T) {
	services, cleanup := setupSongTest(t)
	defer cleanup()

	testContent := `T:Test Get Song
Z:genre Rock
Z:copyright Test Copyright
K:C`

	testFile := "test-get.abc"
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result := services.Import.ImportFile(context.Background(), testFile)
	if result.Error != nil {
		t.Fatalf("Failed to import test song: %v", result.Error)
	}

	// Get all songs to find the ID of our test song
	songs, err := services.Song.List(context.Background())
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

	// Test Get
	song, err := services.Song.Get(context.Background(), testSongID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
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

func TestSongService_Search(t *testing.T) {
	services, cleanup := setupSongTest(t)
	defer cleanup()

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
		err := os.WriteFile(song.filename, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", song.filename, err)
		}

		result := services.Import.ImportFile(context.Background(), song.filename)
		if result.Error != nil {
			t.Fatalf("Failed to import test song %s: %v", song.filename, result.Error)
		}
	}

	// Test search
	results, err := services.Song.Search(context.Background(), "Grace")
	if err != nil {
		t.Fatalf("Search() error = %v", err)
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

func TestSongService_SearchAdvanced(t *testing.T) {
	services, cleanup := setupSongTest(t)
	defer cleanup()

	// Create test songs with different metadata
	testSongs := []struct {
		filename string
		content  string
	}{
		{"title-test.abc", "T:Grace in Title\nZ:genre Folk\nK:C"},
		{"genre-test.abc", "T:Different Song\nZ:genre Grace Genre\nK:C"},
		{"filename-grace.abc", "T:Another Song\nZ:genre Rock\nK:C"},
	}

	for _, song := range testSongs {
		err := os.WriteFile(song.filename, []byte(song.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", song.filename, err)
		}

		result := services.Import.ImportFile(context.Background(), song.filename)
		if result.Error != nil {
			t.Fatalf("Failed to import test song %s: %v", song.filename, result.Error)
		}
	}

	// Test search only in titles
	results, err := services.Song.SearchAdvanced(context.Background(), "Grace", SearchOptions{
		SearchTitle:    true,
		SearchFilename: false,
		SearchGenre:    false,
	})
	if err != nil {
		t.Fatalf("SearchAdvanced() error = %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result for title search, got %d", len(results))
	}
	if len(results) > 0 && results[0].Title != "Grace in Title" {
		t.Errorf("Expected 'Grace in Title', got %s", results[0].Title)
	}

	// Test search only in filenames
	results, err = services.Song.SearchAdvanced(context.Background(), "grace", SearchOptions{
		SearchTitle:    false,
		SearchFilename: true,
		SearchGenre:    false,
	})
	if err != nil {
		t.Fatalf("SearchAdvanced() error = %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result for filename search, got %d", len(results))
	}
	if len(results) > 0 && results[0].Title != "Another Song" {
		t.Errorf("Expected 'Another Song', got %s", results[0].Title)
	}
}
