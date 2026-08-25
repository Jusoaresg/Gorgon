package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/internal/show_aliases/model"
)

var ErrAliasNotFound = errors.New("alias not found")

type ShowAliasesRepositoryInterface interface {
	Create(show model.ShowAlias) (int64, error)
	CreateTx(tx *sqlx.Tx, alias model.ShowAlias) (int64, error)
	UpdateTx(tx *sqlx.Tx, alias model.ShowAlias) error
	GetByID(id int64) (model.ShowAlias, error)
	DeleteByID(id int64) error
	ListByShowID(show_id int64) ([]model.ShowAlias, error)
}

type ShowAliasesRepository struct {
	db *sqlx.DB
}

func NewShowAliasesRepository(db *sqlx.DB) ShowAliasesRepository {
	return ShowAliasesRepository{
		db: db,
	}
}

func (s *ShowAliasesRepository) Create(alias model.ShowAlias) (int64, error) {
	query := `
	INSERT INTO show_aliases (
		show_id,
		alias,
		country,
		source
	) 
	VALUES (
		:show_id,
		:alias,
		:country,
		:source
	) 
	`
	result, err := s.db.NamedExec(query, alias)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *ShowAliasesRepository) CreateTx(tx *sqlx.Tx, alias model.ShowAlias) (int64, error) {
	query := `
	INSERT INTO show_aliases (
		show_id,
		alias,
		country,
		source
	) 
	VALUES (
		:show_id,
		:alias,
		:country,
		:source
	) 
	`
	result, err := tx.NamedExec(query, alias)
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

func (s *ShowAliasesRepository) UpdateTx(tx *sqlx.Tx, alias model.ShowAlias) error {
	query := `
	UPDATE show_aliases SET
		alias = :alias,
		country = :country
	WHERE id = :id
	`
	result, err := tx.NamedExec(query, alias)
	if err != nil {
		tx.Rollback()
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if rows == 0 {
		tx.Rollback()
		return ErrAliasNotFound
	}

	return nil
}

func (s *ShowAliasesRepository) GetByID(id int64) (model.ShowAlias, error) {
	var alias model.ShowAlias
	if err := s.db.Get(&alias, "SELECT * FROM show_aliases WHERE id = ? LIMIT 1", id); err != nil {
		return model.ShowAlias{}, err
	}
	return alias, nil
}

func (s *ShowAliasesRepository) DeleteByID(id int64) error {
	result, err := s.db.Exec("DELETE FROM show_aliases WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrAliasNotFound
	}

	return nil
}

func (s *ShowAliasesRepository) ListByShowID(show_id int64) ([]model.ShowAlias, error) {
	var aliases []model.ShowAlias
	if err := s.db.Select(&aliases, "SELECT * FROM show_aliases where show_id = ?", show_id); err != nil {
		return []model.ShowAlias{}, err
	}
	return aliases, nil
}

var _ ShowAliasesRepositoryInterface = (*ShowAliasesRepository)(nil)
