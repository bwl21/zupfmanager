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
		Config:        req.Config,
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
		ID:                   project.ID,
		Title:                project.Title,
		ShortName:            project.ShortName,
		Config:               project.Config,
		AbcFileDirPreference: project.AbcFileDirPreference,
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
			ID:                   project.ID,
			Title:                project.Title,
			ShortName:            project.ShortName,
			Config:               project.Config,
			AbcFileDirPreference: project.AbcFileDirPreference,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetDefaultConfig returns the default project configuration
// @Summary Get default project configuration
// @Description Get the default project configuration from default-project-config.json
// @Tags projects
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/projects/default-config [get]
func (h *ProjectHandler) GetDefaultConfig(c *gin.Context) {
	config, err := h.services.Config.LoadDefault()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "failed to load default configuration",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, config)
}

// BuildProject starts a project build operation
// @Summary Build project
// @Description Start building a project to generate ABC files, PDFs, and other outputs
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param request body models.BuildProjectRequest false "Build project request"
// @Success 202 {object} models.BuildResultResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/projects/{id}/build [post]
func (h *ProjectHandler) BuildProject(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid project ID",
			Message: "Project ID must be a valid integer",
		})
		return
	}

	var req models.BuildProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Convert to core request
	coreReq := core.BuildProjectRequest{
		ProjectID: projectID,
	}
	if req.OutputDir != nil {
		coreReq.OutputDir = *req.OutputDir
	}
	if req.AbcFileDir != nil {
		coreReq.AbcFileDir = *req.AbcFileDir
	}
	if req.PriorityThreshold != nil {
		coreReq.PriorityThreshold = *req.PriorityThreshold
	}
	if req.SampleID != nil {
		coreReq.SampleID = *req.SampleID
	}

	// Start build using core service
	buildResult, err := h.services.Project.BuildProject(c.Request.Context(), coreReq)
	if err != nil {
		switch err {
		case core.ErrProjectNotFound:
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Project not found",
				Message: "The specified project does not exist",
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
				Error:   "Failed to start build",
				Message: err.Error(),
			})
		}
		return
	}

	// Convert to response
	response := models.BuildResultResponse{
		BuildID:        buildResult.BuildID,
		ProjectID:      buildResult.ProjectID,
		Status:         buildResult.Status,
		OutputDir:      buildResult.OutputDir,
		GeneratedFiles: buildResult.GeneratedFiles,
		StartedAt:      buildResult.StartedAt,
		CompletedAt:    buildResult.CompletedAt,
		Error:          buildResult.Error,
	}

	c.JSON(http.StatusAccepted, response)
}

// GetBuildStatus returns the status of a build operation
// @Summary Get build status
// @Description Get the current status of a project build operation
// @Tags projects
// @Produce json
// @Param id path int true "Project ID"
// @Param buildId path string true "Build ID"
// @Success 200 {object} models.BuildStatusResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/projects/{id}/builds/{buildId}/status [get]
func (h *ProjectHandler) GetBuildStatus(c *gin.Context) {
	buildID := c.Param("buildId")
	if buildID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid build ID",
			Message: "Build ID is required",
		})
		return
	}

	// Get build status using core service
	buildStatus, err := h.services.Project.GetBuildStatus(c.Request.Context(), buildID)
	if err != nil {
		switch err {
		case core.ErrBuildNotFound:
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Build not found",
				Message: "The specified build does not exist",
			})
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Failed to get build status",
				Message: err.Error(),
			})
		}
		return
	}

	// Convert to response
	response := models.BuildStatusResponse{
		Status:      buildStatus.Status,
		Progress:    buildStatus.Progress,
		Message:     buildStatus.Message,
		StartedAt:   buildStatus.StartedAt,
		CompletedAt: buildStatus.CompletedAt,
		Error:       buildStatus.Error,
	}

	c.JSON(http.StatusOK, response)
}

