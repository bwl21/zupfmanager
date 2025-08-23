package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/bwl21/zupfmanager/pkg/api/models"
	"github.com/bwl21/zupfmanager/pkg/core"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportHandler_GetLastImportPath(t *testing.T) {
	// Setup
	services, err := core.NewServices()
	require.NoError(t, err)
	defer services.Close()

	handler := NewImportHandler(services)
	
	// Set a test path
	ctx := context.Background()
	testPath := "/test/import/path"
	err = services.Settings.Set(ctx, "last_import_path", testPath)
	require.NoError(t, err)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/import/last-path", handler.GetLastImportPath)

	// Make request
	req, _ := http.NewRequest("GET", "/api/v1/import/last-path", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, testPath, response["path"])
}

func TestImportHandler_DirectoryImportSavesPath(t *testing.T) {
	// Setup
	services, err := core.NewServices()
	require.NoError(t, err)
	defer services.Close()

	handler := NewImportHandler(services)
	
	// Create test directory and file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.abc")
	abcContent := `X:1
T:Test Song
M:4/4
K:C
C D E F | G A B c |]`
	
	err = os.WriteFile(testFile, []byte(abcContent), 0644)
	require.NoError(t, err)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/import/directory", handler.ImportDirectory)
	router.GET("/api/v1/import/last-path", handler.GetLastImportPath)

	// Import directory
	reqBody := models.ImportDirectoryRequest{
		DirectoryPath: tempDir,
	}
	jsonBody, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/api/v1/import/directory", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check that path was saved
	req2, _ := http.NewRequest("GET", "/api/v1/import/last-path", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	
	var response map[string]string
	err = json.Unmarshal(w2.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, tempDir, response["path"])
}
