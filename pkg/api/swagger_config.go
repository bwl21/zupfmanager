package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// swaggerConfig serves a dynamic swagger config that adapts to the current host
func (s *Server) swaggerConfig(c *gin.Context) {
	host := c.GetHeader("Host")
	if host == "" {
		host = c.Request.Host
	}
	
	scheme := "http"
	
	// Detect if we're behind HTTPS proxy (like Gitpod)
	if c.GetHeader("X-Forwarded-Proto") == "https" ||
		c.GetHeader("X-Forwarded-Ssl") == "on" ||
		strings.Contains(host, "gitpod.io") ||
		c.Request.TLS != nil {
		scheme = "https"
	}

	config := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"title":       "Zupfmanager API",
			"description": "REST API for managing music projects and ABC notation files",
			"version":     "1.0",
		},
		"host":     host,
		"basePath": "/",
		"schemes":  []string{scheme},
		"paths":    getSwaggerPaths(),
		"definitions": getSwaggerDefinitions(),
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.JSON(http.StatusOK, config)
}

func getSwaggerPaths() map[string]interface{} {
	return map[string]interface{}{
		"/health": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"health"},
				"summary":     "Health check",
				"description": "Check if the API server is running",
				"produces":    []string{"application/json"},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "OK",
						"schema": map[string]interface{}{
							"type": "object",
						},
					},
				},
			},
		},
		"/api/v1/import/file": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"import"},
				"summary":     "Import a single ABC file",
				"description": "Import a single ABC notation file into the song database",
				"consumes":    []string{"application/json"},
				"produces":    []string{"application/json"},
				"parameters": []map[string]interface{}{
					{
						"in":       "body",
						"name":     "request",
						"required": true,
						"schema": map[string]interface{}{
							"$ref": "#/definitions/ImportFileRequest",
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Success",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/ImportResponse",
						},
					},
					"400": map[string]interface{}{
						"description": "Bad Request",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/ErrorResponse",
						},
					},
				},
			},
		},
		"/api/v1/import/directory": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"import"},
				"summary":     "Import ABC files from directory",
				"description": "Import all ABC notation files from a specified directory",
				"consumes":    []string{"application/json"},
				"produces":    []string{"application/json"},
				"parameters": []map[string]interface{}{
					{
						"in":       "body",
						"name":     "request",
						"required": true,
						"schema": map[string]interface{}{
							"$ref": "#/definitions/ImportDirectoryRequest",
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Success",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/ImportResponse",
						},
					},
					"400": map[string]interface{}{
						"description": "Bad Request",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/ErrorResponse",
						},
					},
				},
			},
		},
		"/api/v1/projects": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"projects"},
				"summary":     "List all projects",
				"description": "Get a list of all projects",
				"produces":    []string{"application/json"},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Success",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/ProjectListResponse",
						},
					},
				},
			},
			"post": map[string]interface{}{
				"tags":        []string{"projects"},
				"summary":     "Create a new project",
				"description": "Create a new project with title, short name and optional configuration",
				"consumes":    []string{"application/json"},
				"produces":    []string{"application/json"},
				"parameters": []map[string]interface{}{
					{
						"in":       "body",
						"name":     "request",
						"required": true,
						"schema": map[string]interface{}{
							"$ref": "#/definitions/CreateProjectRequest",
						},
					},
				},
				"responses": map[string]interface{}{
					"201": map[string]interface{}{
						"description": "Created",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/ProjectResponse",
						},
					},
					"400": map[string]interface{}{
						"description": "Bad Request",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/ErrorResponse",
						},
					},
				},
			},
		},
		"/api/v1/projects/{id}": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"projects"},
				"summary":     "Get project by ID",
				"description": "Get a specific project by its ID",
				"produces":    []string{"application/json"},
				"parameters": []map[string]interface{}{
					{
						"in":       "path",
						"name":     "id",
						"required": true,
						"type":     "integer",
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Success",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/ProjectResponse",
						},
					},
					"404": map[string]interface{}{
						"description": "Not Found",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/ErrorResponse",
						},
					},
				},
			},
			"put": map[string]interface{}{
				"tags":        []string{"projects"},
				"summary":     "Update project",
				"description": "Update an existing project",
				"consumes":    []string{"application/json"},
				"produces":    []string{"application/json"},
				"parameters": []map[string]interface{}{
					{
						"in":       "path",
						"name":     "id",
						"required": true,
						"type":     "integer",
					},
					{
						"in":       "body",
						"name":     "request",
						"required": true,
						"schema": map[string]interface{}{
							"$ref": "#/definitions/UpdateProjectRequest",
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Success",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/ProjectResponse",
						},
					},
					"400": map[string]interface{}{
						"description": "Bad Request",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/ErrorResponse",
						},
					},
				},
			},
			"delete": map[string]interface{}{
				"tags":        []string{"projects"},
				"summary":     "Delete project",
				"description": "Delete a project by ID",
				"parameters": []map[string]interface{}{
					{
						"in":       "path",
						"name":     "id",
						"required": true,
						"type":     "integer",
					},
				},
				"responses": map[string]interface{}{
					"204": map[string]interface{}{
						"description": "No Content",
					},
					"404": map[string]interface{}{
						"description": "Not Found",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/ErrorResponse",
						},
					},
				},
			},
		},
		"/api/v1/songs": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"songs"},
				"summary":     "List all songs",
				"description": "Get a list of all songs in the database",
				"produces":    []string{"application/json"},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Success",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/SongListResponse",
						},
					},
				},
			},
		},
		"/api/v1/songs/{id}": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"songs"},
				"summary":     "Get song by ID",
				"description": "Get a specific song by its ID",
				"produces":    []string{"application/json"},
				"parameters": []map[string]interface{}{
					{
						"in":       "path",
						"name":     "id",
						"required": true,
						"type":     "integer",
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Success",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/SongResponse",
						},
					},
					"404": map[string]interface{}{
						"description": "Not Found",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/ErrorResponse",
						},
					},
				},
			},
		},
		"/api/v1/songs/search": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"songs"},
				"summary":     "Search songs",
				"description": "Search for songs by query string",
				"produces":    []string{"application/json"},
				"parameters": []map[string]interface{}{
					{
						"in":       "query",
						"name":     "q",
						"required": true,
						"type":     "string",
					},
					{
						"in":       "query",
						"name":     "title",
						"required": false,
						"type":     "boolean",
						"default":  true,
					},
					{
						"in":       "query",
						"name":     "filename",
						"required": false,
						"type":     "boolean",
						"default":  false,
					},
					{
						"in":       "query",
						"name":     "genre",
						"required": false,
						"type":     "boolean",
						"default":  false,
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Success",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/SongListResponse",
						},
					},
					"400": map[string]interface{}{
						"description": "Bad Request",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/ErrorResponse",
						},
					},
				},
			},
		},
	}
}

