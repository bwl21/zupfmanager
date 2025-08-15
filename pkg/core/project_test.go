package core

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestProjectService_CreateProject(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Create default config file for testing
	defaultConfig := `{"test": "value"}`
	err := os.WriteFile("default-project-config.json", []byte(defaultConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create default config file: %v", err)
	}

	service, err := NewProjectService()
	if err != nil {
		t.Fatalf("Failed to create project service: %v", err)
	}
	defer service.Close()

	tests := []struct {
		name    string
		req     CreateProjectRequest
		wantErr bool
	}{
		{
			name: "create project with default config",
			req: CreateProjectRequest{
				Title:         "Test Project",
				ShortName:     "test-proj",
				DefaultConfig: true,
			},
			wantErr: false,
		},
		{
			name: "create project without config",
			req: CreateProjectRequest{
				Title:     "Simple Project",
				ShortName: "simple",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project, err := service.CreateProject(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if project.Title != tt.req.Title {
					t.Errorf("Expected title %s, got %s", tt.req.Title, project.Title)
				}
				if project.ShortName != tt.req.ShortName {
					t.Errorf("Expected short name %s, got %s", tt.req.ShortName, project.ShortName)
				}

				// Check if directory was created
				if _, err := os.Stat(tt.req.ShortName); os.IsNotExist(err) {
					t.Errorf("Project directory %s was not created", tt.req.ShortName)
				}
				if _, err := os.Stat(filepath.Join(tt.req.ShortName, "tpl")); os.IsNotExist(err) {
					t.Errorf("Template directory %s/tpl was not created", tt.req.ShortName)
				}
			}
		})
	}
}

func TestProjectService_ListProjects(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	service, err := NewProjectService()
	if err != nil {
		t.Fatalf("Failed to create project service: %v", err)
	}
	defer service.Close()

	// Create a test project first
	req := CreateProjectRequest{
		Title:     "List Test Project",
		ShortName: "list-test",
	}
	_, err = service.CreateProject(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	projects, err := service.ListProjects(context.Background())
	if err != nil {
		t.Fatalf("ListProjects() error = %v", err)
	}

	if len(projects) == 0 {
		t.Error("Expected at least one project, got none")
	}

	found := false
	for _, p := range projects {
		if p.Title == req.Title && p.ShortName == req.ShortName {
			found = true
			break
		}
	}
	if !found {
		t.Error("Created project not found in list")
	}
}

func TestProjectService_UpdateProject(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	service, err := NewProjectService()
	if err != nil {
		t.Fatalf("Failed to create project service: %v", err)
	}
	defer service.Close()

	// Create a test project first
	createReq := CreateProjectRequest{
		Title:     "Original Title",
		ShortName: "original",
	}
	project, err := service.CreateProject(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	// Update the project
	updateReq := UpdateProjectRequest{
		ID:        project.ID,
		Title:     "Updated Title",
		ShortName: "updated",
	}
	updatedProject, err := service.UpdateProject(context.Background(), updateReq)
	if err != nil {
		t.Fatalf("UpdateProject() error = %v", err)
	}

	if updatedProject.Title != updateReq.Title {
		t.Errorf("Expected title %s, got %s", updateReq.Title, updatedProject.Title)
	}
	if updatedProject.ShortName != updateReq.ShortName {
		t.Errorf("Expected short name %s, got %s", updateReq.ShortName, updatedProject.ShortName)
	}
}
