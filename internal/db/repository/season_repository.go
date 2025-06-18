package repository

import (
	"gorgon/internal/db/model"

	"github.com/jmoiron/sqlx"
)

type SeasonRepositoryInterface interface {
	Create(season model.Season) (int64, error)
	CreateTx(tx *sqlx.Tx, season model.Season) (int64, error)
	DeleteById(id int64) error
	GetById(id int64) (model.Season, error)
	List() ([]model.Season, error)
	ListByShowId(showID int64) ([]model.Season, error)
}

type SeasonRepository struct {
	db *sqlx.DB
}

// TODO: UPDATE METHOD
func NewSeasonRepository(db *sqlx.DB) *SeasonRepository {
	return &SeasonRepository{
		db: db,
	}
}

func (s *SeasonRepository) Create(season model.Season) (int64, error) {
	query := `
	INSERT INTO seasons (
		show_id, 
		season_number
	) 
	VALUES (
		:show_id, 
		:season_number
	) 
	`
	result, err := s.db.NamedExec(query, season)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SeasonRepository) CreateTx(tx *sqlx.Tx, season model.Season) (int64, error) {
	query := `
	INSERT INTO seasons (
		show_id, 
		season_number
	) 
	VALUES (
		:show_id, 
		:season_number
	) 
	`
	result, err := tx.NamedExec(query, season)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, nil
}

func (s *SeasonRepository) GetById(id int64) (model.Season, error) {
	var season model.Season
	if err := s.db.Get(&season, "SELECT * FROM seasons WHERE id = ? LIMIT 1", id); err != nil {
		return model.Season{}, err
	}
	return season, nil
}

func (s *SeasonRepository) DeleteById(id int64) error {
	if _, err := s.db.Exec("DELETE FROM seasons WHERE id = ?", id); err != nil {
		return err
	}
	return nil
}

func (s *SeasonRepository) List() ([]model.Season, error) {
	var seasons []model.Season
	if err := s.db.Select(&seasons, "SELECT * FROM seasons"); err != nil {
		return []model.Season{}, err
	}
	return seasons, nil
}

func (s *SeasonRepository) ListByShowId(showID int64) ([]model.Season, error) {
	var seasons []model.Season
	if err := s.db.Select(&seasons, "SELECT * FROM seasons WHERE show_id = ?", showID); err != nil {
		return []model.Season{}, err
	}
	return seasons, nil
}

var _ SeasonRepositoryInterface = (*SeasonRepository)(nil)
