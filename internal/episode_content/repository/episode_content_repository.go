package repository

import (
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/internal/episode_content/model"

	"github.com/jmoiron/sqlx"
)

type EpisodeContentRepository struct {
	db *sqlx.DB
}

func NewEpisodeContentRepository() EpisodeContentRepository {
	return EpisodeContentRepository{
		db: config.GetSQLite(),
	}
}

func (s *EpisodeContentRepository) Create(content model.EpisodeContent) error {
	query := `
	INSERT INTO episode_content (
		episode_id,
		name,
	        file_path,
		size,
		is_seed
	) 
	VALUES (
		:episode_id,
		:name,
		:file_path,
		:size,
		:is_seed
	) 
	`
	if _, err := s.db.NamedExec(query, content); err != nil {
		return err
	}

	return nil
}

func (s *EpisodeContentRepository) CreateTx(tx *sqlx.Tx, content model.EpisodeContent) error {
	query := `
	INSERT INTO episode_content (
		episode_id,
		name,
	        file_path,
		size,
		is_seed
	) 
	VALUES (
		:episode_id,
		:name,
		:file_path,
		:size,
		:is_seed
	) 
	`
	if _, err := tx.NamedExec(query, content); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *EpisodeContentRepository) GetById(id int) (model.EpisodeContent, error) {
	var content model.EpisodeContent
	if err := s.db.Get(&content, "SELECT * FROM episode_content WHERE id = ? LIMIT 1", id); err != nil {
		return model.EpisodeContent{}, err
	}
	return content, nil
}

func (s *EpisodeContentRepository) GetByEpisodeId(episodeId int64) (model.EpisodeContent, error) {
	var content model.EpisodeContent
	if err := s.db.Get(&content, "SELECT * FROM episode_content WHERE episode_id = ? LIMIT 1", episodeId); err != nil {
		return model.EpisodeContent{}, err
	}
	return content, nil
}

func (s *EpisodeContentRepository) DeleteById(id int) error {
	if _, err := s.db.Exec("DELETE FROM episode_content WHERE id = ?", id); err != nil {
		return err
	}
	return nil
}

func (s *EpisodeContentRepository) List() ([]model.EpisodeContent, error) {
	var contents []model.EpisodeContent
	if err := s.db.Select(&contents, "SELECT * FROM episodes"); err != nil {
		return []model.EpisodeContent{}, err
	}
	return contents, nil
}

func (s *EpisodeContentRepository) ListByEpisodeId(episodeID int64) ([]model.EpisodeContent, error) {
	var contents []model.EpisodeContent
	if err := s.db.Select(&contents, "SELECT * FROM episode_content WHERE episode_id = ?", episodeID); err != nil {
		return []model.EpisodeContent{}, err
	}
	return contents, nil
}

func (s *EpisodeContentRepository) Update(content model.EpisodeContent) error {
	query := `
	UPDATE shows SET
		episode_id = :episode_id,
		name = :name,
		size = :size,
		file_path = :file_path,
		is_seed = :is_seed,
	WHERE id = :id
	`
	_, err := s.db.NamedExec(query, content)
	if err != nil {
		return err
	}
	return nil
}

func (s *EpisodeContentRepository) UpdateTx(tx *sqlx.Tx, content model.EpisodeContent) error {
	query := `
	UPDATE shows SET
		episode_id = :episode_id,
		name = :name,
		file_path = :file_path,
		size = :size,
		is_seed = :is_seed,
	WHERE id = :id
	`
	_, err := tx.NamedExec(query, content)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
