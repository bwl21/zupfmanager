package core

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "no validation errors"
	}
	
	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// HasErrors returns true if there are validation errors
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

var (
	// alphanumRegex matches alphanumeric characters, hyphens, and underscores
	alphanumRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// ValidateCreateProjectRequest validates a create project request
func ValidateCreateProjectRequest(req CreateProjectRequest) error {
	var errors ValidationErrors

	// Validate Title
	if strings.TrimSpace(req.Title) == "" {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "title is required and cannot be empty",
		})
	}

	// Validate ShortName
	shortName := strings.TrimSpace(req.ShortName)
	if shortName == "" {
		errors = append(errors, ValidationError{
			Field:   "short_name",
			Message: "short_name is required and cannot be empty",
		})
	} else if len(shortName) > 50 {
		errors = append(errors, ValidationError{
			Field:   "short_name",
			Message: "short_name cannot be longer than 50 characters",
		})
	} else if !alphanumRegex.MatchString(shortName) {
		errors = append(errors, ValidationError{
			Field:   "short_name",
			Message: "short_name can only contain alphanumeric characters, hyphens, and underscores",
		})
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateBuildProjectRequest validates a build project request
func ValidateBuildProjectRequest(req BuildProjectRequest) error {
	var errors ValidationErrors

	// Validate ProjectID
	if req.ProjectID <= 0 {
		errors = append(errors, ValidationError{
			Field:   "project_id",
			Message: "project_id must be a positive integer",
		})
	}

	// Validate PriorityThreshold if provided
	if req.PriorityThreshold != 0 && (req.PriorityThreshold < 1 || req.PriorityThreshold > 4) {
		errors = append(errors, ValidationError{
			Field:   "priority_threshold",
			Message: "priority_threshold must be between 1 and 4",
		})
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateAddSongToProjectRequest validates an add song to project request
func ValidateAddSongToProjectRequest(req AddSongToProjectRequest) error {
	var errors ValidationErrors

	// Validate ProjectID
	if req.ProjectID <= 0 {
		errors = append(errors, ValidationError{
			Field:   "project_id",
			Message: "project_id must be a positive integer",
		})
	}

	// Validate SongID
	if req.SongID <= 0 {
		errors = append(errors, ValidationError{
			Field:   "song_id",
			Message: "song_id must be a positive integer",
		})
	}

	// Validate Difficulty if provided
	if req.Difficulty != nil {
		validDifficulties := []string{"easy", "medium", "hard", "expert"}
		valid := false
		for _, d := range validDifficulties {
			if *req.Difficulty == d {
				valid = true
				break
			}
		}
		if !valid {
			errors = append(errors, ValidationError{
				Field:   "difficulty",
				Message: "difficulty must be one of: easy, medium, hard, expert",
			})
		}
	}

	// Validate Priority if provided
	if req.Priority != nil {
		if *req.Priority < 1 || *req.Priority > 4 {
			errors = append(errors, ValidationError{
				Field:   "priority",
				Message: "priority must be between 1 and 4",
			})
		}
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateUpdateProjectSongRequest validates an update project song request
func ValidateUpdateProjectSongRequest(req UpdateProjectSongRequest) error {
	var errors ValidationErrors

	// Validate ProjectID
	if req.ProjectID <= 0 {
		errors = append(errors, ValidationError{
			Field:   "project_id",
			Message: "project_id must be a positive integer",
		})
	}

	// Validate SongID
	if req.SongID <= 0 {
		errors = append(errors, ValidationError{
			Field:   "song_id",
			Message: "song_id must be a positive integer",
		})
	}

	// Validate Difficulty if provided
	if req.Difficulty != nil {
		validDifficulties := []string{"easy", "medium", "hard", "expert"}
		valid := false
		for _, d := range validDifficulties {
			if *req.Difficulty == d {
				valid = true
				break
			}
		}
		if !valid {
			errors = append(errors, ValidationError{
				Field:   "difficulty",
				Message: "difficulty must be one of: easy, medium, hard, expert",
			})
		}
	}

	// Validate Priority if provided
	if req.Priority != nil {
		if *req.Priority < 1 || *req.Priority > 4 {
			errors = append(errors, ValidationError{
				Field:   "priority",
				Message: "priority must be between 1 and 4",
			})
		}
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateUpdateProjectRequest validates an update project request
func ValidateUpdateProjectRequest(req UpdateProjectRequest) error {
	var errors ValidationErrors

	// Validate ID
	if req.ID <= 0 {
		errors = append(errors, ValidationError{
			Field:   "id",
			Message: "id must be a positive integer",
		})
	}

	// Validate Title
	if strings.TrimSpace(req.Title) == "" {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "title is required and cannot be empty",
		})
	}

	// Validate ShortName
	shortName := strings.TrimSpace(req.ShortName)
	if shortName == "" {
		errors = append(errors, ValidationError{
			Field:   "short_name",
			Message: "short_name is required and cannot be empty",
		})
	} else if len(shortName) > 50 {
		errors = append(errors, ValidationError{
			Field:   "short_name",
			Message: "short_name cannot be longer than 50 characters",
		})
	} else if !alphanumRegex.MatchString(shortName) {
		errors = append(errors, ValidationError{
			Field:   "short_name",
			Message: "short_name can only contain alphanumeric characters, hyphens, and underscores",
		})
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}