// ListBuilds returns all builds for a project
// @Summary List project builds
// @Description Get all build operations for a project
// @Tags projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} models.BuildListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/projects/{id}/builds [get]
func (h *ProjectHandler) ListBuilds(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid project ID",
			Message: "Project ID must be a valid integer",
		})
		return
	}

	// List builds using core service
	builds, err := h.services.Project.ListBuilds(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to list builds",
			Message: err.Error(),
		})
		return
	}

	// Convert to response
	responses := make([]models.BuildResultResponse, len(builds))
	for i, build := range builds {
		responses[i] = models.BuildResultResponse{
			BuildID:        build.BuildID,
			ProjectID:      build.ProjectID,
			Status:         build.Status,
			OutputDir:      build.OutputDir,
			GeneratedFiles: build.GeneratedFiles,
			StartedAt:      build.StartedAt,
			CompletedAt:    build.CompletedAt,
			Error:          build.Error,
		}
	}

	c.JSON(http.StatusOK, models.BuildListResponse{
		Builds: responses,
		Total:  len(responses),
	})
}

// ClearBuildHistory clears all build history for a project
// @Summary Clear build history
// @Description Remove all build history for a project
// @Tags projects
// @Param id path int true "Project ID"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/projects/{id}/builds [delete]
func (h *ProjectHandler) ClearBuildHistory(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid project ID",
			Message: "Project ID must be a valid integer",
		})
		return
	}

	// Clear build history using core service
	err = h.services.Project.ClearBuildHistory(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to clear build history",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{
		Message: "Build history cleared successfully",
	})
}

// GetBuildDefaults returns default values for build configuration
// @Summary Get build defaults
// @Description Get default values for build configuration including abc_file_dir from last import
// @Tags projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} models.BuildDefaultsResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/projects/{id}/build/defaults [get]
func (h *ProjectHandler) GetBuildDefaults(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid project ID",
			Message: "Project ID must be a valid integer",
		})
		return
	}

	// Get the project to check if it exists and has abc_file_dir config
	project, err := h.services.Project.Get(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Project not found",
			Message: "The specified project does not exist",
		})
		return
	}

	defaults := models.BuildDefaultsResponse{
		OutputDir:         project.ShortName,
		AbcFileDir:        "",
		PriorityThreshold: 4,
		SampleID:          "",
	}

	// Priority order: abc_file_dir_preference > abc_file_dir (from config) > last import directory
	if project.AbcFileDirPreference != "" {
		defaults.AbcFileDir = project.AbcFileDirPreference
	} else if abcFileDir, ok := project.Config["abc_file_dir"].(string); ok && abcFileDir != "" {
		defaults.AbcFileDir = abcFileDir
	} else {
		// Try to get the last import directory
		lastImportDir, err := core.GetLastImportDir()
		if err == nil && lastImportDir != "" {
			defaults.AbcFileDir = lastImportDir
		}
	}

	c.JSON(http.StatusOK, defaults)
}

// UpdateAbcFileDirPreference updates the abc_file_dir preference for a project
// @Summary Update ABC file directory preference
// @Description Update the preferred ABC file directory for a project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param request body models.UpdateAbcFileDirRequest true "Update ABC file dir request"
// @Success 200 {object} models.ProjectResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/projects/{id}/abc-file-dir [put]
func (h *ProjectHandler) UpdateAbcFileDirPreference(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid project ID",
			Message: "Project ID must be a valid integer",
		})
		return
	}

	var req models.UpdateAbcFileDirRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid request",
			Message: err.Error(),
		})
		return
	}

	// Update the abc_file_dir_preference field directly
	_, err = h.services.DB().Project.UpdateOneID(projectID).SetAbcFileDirPreference(req.AbcFileDir).Save(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update project",
			Message: err.Error(),
		})
		return
	}

	// Return the updated project
	updatedProject, err := h.services.Project.Get(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to retrieve updated project",
			Message: err.Error(),
		})
		return
	}

	response := models.ProjectResponse{
		ID:                   updatedProject.ID,
		Title:                updatedProject.Title,
		ShortName:            updatedProject.ShortName,
		Config:               updatedProject.Config,
		AbcFileDirPreference: updatedProject.AbcFileDirPreference,
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
		ID:                   project.ID,
		Title:                project.Title,
		ShortName:            project.ShortName,
		Config:               project.Config,
		AbcFileDirPreference: project.AbcFileDirPreference,
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
		Config:        req.Config,
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
		ID:                   project.ID,
		Title:                project.Title,
		ShortName:            project.ShortName,
		Config:               project.Config,
		AbcFileDirPreference: project.AbcFileDirPreference,
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
