package repository

import (
	"github.com/jmoiron/sqlx"
)

type ShowSearchPatternsRepositoryInterface interface {
	Replace(showID int64, patterns []string) error
	GetByShowID(showID int64) ([]string, error)
}

type ShowSearchPatternsRepository struct {
	db *sqlx.DB
}

func NewShowSearchPatternsRepository(db *sqlx.DB) ShowSearchPatternsRepository {
	return ShowSearchPatternsRepository{
		db: db,
	}
}

func (s *ShowSearchPatternsRepository) Replace(showID int64, patterns []string) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	if _, err := tx.Exec("DELETE FROM show_search_patterns WHERE show_id = ?", showID); err != nil {
		tx.Rollback()
		return err
	}

	for position, pattern := range patterns {
		if _, err := tx.Exec(
			"INSERT INTO show_search_patterns (show_id, pattern, position) VALUES (?, ?, ?)",
			showID, pattern, position,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *ShowSearchPatternsRepository) GetByShowID(showID int64) ([]string, error) {
	var patterns []string
	if err := s.db.Select(&patterns, "SELECT pattern FROM show_search_patterns WHERE show_id = ? ORDER BY position ASC", showID); err != nil {
		return nil, err
	}
	if patterns == nil {
		patterns = []string{}
	}
	return patterns, nil
}

var _ ShowSearchPatternsRepositoryInterface = (*ShowSearchPatternsRepository)(nil)
