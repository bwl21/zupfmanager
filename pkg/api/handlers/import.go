package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/bwl21/zupfmanager/pkg/api/models"
	"github.com/bwl21/zupfmanager/pkg/core"
	"github.com/gin-gonic/gin"
)

// ImportHandler handles import-related API endpoints
type ImportHandler struct {
	services *core.Services
}

// NewImportHandler creates a new import handler
func NewImportHandler(services *core.Services) *ImportHandler {
	return &ImportHandler{
		services: services,
	}
}

// ImportFile imports a single ABC file
// @Summary Import a single ABC file
// @Description Import a single ABC notation file into the song database
// @Tags import
// @Accept json
// @Produce json
// @Param request body models.ImportFileRequest true "Import file request"
// @Success 200 {object} models.ImportResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/import/file [post]
func (h *ImportHandler) ImportFile(c *gin.Context) {
	var req models.ImportFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid request",
			Message: err.Error(),
		})
		return
	}

	// Validate file extension
	if filepath.Ext(req.FilePath) != ".abc" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid file type",
			Message: "only .abc files are supported",
		})
		return
	}

	// Import the file
	result := h.services.Import.ImportFile(c.Request.Context(), req.FilePath)
	
	// Convert core.ImportResult to API model
	apiResult := models.ImportResult{
		Filename: result.Filename,
		Title:    result.Title,
		Action:   result.Action,
		Changes:  result.Changes,
	}
	
	if result.Error != nil {
		apiResult.Error = result.Error.Error()
	}

	// Create summary
	summary := models.ImportSummary{
		Total: 1,
	}
	
	success := result.Error == nil
	if success {
		switch result.Action {
		case "created":
			summary.Created = 1
		case "updated":
			summary.Updated = 1
		case "unchanged":
			summary.Unchanged = 1
		}
	} else {
		summary.Errors = 1
	}

	response := models.ImportResponse{
		Success: success,
		Results: []models.ImportResult{apiResult},
		Summary: summary,
	}

	if success {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusInternalServerError, response)
	}
}

// ImportDirectory imports all ABC files from a directory
// @Summary Import ABC files from directory
// @Description Import all ABC notation files from a specified directory
// @Tags import
// @Accept json
// @Produce json
// @Param request body models.ImportDirectoryRequest true "Import directory request"
// @Success 200 {object} models.ImportResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/import/directory [post]
func (h *ImportHandler) ImportDirectory(c *gin.Context) {
	var req models.ImportDirectoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid request",
			Message: err.Error(),
		})
		return
	}

	// Import the directory
	results, err := h.services.Import.ImportDirectory(c.Request.Context(), req.DirectoryPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "import failed",
			Message: err.Error(),
		})
		return
	}

	// Convert core.ImportResult slice to API models
	apiResults := make([]models.ImportResult, len(results))
	summary := models.ImportSummary{
		Total: len(results),
	}

	for i, result := range results {
		apiResults[i] = models.ImportResult{
			Filename: result.Filename,
			Title:    result.Title,
			Action:   result.Action,
			Changes:  result.Changes,
		}
		
		if result.Error != nil {
			apiResults[i].Error = result.Error.Error()
			summary.Errors++
		} else {
			switch result.Action {
			case "created":
				summary.Created++
			case "updated":
				summary.Updated++
			case "unchanged":
				summary.Unchanged++
			}
		}
	}

	response := models.ImportResponse{
		Success: summary.Errors == 0,
		Results: apiResults,
		Summary: summary,
	}

	c.JSON(http.StatusOK, response)
}
