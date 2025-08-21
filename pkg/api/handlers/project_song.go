package handlers

import (
	"net/http"
	"strconv"

	"github.com/bwl21/zupfmanager/pkg/api/models"
	"github.com/bwl21/zupfmanager/pkg/core"
	"github.com/gin-gonic/gin"
)

// ProjectSongHandler handles project-song relationship operations
type ProjectSongHandler struct {
	services *core.Services
}

// NewProjectSongHandler creates a new project-song handler
func NewProjectSongHandler(services *core.Services) *ProjectSongHandler {
	return &ProjectSongHandler{
		services: services,
	}
}

// AddSongToProject adds a song to a project
// @Summary Add song to project
// @Description Add an existing song to a project with optional difficulty and priority
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param songId path int true "Song ID"
// @Param request body models.AddSongToProjectRequest false "Add song request"
// @Success 201 {object} models.ProjectSongResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse "Song already in project"
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/projects/{id}/songs/{songId} [post]
func (h *ProjectSongHandler) AddSongToProject(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid project ID",
			Message: "Project ID must be a valid integer",
		})
		return
	}

	songID, err := strconv.Atoi(c.Param("songId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid song ID",
			Message: "Song ID must be a valid integer",
		})
		return
	}

	var req models.AddSongToProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Convert to core request
	coreReq := core.AddSongToProjectRequest{
		ProjectID:  projectID,
		SongID:     songID,
		Difficulty: req.Difficulty,
		Priority:   req.Priority,
		Comment:    req.Comment,
	}

	// Add song to project using core service
	projectSong, err := h.services.Project.AddSongToProject(c.Request.Context(), coreReq)
	if err != nil {
		switch err {
		case core.ErrProjectNotFound:
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Project not found",
				Message: "The specified project does not exist",
			})
		case core.ErrSongNotFound:
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Song not found",
				Message: "The specified song does not exist",
			})
		case core.ErrSongAlreadyInProject:
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Error:   "Song already in project",
				Message: "This song is already added to the project",
			})
		default:
			// Check if it's a validation error
			if validationErr, ok := err.(core.ValidationErrors); ok {
				details := make(map[string]string)
				for _, ve := range validationErr {
					details[ve.Field] = ve.Message
				}
				c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Error:   "validation failed",
					Message: err.Error(),
					Details: details,
				})
				return
			}
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Failed to add song to project",
				Message: err.Error(),
			})
		}
		return
	}

	// Convert to response
	response := models.ProjectSongResponse{
		ID:         projectSong.ID,
		ProjectID:  projectSong.ProjectID,
		SongID:     projectSong.SongID,
		Difficulty: projectSong.Difficulty,
		Priority:   projectSong.Priority,
		Comment:    projectSong.Comment,
	}

	// Add related entities if available
	if projectSong.Song != nil {
		response.Song = &models.SongResponse{
			ID:        projectSong.Song.ID,
			Title:     projectSong.Song.Title,
			Filename:  projectSong.Song.Filename,
			Genre:     projectSong.Song.Genre,
			Copyright: projectSong.Song.Copyright,
			Tocinfo:   projectSong.Song.Tocinfo,
		}
	}

	c.JSON(http.StatusCreated, response)
}

// RemoveSongFromProject removes a song from a project
// @Summary Remove song from project
// @Description Remove a song from a project
// @Tags projects
// @Param id path int true "Project ID"
// @Param songId path int true "Song ID"
// @Success 204 "Song removed successfully"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/projects/{id}/songs/{songId} [delete]
func (h *ProjectSongHandler) RemoveSongFromProject(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid project ID",
			Message: "Project ID must be a valid integer",
		})
		return
	}

	songID, err := strconv.Atoi(c.Param("songId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid song ID",
			Message: "Song ID must be a valid integer",
		})
		return
	}

	// Remove song from project using core service
	err = h.services.Project.RemoveSongFromProject(c.Request.Context(), projectID, songID)
	if err != nil {
		switch err {
		case core.ErrProjectSongNotFound:
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Project-song relationship not found",
				Message: "The song is not in the specified project",
			})
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Database error",
				Message: err.Error(),
			})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateProjectSong updates a project-song relationship
