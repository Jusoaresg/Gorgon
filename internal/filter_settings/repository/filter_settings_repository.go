package repository

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/internal/filter_settings/model"
)

const filterSettingsKey = "filters"

type FilterSettingsRepositoryInterface interface {
	Get() (model.FilterSettings, error)
	Save(settings model.FilterSettings) error
}

type FilterSettingsRepository struct {
	db *sqlx.DB
}

func NewFilterSettingsRepository(db *sqlx.DB) FilterSettingsRepository {
	return FilterSettingsRepository{
		db: db,
	}
}

func (s *FilterSettingsRepository) Get() (model.FilterSettings, error) {
	var value string
	err := s.db.Get(&value, "SELECT value FROM app_settings WHERE key = ? LIMIT 1", filterSettingsKey)
	if err != nil {
		return model.DefaultFilterSettings(), nil
	}

	var stored struct {
		DefaultFilterProfileID *int64 `json:"default_filter_profile_id"`
		UseAliases             *bool  `json:"use_aliases"`
		OnlyLatin              *bool  `json:"only_latin"`
	}
	if err := json.Unmarshal([]byte(value), &stored); err != nil {
		return model.DefaultFilterSettings(), nil
	}

	settings := model.DefaultFilterSettings()
	if stored.DefaultFilterProfileID != nil {
		settings.DefaultFilterProfileID = stored.DefaultFilterProfileID
	}
	if stored.UseAliases != nil {
		settings.UseAliases = *stored.UseAliases
	}
	if stored.OnlyLatin != nil {
		settings.OnlyLatin = *stored.OnlyLatin
	}
	return settings, nil
}

func (s *FilterSettingsRepository) Save(settings model.FilterSettings) error {
	value, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT INTO app_settings (key, value)
		VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, filterSettingsKey, string(value))
	return err
}

var _ FilterSettingsRepositoryInterface = (*FilterSettingsRepository)(nil)
