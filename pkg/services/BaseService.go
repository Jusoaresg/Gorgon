package services

import (
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

func (b *BaseService) List(model interface{}) error {
	return b.DB.Find(model).Error
}

func (b *BaseService) Delete(id int, model interface{}) error {
	return b.DB.Model(model).Where("aid = ?", id).Delete(model).Error
}

func (b *BaseService) DeletePermanently(id int, model interface{}) error {
	return b.DB.Unscoped().Where("aid = ?", id).Delete(&model).Error
}