// @Summary Update project-song relationship
// @Description Update difficulty and priority of a song in a project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param songId path int true "Song ID"
// @Param request body models.UpdateProjectSongRequest true "Update request"
// @Success 200 {object} models.ProjectSongResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/projects/{id}/songs/{songId} [put]
func (h *ProjectSongHandler) UpdateProjectSong(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid project ID",
			Message: "Project ID must be a valid integer",
		})
		return
	}

	songID, err := strconv.Atoi(c.Param("songId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid song ID",
			Message: "Song ID must be a valid integer",
		})
		return
	}

	var req models.UpdateProjectSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Convert to core request
	coreReq := core.UpdateProjectSongRequest{
		ProjectID:  projectID,
		SongID:     songID,
		Difficulty: req.Difficulty,
		Priority:   req.Priority,
		Comment:    req.Comment,
	}

	// Update project-song using core service
	projectSong, err := h.services.Project.UpdateProjectSong(c.Request.Context(), coreReq)
	if err != nil {
		switch err {
		case core.ErrProjectSongNotFound:
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Project-song relationship not found",
				Message: "The song is not in the specified project",
			})
		default:
			// Check if it's a validation error
			if validationErr, ok := err.(core.ValidationErrors); ok {
				details := make(map[string]string)
				for _, ve := range validationErr {
					details[ve.Field] = ve.Message
				}
				c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Error:   "validation failed",
					Message: err.Error(),
					Details: details,
				})
				return
			}
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Failed to update project-song relationship",
				Message: err.Error(),
			})
		}
		return
	}

	// Convert to response
	response := models.ProjectSongResponse{
		ID:         projectSong.ID,
		ProjectID:  projectSong.ProjectID,
		SongID:     projectSong.SongID,
		Difficulty: projectSong.Difficulty,
		Priority:   projectSong.Priority,
		Comment:    projectSong.Comment,
	}

	// Add related entities if available
	if projectSong.Song != nil {
		response.Song = &models.SongResponse{
			ID:        projectSong.Song.ID,
			Title:     projectSong.Song.Title,
			Filename:  projectSong.Song.Filename,
			Genre:     projectSong.Song.Genre,
			Copyright: projectSong.Song.Copyright,
			Tocinfo:   projectSong.Song.Tocinfo,
		}
	}

	c.JSON(http.StatusOK, response)
}

// ListProjectSongs lists all songs in a project
// @Summary List project songs
// @Description Get all songs in a project with their relationships
// @Tags projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} models.ProjectSongsResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/projects/{id}/songs [get]
func (h *ProjectSongHandler) ListProjectSongs(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid project ID",
			Message: "Project ID must be a valid integer",
		})
		return
	}

	// List project songs using core service
	projectSongs, err := h.services.Project.ListProjectSongs(c.Request.Context(), projectID)
	if err != nil {
		switch err {
		case core.ErrProjectNotFound:
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Project not found",
				Message: "The specified project does not exist",
			})
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Database error",
				Message: err.Error(),
			})
		}
		return
	}

	// Convert to response
	responses := make([]models.ProjectSongResponse, len(projectSongs))
	for i, ps := range projectSongs {
		responses[i] = models.ProjectSongResponse{
			ID:         ps.ID,
			ProjectID:  ps.ProjectID,
			SongID:     ps.SongID,
			Difficulty: ps.Difficulty,
			Priority:   ps.Priority,
			Comment:    ps.Comment,
		}

		// Add song details if available
		if ps.Song != nil {
			songResponse := &models.SongResponse{
				ID:        ps.Song.ID,
				Title:     ps.Song.Title,
				Filename:  ps.Song.Filename,
				Genre:     ps.Song.Genre,
				Copyright: ps.Song.Copyright,
				Tocinfo:   ps.Song.Tocinfo,
			}

			// Add project associations if available
			if ps.Song.Projects != nil {
				projectRefs := make([]models.ProjectReference, len(ps.Song.Projects))
				for j, proj := range ps.Song.Projects {
					projectRefs[j] = models.ProjectReference{
						ID:        proj.ID,
						Title:     proj.Title,
						ShortName: proj.ShortName,
					}
				}
				songResponse.Projects = projectRefs
			}

			responses[i].Song = songResponse
		}
	}

	c.JSON(http.StatusOK, models.ProjectSongsResponse{
		ProjectSongs: responses,
		Total:        len(responses),
	})
}
