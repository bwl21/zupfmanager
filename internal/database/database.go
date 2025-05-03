package database

import (
	"context"
	"log/slog"

	"github.com/bwl21/zupfmanager/internal/ent"
	"github.com/bwl21/zupfmanager/internal/ent/migrate"

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
	if err := c.Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		return err
	}

	slog.Debug("Database initialized successfully")
	return nil
}

// Close closes the database connection
func (c *Client) Close() error {
	return c.Client.Close()
}
