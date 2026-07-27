package repository

import (
	"github.com/jusoaresg/gorgon/internal/indexer/model"

	"github.com/jmoiron/sqlx"
)

type IndexerRepository struct {
	db *sqlx.DB
}

func NewIndexerRepository(db *sqlx.DB) *IndexerRepository {
	return &IndexerRepository{
		db: db,
	}
}

func (s *IndexerRepository) Create(indexer model.Indexer) error {
	query := `
	INSERT INTO indexers (
		indexer_id, 
		name,
		enabled,
		definition_name,
		indexers_urls,
		language
	) 
	VALUES (
		:indexer_id, 
		:name,
		:enabled,
		:definition_name,
		:indexers_urls,
		:language
	) 
	`
	if _, err := s.db.NamedExec(query, indexer); err != nil {
		return err
	}

	return nil
}

func (s *IndexerRepository) GetById(id int) (model.Indexer, error) {
	var indexer model.Indexer
	if err := s.db.Get(&indexer, "SELECT * FROM indexers WHERE id = ? LIMIT 1", id); err != nil {
		return model.Indexer{}, err
	}
	return indexer, nil
}

func (s *IndexerRepository) DeleteById(id int) error {
	if _, err := s.db.Exec("DELETE FROM indexers WHERE id = ?", id); err != nil {
		return err
	}
	return nil
}

func (s *IndexerRepository) List() ([]model.Indexer, error) {
	var indexers []model.Indexer
	if err := s.db.Select(&indexers, "SELECT * FROM seasons"); err != nil {
		return []model.Indexer{}, err
	}
	return indexers, nil
}
