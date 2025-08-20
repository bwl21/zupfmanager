package handlers

import (
	"net/http"
	"strconv"

	"github.com/bwl21/zupfmanager/pkg/api/models"
	"github.com/bwl21/zupfmanager/pkg/core"
	"github.com/gin-gonic/gin"
)

// SongHandler handles song-related API endpoints
type SongHandler struct {
	services *core.Services
}

// NewSongHandler creates a new song handler
func NewSongHandler(services *core.Services) *SongHandler {
	return &SongHandler{
		services: services,
	}
}

// ListSongs lists all songs
// @Summary List all songs
// @Description Get a list of all songs in the database
// @Tags songs
// @Produce json
// @Success 200 {object} models.SongListResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/songs [get]
func (h *SongHandler) ListSongs(c *gin.Context) {
	songs, err := h.services.Song.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "failed to list songs",
			Message: err.Error(),
		})
		return
	}

	response := models.SongListResponse{
		Songs: make([]models.SongResponse, len(songs)),
		Count: len(songs),
	}

	for i, song := range songs {
		response.Songs[i] = models.SongResponse{
			ID:        song.ID,
			Title:     song.Title,
			Filename:  song.Filename,
			Genre:     song.Genre,
			Copyright: song.Copyright,
			Tocinfo:   song.Tocinfo,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GeneratePreview generates preview PDFs for a song
// @Summary Generate preview PDFs
// @Description Generate preview PDFs for a specific song
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param request body models.GeneratePreviewRequest true "Preview generation request"
// @Success 200 {object} models.GeneratePreviewResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/songs/{id}/generate-preview [post]
func (h *SongHandler) GeneratePreview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid song ID",
			Message: "song ID must be a number",
		})
		return
	}

	var request models.GeneratePreviewRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Set the song ID from the URL parameter
	coreRequest := core.GeneratePreviewRequest{
		SongID:     id,
		AbcFileDir: request.AbcFileDir,
		Config:     request.Config,
	}

	result, err := h.services.Song.GeneratePreview(c.Request.Context(), coreRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "failed to generate preview",
			Message: err.Error(),
		})
		return
	}

	response := models.GeneratePreviewResponse{
		PDFFiles:   result.PDFFiles,
		PreviewDir: result.PreviewDir,
	}

	c.JSON(http.StatusOK, response)
}

// ListPreviewPDFs lists available preview PDFs for a song
// @Summary List preview PDFs
// @Description Get a list of available preview PDFs for a song
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} models.PreviewPDFListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/songs/{id}/preview-pdfs [get]
func (h *SongHandler) ListPreviewPDFs(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid song ID",
			Message: "song ID must be a number",
		})
		return
	}

	pdfs, err := h.services.Song.ListPreviewPDFs(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "failed to list preview PDFs",
			Message: err.Error(),
		})
		return
	}

	response := models.PreviewPDFListResponse{
		PDFs:  make([]models.PreviewPDFResponse, len(pdfs)),
		Count: len(pdfs),
	}

	for i, pdf := range pdfs {
		response.PDFs[i] = models.PreviewPDFResponse{
			Filename:  pdf.Filename,
			Size:      pdf.Size,
			CreatedAt: pdf.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetPreviewPDF serves a preview PDF file
// @Summary Get preview PDF
// @Description Download or view a specific preview PDF file
// @Tags songs
// @Produce application/pdf
// @Param id path int true "Song ID"
// @Param filename path string true "PDF filename"
// @Param abc_file_dir query string true "ABC file directory"
// @Success 200 {file} application/pdf
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/songs/{id}/preview-pdf/{filename} [get]
func (h *SongHandler) GetPreviewPDF(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid song ID",
			Message: "song ID must be a number",
		})
		return
	}

	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing filename",
			Message: "filename parameter is required",
		})
		return
	}

	abcFileDir := c.Query("abc_file_dir")
	if abcFileDir == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing abc_file_dir",
			Message: "abc_file_dir query parameter is required",
		})
		return
	}

	filePath, err := h.services.Song.GetPreviewPDFFromDir(c.Request.Context(), id, filename, abcFileDir)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "PDF not found",
			Message: err.Error(),
		})
		return
	}

	// Serve the PDF file
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "inline; filename=\""+filename+"\"")
	c.File(filePath)
}

// CleanupPreviewPDFs removes all preview PDFs for a song
// @Summary Cleanup preview PDFs
// @Description Remove all preview PDFs for a specific song
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/songs/{id}/preview-pdfs [delete]
func (h *SongHandler) CleanupPreviewPDFs(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid song ID",
			Message: "song ID must be a number",
		})
		return
	}

	err = h.services.Song.CleanupPreviewPDFs(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "failed to cleanup preview PDFs",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.MessageResponse{
		Message: "Preview PDFs deleted successfully",
	})
}

// GetSong gets a song by ID
// @Summary Get song by ID
// @Description Get a specific song by its ID
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} models.SongResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/songs/{id} [get]
func (h *SongHandler) GetSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid song ID",
			Message: "song ID must be a number",
		})
		return
	}

	song, err := h.services.Song.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "song not found",
			Message: err.Error(),
		})
		return
	}

	response := models.SongResponse{
		ID:        song.ID,
		Title:     song.Title,
		Filename:  song.Filename,
		Genre:     song.Genre,
		Copyright: song.Copyright,
		Tocinfo:   song.Tocinfo,
	}

	c.JSON(http.StatusOK, response)
}

// SearchSongs searches for songs
// @Summary Search songs
// @Description Search for songs by query string
// @Tags songs
// @Produce json
// @Param q query string true "Search query"
// @Param title query bool false "Search in title" default(true)
// @Param filename query bool false "Search in filename" default(false)
// @Param genre query bool false "Search in genre" default(false)
// @Success 200 {object} models.SongListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/songs/search [get]
func (h *SongHandler) SearchSongs(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing search query",
			Message: "query parameter 'q' is required",
		})
		return
	}

	// Parse search options
	searchTitle := c.DefaultQuery("title", "true") == "true"
	searchFilename := c.DefaultQuery("filename", "false") == "true"
	searchGenre := c.DefaultQuery("genre", "false") == "true"

	// If no specific fields are selected, default to title search
	if !searchTitle && !searchFilename && !searchGenre {
		searchTitle = true
	}

	options := core.SearchOptions{
		SearchTitle:    searchTitle,
		SearchFilename: searchFilename,
		SearchGenre:    searchGenre,
	}

	var songs []*core.Song
	var err error

	// Use advanced search if specific options are set
	if searchFilename || searchGenre || !searchTitle {
		songs, err = h.services.Song.SearchAdvanced(c.Request.Context(), query, options)
	} else {
		songs, err = h.services.Song.Search(c.Request.Context(), query)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "search failed",
			Message: err.Error(),
		})
		return
	}

	response := models.SongListResponse{
		Songs: make([]models.SongResponse, len(songs)),
		Count: len(songs),
	}

	for i, song := range songs {
		response.Songs[i] = models.SongResponse{
			ID:        song.ID,
			Title:     song.Title,
			Filename:  song.Filename,
			Genre:     song.Genre,
			Copyright: song.Copyright,
			Tocinfo:   song.Tocinfo,
		}
	}

	c.JSON(http.StatusOK, response)
}
