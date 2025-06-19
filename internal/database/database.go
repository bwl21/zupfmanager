package database

import (
	"context"
	"log/slog"

	"github.com/bwl21/zupfmanager/internal/ent"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Client is the database client
type Client struct {
	*ent.Client
}

// New creates a new database client
func New() (*Client, error) {
	// Create SQLite driver
	driver, err := sql.Open(dialect.SQLite, "file:zupfmanager.db?mode=rwc&cache=shared&_fk=1")
	if err != nil {
		return nil, err
	}

	// Create ent client
	client := ent.NewClient(ent.Driver(driver))
	clnt := &Client{Client: client}
	err = clnt.Init()
	if err != nil {
		return nil, err
	}

	return clnt, nil
}

// Init initializes the database
func (c *Client) Init() error {
	// Run migrations
	slog.Info("Running database migrations")
	if err := c.Schema.Create(
		context.Background(),
	); err != nil {
		slog.Error("Database migration failed", "error", err)
		return err
	}
	slog.Info("Database initialized/migrated successfully")

	return nil
}

// CreateOrUpdateProject creates a new project or updates an existing project
func (c *Client) CreateOrUpdateProject(ctx context.Context, projectID int, title, shortName string, config map[string]interface{}) (*ent.Project, error) {
	if projectID == 0 {
		// Create a new project
		project, err := c.Project.Create().
			SetTitle(title).
			SetShortName(shortName).
			SetConfig(config).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		slog.Info("Created new project", "id", project.ID, "title", project.Title)
		return project, nil
	} else {
		// Update an existing project
		project, err := c.Project.Get(ctx, projectID)
		if err != nil {
			return nil, err
		}
		project, err = project.Update().
			SetTitle(title).
			SetShortName(shortName).
			SetConfig(config).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		slog.Info("Updated project", "id", projectID, "title", project.Title)
		return project, nil
	}
}

// GetProject gets a project by ID
func (c *Client) GetProject(ctx context.Context, projectID int) (*ent.Project, error) {
	project, err := c.Project.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return project, nil
}

// Close closes the database connection
func (c *Client) Close() error {
	return c.Client.Close()
}
