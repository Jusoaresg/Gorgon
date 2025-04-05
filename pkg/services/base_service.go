package services

import (
	"fmt"
	"gorgon/config"

	"gorm.io/gorm"
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

func (b *BaseService) Add(data interface{}) error {
	return b.DB.Create(data).Error
}

func (b *BaseService) Get(model interface{}, id int) error {
	return b.DB.First(model, "id = ?", id).Error
}

func (b *BaseService) GetWithPreload(model interface{}, id int, relations ...string) error {
	query := b.DB
	for _, relation := range relations {
		query = query.Preload(relation)
	}
	return query.First(model, "id = ?", id).Error
}

func (b *BaseService) GetShowsByIdentification(model interface{}, identification string, ids []int) error {
	return b.DB.Where(fmt.Sprintf("%s IN ?", identification), ids).Find(model).Error
}

// Always use a slice as model
func (b *BaseService) List(model interface{}) error {
	return b.DB.Find(model).Error
}

// Always use a slice as model
func (b *BaseService) ListWithPreload(model interface{}, relations ...string) error {
	query := b.DB
	for _, relation := range relations {
		query = query.Preload(relation)
	}
	return query.Find(model).Error
}

func (b *BaseService) ListByNameWithPreload(name string, model interface{}, relations ...string) error {
	query := b.DB
	for _, relation := range relations {
		query = query.Preload(relation)
	}
	query.Where("Name LIKE ?", "%"+name+"%")
	return query.Find(model).Error
}

func (b *BaseService) Delete(id int, model interface{}) error {
	return b.DB.Model(model).Where("id = ?", id).Delete(model).Error
}

func (b *BaseService) DeletePermanently(id int, model interface{}) error {
	return b.DB.Unscoped().Where("aid = ?", id).Delete(&model).Error
}
