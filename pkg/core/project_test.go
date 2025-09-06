package core

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func setupProjectTest(t *testing.T) (*Services, func()) {
	// Create temporary directory for test
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	cleanup := func() {
		os.Chdir(oldWd)
	}
	os.Chdir(tempDir)

	// Create default config file for testing
	defaultConfig := `{"test": "value"}`
	err := os.WriteFile("default-project-config.json", []byte(defaultConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create default config file: %v", err)
	}

	services, err := NewServices()
	if err != nil {
		t.Fatalf("Failed to create services: %v", err)
	}

	return services, func() {
		services.Close()
		cleanup()
	}
}

func TestProjectService_Create(t *testing.T) {
	services, cleanup := setupProjectTest(t)
	defer cleanup()

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
		{
			name: "create project with invalid short name",
			req: CreateProjectRequest{
				Title:     "Invalid Project",
				ShortName: "invalid name with spaces",
			},
			wantErr: true,
		},
		{
			name: "create project with empty title",
			req: CreateProjectRequest{
				Title:     "",
				ShortName: "empty-title",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project, err := services.Project.Create(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
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

func TestProjectService_List(t *testing.T) {
	services, cleanup := setupProjectTest(t)
	defer cleanup()

	// Create a test project first
	req := CreateProjectRequest{
		Title:     "List Test Project",
		ShortName: "list-test",
	}
	_, err := services.Project.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	projects, err := services.Project.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
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

func TestProjectService_Update(t *testing.T) {
	services, cleanup := setupProjectTest(t)
	defer cleanup()

	// Create a test project first
	createReq := CreateProjectRequest{
		Title:     "Original Title",
		ShortName: "original",
	}
	project, err := services.Project.Create(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	// Update the project
	updateReq := UpdateProjectRequest{
		ID:        project.ID,
		Title:     "Updated Title",
		ShortName: "updated",
	}
	updatedProject, err := services.Project.Update(context.Background(), updateReq)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if updatedProject.Title != updateReq.Title {
		t.Errorf("Expected title %s, got %s", updateReq.Title, updatedProject.Title)
	}
	if updatedProject.ShortName != updateReq.ShortName {
		t.Errorf("Expected short name %s, got %s", updateReq.ShortName, updatedProject.ShortName)
	}
}

func TestProjectService_Get(t *testing.T) {
	services, cleanup := setupProjectTest(t)
	defer cleanup()

	// Create a test project first
	createReq := CreateProjectRequest{
		Title:     "Get Test Project",
		ShortName: "get-test",
	}
	project, err := services.Project.Create(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	// Get the project
	retrievedProject, err := services.Project.Get(context.Background(), project.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if retrievedProject.ID != project.ID {
		t.Errorf("Expected ID %d, got %d", project.ID, retrievedProject.ID)
	}
	if retrievedProject.Title != project.Title {
		t.Errorf("Expected title %s, got %s", project.Title, retrievedProject.Title)
	}
}

func TestProjectService_Delete(t *testing.T) {
	services, cleanup := setupProjectTest(t)
	defer cleanup()

	// Create a test project first
	createReq := CreateProjectRequest{
		Title:     "Delete Test Project",
		ShortName: "delete-test",
	}
	project, err := services.Project.Create(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	// Delete the project
	err = services.Project.Delete(context.Background(), project.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify it's deleted
	_, err = services.Project.Get(context.Background(), project.ID)
	if err == nil {
		t.Error("Expected error when getting deleted project, got none")
	}
}
