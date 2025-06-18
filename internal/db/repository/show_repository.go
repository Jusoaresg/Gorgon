package repository

import (
	"errors"
	"gorgon/internal/db/model"

	"github.com/jmoiron/sqlx"
)

var ErrShowNotFound = errors.New("show not found")

type ShowRepositoryInterface interface {
	Create(show model.Show) (int64, error)
	CreateTx(tx *sqlx.Tx, show model.Show) (int64, error)
	GetById(id int64) (model.Show, error)
	GetByTvMazeID(tvMazeID int64) (model.Show, error)
	DeleteById(id int64) error
	List() ([]model.Show, error)
	UpdateByTvMazeID(show model.Show) error
	UpdateTxByTvMazeID(tx *sqlx.Tx, show model.Show) error
}

type ShowRepository struct {
	db *sqlx.DB
}

func NewShowRepository(db *sqlx.DB) *ShowRepository {
	return &ShowRepository{
		db: db,
	}
}

func (s *ShowRepository) Create(show model.Show) (int64, error) {
	query := `
	INSERT INTO shows (
		tv_maze_id, 
		name, 
		type, 
		language, 
		status, 
		premiered, 
		ended, 
		rating, 
		summary, 
		updated, 
		tv_rage, 
		the_tvdbd, 
		imdb, 
		image_original, 
		image_medium, 
		genres
	) 
	VALUES (
		:tv_maze_id, 
		:name, 
		:type, 
		:language, 
		:status, 
		:premiered, 
		:ended, 
		:rating, 
		:summary, 
		:updated, 
		:tv_rage, 
		:the_tvdbd, 
		:imdb, 
		:image_original, 
		:image_medium, 
		:genres
	) 
	`
	result, err := s.db.NamedExec(query, show)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *ShowRepository) CreateTx(tx *sqlx.Tx, show model.Show) (int64, error) {
	query := `
	INSERT INTO shows (
		tv_maze_id, 
		name, 
		type, 
		language, 
		status, 
		premiered, 
		ended, 
		rating, 
		summary, 
		updated, 
		tv_rage, 
		the_tvdbd, 
		imdb, 
		image_original, 
		image_medium, 
		genres
	) 
	VALUES (
		:tv_maze_id, 
		:name, 
		:type, 
		:language, 
		:status, 
		:premiered, 
		:ended, 
		:rating, 
		:summary, 
		:updated, 
		:tv_rage, 
		:the_tvdbd, 
		:imdb, 
		:image_original, 
		:image_medium, 
		:genres
	) 
	`
	result, err := tx.NamedExec(query, show)
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

func (s *ShowRepository) GetById(id int64) (model.Show, error) {
	var show model.Show
	if err := s.db.Get(&show, "SELECT * FROM shows WHERE id = ? LIMIT 1", id); err != nil {
		return model.Show{}, err
	}
	return show, nil
}

func (s *ShowRepository) GetByTvMazeID(tvMazeID int64) (model.Show, error) {
	var show model.Show
	if err := s.db.Get(&show, "SELECT * FROM shows WHERE tv_maze_id = ? LIMIT 1", tvMazeID); err != nil {
		return model.Show{}, err
	}
	return show, nil
}

func (s *ShowRepository) DeleteById(id int64) error {
	result, err := s.db.Exec("DELETE FROM shows WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrShowNotFound
	}

	return nil
}

func (s *ShowRepository) List() ([]model.Show, error) {
	var shows []model.Show
	if err := s.db.Select(&shows, "SELECT * FROM shows"); err != nil {
		return []model.Show{}, err
	}
	return shows, nil
}

func (s *ShowRepository) UpdateByTvMazeID(show model.Show) error {
	query := `
	UPDATE shows SET
		name = :name,
		type = :type,
		language = :language,
		status = :status,
		premiered = :premiered,
		ended = :ended,
		rating = :rating,
		summary = :summary,
		updated = :updated,
		tv_rage = :tv_rage,
		the_tvdbd = :the_tvdbd,
		imdb = :imdb,
		image_original = :image_original,
		image_medium = :image_medium,
		genres = :genres
	WHERE tv_maze_id = :tv_maze_id
	`

	result, err := s.db.NamedExec(query, show)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrShowNotFound
	}
	return nil
}

func (s *ShowRepository) UpdateTxByTvMazeID(tx *sqlx.Tx, show model.Show) error {
	query := `
	UPDATE shows SET
		name = :name,
		type = :type,
		language = :language,
		status = :status,
		premiered = :premiered,
		ended = :ended,
		rating = :rating,
		summary = :summary,
		updated = :updated,
		tv_rage = :tv_rage,
		the_tvdbd = :the_tvdbd,
		imdb = :imdb,
		image_original = :image_original,
		image_medium = :image_medium,
		genres = :genres
	WHERE tv_maze_id = :tv_maze_id
	`
	_, err := tx.NamedExec(query, show)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

var _ ShowRepositoryInterface = (*ShowRepository)(nil)
