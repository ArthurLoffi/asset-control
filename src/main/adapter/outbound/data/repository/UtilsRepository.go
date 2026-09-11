package repository

import (
	"asset-control/src/main/adapter/outbound/data/entities"

	"gorm.io/gorm"
)

type UtilsRepository struct {
	db *gorm.DB
}

func NewUtilsRepository(db *gorm.DB) *UtilsRepository {
	return &UtilsRepository{
		db: db,
	}
}

func (ur *UtilsRepository) BranchIdToDescription(id uint64) (string, error) {
	var description string

	err := ur.db.Model(&entities.Branch{}).
		Select("name").
		Where("id = ?", id).
		First(&description).Error

	if err != nil {
		return "", err
	}

	return description, nil
}

func (ur *UtilsRepository) DepartmentIdToDescription(id uint64) (string, error) {
	var description string

	err := ur.db.Model(&entities.Department{}).
		Select("name").
		Where("id = ?", id).
		First(&description).Error

	if err != nil {
		return "", err
	}

	return description, nil
}
