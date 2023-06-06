package service

import (
	"fmt"
	"gorm.io/gorm"
)

type Entity interface {
	TableName() string
}

type Repository struct {
	entity Entity
	db     *gorm.DB
}

type IRepository interface {
	FindOne(id int, dbTransaction *gorm.DB) (any, error)
	FindAll(applyFilterAndSort func(db *gorm.DB) *gorm.DB, dbTransaction *gorm.DB) ([]map[string]any, error)
	Create(form any, dbTransaction *gorm.DB) (any, error)
	Update(id int, form any, dbTransaction *gorm.DB) (any, error)
	Delete(id int, dbTransaction *gorm.DB) error
}

func NewRepository(entity Entity, db *gorm.DB) IRepository {
	return &Repository{
		entity: entity,
		db:     db,
	}
}

func (r *Repository) getDB(dbTransaction *gorm.DB) *gorm.DB {
	if dbTransaction != nil {
		return dbTransaction
	}
	return r.db
}

func (r *Repository) FindOne(id int, dbTransaction *gorm.DB) (any, error) {
	db := r.getDB(dbTransaction)

	entity := r.entity
	err := db.Model(entity).Where("id = ?", id).First(entity).Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *Repository) FindAll(applyFilterAndSort func(db *gorm.DB) *gorm.DB, dbTransaction *gorm.DB) ([]map[string]any, error) {
	db := r.getDB(dbTransaction)

	var count int64
	err := db.Model(r.entity).Count(&count).Error
	if err != nil {
		return nil, err
	}

	var entities []map[string]any
	query := db.Table(r.entity.TableName())
	query = applyFilterAndSort(query)
	res := query.Find(&entities)
	if res.Error != nil {
		fmt.Printf("error, %+v\n", res.Error)
		return nil, err
	}

	return entities, nil
}

func (r *Repository) Create(form any, dbTransaction *gorm.DB) (any, error) {
	db := r.getDB(dbTransaction)

	result := db.Table(r.entity.TableName()).Select("*").Create(form)
	if result.Error != nil {
		return nil, result.Error
	}

	return form, nil
}

func (r *Repository) Update(id int, form any, dbTransaction *gorm.DB) (any, error) {
	db := r.getDB(dbTransaction)

	entity := r.entity
	err := db.Model(entity).Where("id = ?", id).First(entity).Error
	if err != nil {
		return nil, err
	}

	result := db.Save(form)
	if result.Error != nil {
		return nil, result.Error
	}

	return entity, nil
}

func (r *Repository) Delete(id int, dbTransaction *gorm.DB) error {
	db := r.getDB(dbTransaction)

	entity := r.entity
	err := db.Model(entity).Where("id = ?", id).First(entity).Error
	if err != nil {
		return err
	}

	result := db.Delete(entity)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
