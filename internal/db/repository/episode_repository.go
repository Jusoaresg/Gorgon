package repository

import (
	"gorgon/config"
	"gorgon/internal/db/model"

	"github.com/jmoiron/sqlx"
)

type EpisodeRepository struct {
	db *sqlx.DB
}

func NewEpisodeRepository() EpisodeRepository {
	return EpisodeRepository{
		db: config.GetSQLite(),
	}
}

func (s *EpisodeRepository) Create(episode model.Episode) error {
	query := `
	INSERT INTO episodes (
		show_id, 
		season_id,
		name, 
		summary,
		type,
		number,
		season,
		airstamp,
		file_path,
		tracking,
		torrent_hash
	) 
	VALUES (
		:show_id, 
		:season_id,
		:name, 
		:summary,
		:type,
		:number,
		:season,
		:airstamp,
		:file_path,
		:tracking,
		:torrent_hash
	) 
	`
	if _, err := s.db.NamedExec(query, episode); err != nil {
		return err
	}

	return nil
}

func (s *EpisodeRepository) CreateTx(tx *sqlx.Tx, episode model.Episode) error {
	query := `
	INSERT INTO episodes (
		show_id, 
		season_id,
		name, 
		summary,
		type,
		number,
		season,
		airstamp,
		file_path,
		tracking,
		torrent_hash
	) 
	VALUES (
		:show_id, 
		:season_id,
		:name, 
		:summary,
		:type,
		:number,
		:season,
		:airstamp,
		:file_path,
		:tracking,
		:torrent_hash
	) 
	`
	if _, err := tx.NamedExec(query, episode); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *EpisodeRepository) GetById(id int) (model.Episode, error) {
	var episode model.Episode
	if err := s.db.Get(&episode, "SELECT * FROM episodes WHERE id = ? LIMIT 1", id); err != nil {
		return model.Episode{}, err
	}
	return episode, nil
}

func (s *EpisodeRepository) DeleteById(id string) error {
	if _, err := s.db.Exec("DELETE FROM episodes WHERE id = ?", id); err != nil {
		return err
	}
	return nil
}

func (s *EpisodeRepository) List() ([]model.Episode, error) {
	var episodes []model.Episode
	if err := s.db.Select(&episodes, "SELECT * FROM episodes"); err != nil {
		return []model.Episode{}, err
	}
	return episodes, nil
}

func (s *EpisodeRepository) ListByTracking(tracking string) ([]model.Episode, error) {
	var episodes []model.Episode
	if err := s.db.Select(&episodes, "SELECT * FROM episodes WHERE tracking = ?", tracking); err != nil {
		return []model.Episode{}, err
	}
	return episodes, nil
}

func (s *EpisodeRepository) ListByShowId(showID int64) ([]model.Episode, error) {
	var episodes []model.Episode
	if err := s.db.Select(&episodes, "SELECT * FROM episodes WHERE show_id = ?", showID); err != nil {
		return []model.Episode{}, err
	}
	return episodes, nil
}

func (s *EpisodeRepository) Update(episode model.Episode) error {
	query := `
	UPDATE episodes SET
		show_id = :show_id,
		season_id = :season_id,
		name = :name,
		summary = :summary,
		type = :type,
		number = :number,
		season = :season,
		airstamp = :airstamp,
		file_path = :file_path,
		tracking = :tracking,
		torrent_hash = :torrent_hash
	WHERE id = :id
	`
	_, err := s.db.NamedExec(query, episode)
	if err != nil {
		return err
	}
	return nil
}

func (s *EpisodeRepository) UpdateTx(tx *sqlx.Tx, episode model.Episode) error {
	query := `
	UPDATE episodes SET
		show_id = :show_id,
		season_id = :season_id,
		name = :name,
		summary = :summary,
		type = :type,
		number = :number,
		season = :season,
		airstamp = :airstamp,
		file_path = :file_path,
		tracking = :tracking,
		torrent_hash = :torrent_hash
	WHERE id = :id
	`
	_, err := tx.NamedExec(query, episode)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
