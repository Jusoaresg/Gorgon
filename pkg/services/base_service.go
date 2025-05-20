package services

import (
	"fmt"
	"gorgon/config"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseService struct {
	DB *gorm.DB
}

func NewBaseService() (b *BaseService) {
	db := config.GetSQLite()
	if db == nil {
		panic("Database not initialized")
	}
	return &BaseService{DB: db}
}

func (b *BaseService) Add(data any) error {
	return b.DB.Create(data).Error
}

func (b *BaseService) Get(model any, id int) error {
	return b.DB.First(model, "id = ?", id).Error
}

func (b *BaseService) UpdateByID(id int, model any) error {
	return b.DB.Model(model).Where("id = ?", id).Updates(model).Error
}

func (b *BaseService) UpdateByIDWithSelect(id int, model any, selects ...string) error {
	return b.DB.Model(model).Where("id = ?", id).Select(selects).Updates(model).Error
}

func (b *BaseService) GetWithPreload(model any, id int, relations ...string) error {
	query := b.DB
	for _, relation := range relations {
		query = query.Preload(relation)
	}
	return query.First(model, "id = ?", id).Error
}

func (b *BaseService) GetShowsByIdentification(model any, identification string, ids []int) error {
	return b.DB.Where(fmt.Sprintf("%s IN ?", identification), ids).Find(model).Error
}

// Always use a slice as model
func (b *BaseService) List(model any) error {
	return b.DB.Find(model).Error
}

func (b *BaseService) ListWithIdentification(model any, identification string, value string) error {
	return b.DB.Where(fmt.Sprintf("%s = ?", identification), value).Find(model).Error
}

// Always use a slice as model
func (b *BaseService) ListWithPreload(model any, relations ...string) error {
	query := b.DB
	for _, relation := range relations {
		query = query.Preload(relation)
	}
	return query.Find(model).Error
}

func (b *BaseService) ListByNameWithPreload(name string, model any, relations ...string) error {
	query := b.DB
	for _, relation := range relations {
		query = query.Preload(relation)
	}
	query.Where("Name LIKE ?", "%"+name+"%")
	return query.Find(model).Error
}

func (b *BaseService) Delete(id int, model any) error {
	return b.DB.Unscoped().Model(model).Where("id = ?", id).Select(clause.Associations).Delete(model).Error
}

func (b *BaseService) DeletePermanently(id int, model any) error {
	return b.DB.Unscoped().Where("aid = ?", id).Select(clause.Associations).Delete(&model).Error
}
