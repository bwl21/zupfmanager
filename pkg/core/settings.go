package core

import (
	"context"
	"fmt"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/bwl21/zupfmanager/internal/ent"
	"github.com/bwl21/zupfmanager/internal/ent/setting"
)

// settingsService implements SettingsService interface
type settingsService struct {
	db *database.Client
}

// NewSettingsService creates a new settings service
func NewSettingsService(db *database.Client) SettingsService {
	return &settingsService{
		db: db,
	}
}

// Get retrieves a setting value by key (alias for GetSetting)
func (s *settingsService) Get(ctx context.Context, key string) (string, error) {
	return s.GetSetting(ctx, key)
}

// Set sets a setting value by key (alias for SetSetting)
func (s *settingsService) Set(ctx context.Context, key, value string) error {
	return s.SetSetting(ctx, key, value)
}

// GetSetting retrieves a setting value by key
func (s *settingsService) GetSetting(ctx context.Context, key string) (string, error) {
	settingEntity, err := s.db.Setting.Query().
		Where(setting.KeyEQ(key)).
		Only(ctx)
	
	if err != nil {
		if ent.IsNotFound(err) {
			return "", nil // Return empty string if setting doesn't exist
		}
		return "", fmt.Errorf("failed to get setting %s: %w", key, err)
	}
	
	return settingEntity.Value, nil
}

// SetSetting sets a setting value by key
func (s *settingsService) SetSetting(ctx context.Context, key, value string) error {
	// Try to update existing setting first
	updated, err := s.db.Setting.Update().
		Where(setting.KeyEQ(key)).
		SetValue(value).
		Save(ctx)
	
	if err != nil {
		return fmt.Errorf("failed to update setting %s: %w", key, err)
	}
	
	// If no rows were updated, create new setting
	if updated == 0 {
		_, err = s.db.Setting.Create().
			SetKey(key).
			SetValue(value).
			Save(ctx)
		
		if err != nil {
			return fmt.Errorf("failed to create setting %s: %w", key, err)
		}
	}
	
	return nil
}

// GetLastImportPath retrieves the last import path
func (s *settingsService) GetLastImportPath(ctx context.Context) (string, error) {
	return s.GetSetting(ctx, "last_import_path")
}

// SetLastImportPath sets the last import path
func (s *settingsService) SetLastImportPath(ctx context.Context, path string) error {
	return s.SetSetting(ctx, "last_import_path", path)
}
