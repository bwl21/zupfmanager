package handlers

import (
	"net/http"
	"strconv"

	"github.com/bwl21/zupfmanager/pkg/api/models"
	"github.com/bwl21/zupfmanager/pkg/core"
	"github.com/gin-gonic/gin"
)

// ProjectHandler handles project-related API endpoints
type ProjectHandler struct {
	services *core.Services
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(services *core.Services) *ProjectHandler {
	return &ProjectHandler{
		services: services,
	}
}

// CreateProject creates a new project
// @Summary Create a new project
// @Description Create a new project with title, short name and optional configuration
// @Tags projects
// @Accept json
// @Produce json
// @Param request body models.CreateProjectRequest true "Create project request"
// @Success 201 {object} models.ProjectResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid request",
			Message: err.Error(),
		})
		return
	}

	// Convert to core request
	coreReq := core.CreateProjectRequest{
		Title:         req.Title,
		ShortName:     req.ShortName,
		ConfigFile:    req.ConfigFile,
		DefaultConfig: req.DefaultConfig,
	}

	project, err := h.services.Project.Create(c.Request.Context(), coreReq)
	if err != nil {
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
			Error:   "failed to create project",
			Message: err.Error(),
		})
		return
	}

	response := models.ProjectResponse{
		ID:        project.ID,
		Title:     project.Title,
		ShortName: project.ShortName,
		Config:    project.Config,
	}

	c.JSON(http.StatusCreated, response)
}

// ListProjects lists all projects
// @Summary List all projects
// @Description Get a list of all projects
// @Tags projects
// @Produce json
// @Success 200 {object} models.ProjectListResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/projects [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	// Debug: Check if handler is being called
	if h == nil || h.services == nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "handler not initialized",
			Message: "project handler or services is nil",
		})
		return
	}

	projects, err := h.services.Project.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "failed to list projects",
			Message: err.Error(),
		})
		return
	}

	response := models.ProjectListResponse{
		Projects: make([]models.ProjectResponse, len(projects)),
		Count:    len(projects),
	}

	for i, project := range projects {
		response.Projects[i] = models.ProjectResponse{
			ID:        project.ID,
			Title:     project.Title,
			ShortName: project.ShortName,
			Config:    project.Config,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetProject gets a project by ID
// @Summary Get project by ID
// @Description Get a specific project by its ID
// @Tags projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} models.ProjectResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid project ID",
			Message: "project ID must be a number",
		})
		return
	}

	project, err := h.services.Project.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "project not found",
			Message: err.Error(),
		})
		return
	}

	response := models.ProjectResponse{
		ID:        project.ID,
		Title:     project.Title,
		ShortName: project.ShortName,
		Config:    project.Config,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProject updates a project
// @Summary Update project
// @Description Update an existing project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param request body models.UpdateProjectRequest true "Update project request"
// @Success 200 {object} models.ProjectResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid project ID",
			Message: "project ID must be a number",
		})
		return
	}

	var req models.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid request",
			Message: err.Error(),
		})
		return
	}

	// Convert to core request
	coreReq := core.UpdateProjectRequest{
		ID:            id,
		Title:         req.Title,
		ShortName:     req.ShortName,
		ConfigFile:    req.ConfigFile,
		DefaultConfig: req.DefaultConfig,
	}

	project, err := h.services.Project.Update(c.Request.Context(), coreReq)
	if err != nil {
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
			Error:   "failed to update project",
			Message: err.Error(),
		})
		return
	}

	response := models.ProjectResponse{
		ID:        project.ID,
		Title:     project.Title,
		ShortName: project.ShortName,
		Config:    project.Config,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteProject deletes a project
// @Summary Delete project
// @Description Delete a project by ID
// @Tags projects
// @Param id path int true "Project ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid project ID",
			Message: "project ID must be a number",
		})
		return
	}

	err = h.services.Project.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "failed to delete project",
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