func getSwaggerDefinitions() map[string]interface{} {
	return map[string]interface{}{
		"ImportFileRequest": map[string]interface{}{
			"type": "object",
			"required": []string{"file_path"},
			"properties": map[string]interface{}{
				"file_path": map[string]interface{}{
					"type":    "string",
					"example": "/path/to/song.abc",
				},
			},
		},
		"ImportDirectoryRequest": map[string]interface{}{
			"type": "object",
			"required": []string{"directory_path"},
			"properties": map[string]interface{}{
				"directory_path": map[string]interface{}{
					"type":    "string",
					"example": "/path/to/songs/",
				},
			},
		},
		"ImportResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"success": map[string]interface{}{
					"type":    "boolean",
					"example": true,
				},
				"results": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"$ref": "#/definitions/ImportResult",
					},
				},
				"summary": map[string]interface{}{
					"$ref": "#/definitions/ImportSummary",
				},
			},
		},
		"ImportResult": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"filename": map[string]interface{}{
					"type":    "string",
					"example": "song.abc",
				},
				"title": map[string]interface{}{
					"type":    "string",
					"example": "Amazing Grace",
				},
				"action": map[string]interface{}{
					"type":    "string",
					"example": "created",
				},
				"changes": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
				"error": map[string]interface{}{
					"type":    "string",
					"example": "file not found",
				},
			},
		},
		"ImportSummary": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"total": map[string]interface{}{
					"type":    "integer",
					"example": 10,
				},
				"created": map[string]interface{}{
					"type":    "integer",
					"example": 7,
				},
				"updated": map[string]interface{}{
					"type":    "integer",
					"example": 2,
				},
				"unchanged": map[string]interface{}{
					"type":    "integer",
					"example": 1,
				},
				"errors": map[string]interface{}{
					"type":    "integer",
					"example": 0,
				},
			},
		},
		"CreateProjectRequest": map[string]interface{}{
			"type": "object",
			"required": []string{"title", "short_name"},
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type":    "string",
					"example": "My Music Project",
				},
				"short_name": map[string]interface{}{
					"type":    "string",
					"example": "my-project",
				},
				"config_file": map[string]interface{}{
					"type":    "string",
					"example": "/path/to/config.json",
				},
				"default_config": map[string]interface{}{
					"type":    "boolean",
					"example": true,
				},
			},
		},
		"UpdateProjectRequest": map[string]interface{}{
			"type": "object",
			"required": []string{"title", "short_name"},
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type":    "string",
					"example": "Updated Project",
				},
				"short_name": map[string]interface{}{
					"type":    "string",
					"example": "updated-project",
				},
				"config_file": map[string]interface{}{
					"type":    "string",
					"example": "/path/to/config.json",
				},
				"default_config": map[string]interface{}{
					"type":    "boolean",
					"example": false,
				},
			},
		},
		"ProjectResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":    "integer",
					"example": 1,
				},
				"title": map[string]interface{}{
					"type":    "string",
					"example": "My Music Project",
				},
				"short_name": map[string]interface{}{
					"type":    "string",
					"example": "my-project",
				},
				"config": map[string]interface{}{
					"type": "object",
				},
			},
		},
		"ProjectListResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"projects": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"$ref": "#/definitions/ProjectResponse",
					},
				},
				"count": map[string]interface{}{
					"type":    "integer",
					"example": 5,
				},
			},
		},
		"SongResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":    "integer",
					"example": 1,
				},
				"title": map[string]interface{}{
					"type":    "string",
					"example": "Amazing Grace",
				},
				"filename": map[string]interface{}{
					"type":    "string",
					"example": "amazing_grace.abc",
				},
				"genre": map[string]interface{}{
					"type":    "string",
					"example": "Hymn",
				},
				"copyright": map[string]interface{}{
					"type":    "string",
					"example": "Public Domain",
				},
				"tocinfo": map[string]interface{}{
					"type":    "string",
					"example": "John Newton",
				},
			},
		},
		"SongListResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"songs": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"$ref": "#/definitions/SongResponse",
					},
				},
				"count": map[string]interface{}{
					"type":    "integer",
					"example": 10,
				},
			},
		},
		"ErrorResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"error": map[string]interface{}{
					"type":    "string",
					"example": "validation failed",
				},
				"message": map[string]interface{}{
					"type":    "string",
					"example": "Invalid input provided",
				},
				"details": map[string]interface{}{
					"type": "object",
					"additionalProperties": map[string]interface{}{
						"type": "string",
					},
				},
			},
		},
	}
}
