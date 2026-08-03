package repository

import (
	"errors"

	"github.com/jusoaresg/gorgon/internal/episode_content/model"

	"github.com/jmoiron/sqlx"
)

var ErrShowNotFound = errors.New("show not found")

type EpisodeContentRepository struct {
	db *sqlx.DB
}

func NewEpisodeContentRepository(db *sqlx.DB) EpisodeContentRepository {
	return EpisodeContentRepository{
		db: db,
	}
}

func (s *EpisodeContentRepository) Create(content model.EpisodeContent) (int64, error) {
	query := `
	INSERT INTO episode_contents (
		episode_id,
		name,
	        file_path,
		size
	) 
	VALUES (
		:episode_id,
		:name,
		:file_path,
		:size
	) 
	`
	result, err := s.db.NamedExec(query, content)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *EpisodeContentRepository) CreateTx(tx *sqlx.Tx, content model.EpisodeContent) (int64, error) {
	query := `
	INSERT INTO episode_contents (
		episode_id,
		name,
	        file_path,
		size
	) 
	VALUES (
		:episode_id,
		:name,
		:file_path,
		:size
	) 
	`
	result, err := tx.NamedExec(query, content)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *EpisodeContentRepository) GetById(id int64) (model.EpisodeContent, error) {
	var content model.EpisodeContent
	if err := s.db.Get(&content, "SELECT * FROM episode_contents WHERE id = ? LIMIT 1", id); err != nil {
		return model.EpisodeContent{}, err
	}
	return content, nil
}

func (s *EpisodeContentRepository) GetByEpisodeId(episodeId int64) (model.EpisodeContent, error) {
	var content model.EpisodeContent
	if err := s.db.Get(&content, "SELECT * FROM episode_contents WHERE episode_id = ? LIMIT 1", episodeId); err != nil {
		return model.EpisodeContent{}, err
	}
	return content, nil
}

func (s *EpisodeContentRepository) DeleteById(id int64) error {
	if _, err := s.db.Exec("DELETE FROM episode_contents WHERE id = ?", id); err != nil {
		return err
	}
	return nil
}

func (s *EpisodeContentRepository) List() ([]model.EpisodeContent, error) {
	var contents []model.EpisodeContent
	if err := s.db.Select(&contents, "SELECT * FROM episode_contents"); err != nil {
		return []model.EpisodeContent{}, err
	}
	return contents, nil
}

func (s *EpisodeContentRepository) ListByEpisodeId(episodeID int64) ([]model.EpisodeContent, error) {
	var contents []model.EpisodeContent
	if err := s.db.Select(&contents, "SELECT * FROM episode_contents WHERE episode_id = ?", episodeID); err != nil {
		return []model.EpisodeContent{}, err
	}
	return contents, nil
}

func (s *EpisodeContentRepository) Update(content model.EpisodeContent) error {
	query := `
	UPDATE episode_contents SET
		episode_id = :episode_id,
		name = :name,
		size = :size,
		file_path = :file_path
	WHERE id = :id
	`
	result, err := s.db.NamedExec(query, content)
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

func (s *EpisodeContentRepository) UpdateTx(tx *sqlx.Tx, content model.EpisodeContent) error {
	query := `
	UPDATE episode_contents SET
		episode_id = :episode_id,
		name = :name,
		file_path = :file_path,
		size = :size
	WHERE id = :id
	`
	_, err := tx.NamedExec(query, content)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
