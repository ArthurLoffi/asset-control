package repository

import (
	"asset-control/src/main/adapter/outbound/data/entities"
	"asset-control/src/main/domain/dto"
	"errors"

	"gorm.io/gorm"
)

type PhoneRepository struct {
	db *gorm.DB
}

func NewPhoneRepository(db *gorm.DB) *PhoneRepository {
	return &PhoneRepository{db: db}
}

func (pr *PhoneRepository) CreatePhone(phone dto.CreatePhoneDto) error {
	phoneEntity := entities.Phone{
		Equipament:   phone.Equipament,
		Model:        phone.Model,
		IMEI:         phone.IMEI,
		SerialNumber: phone.SerialNumber,
		HasChip:      phone.HasChip,
		Situation:    phone.Situation,
		EmployeeID:   phone.EmployeeID,
	}

	if err := pr.db.Create(&phoneEntity).Error; err != nil {
		return err
	}

	return nil
}

func (pr *PhoneRepository) UpdatePhone(
	id uint64,
	phone dto.CreatePhoneDto,
) (dto.PhoneDto, error) {
	var phoneEntity entities.Phone

	if err := pr.db.First(&phoneEntity, id).Error; err != nil {
		return dto.PhoneDto{}, err
	}

	phoneEntity.Equipament = phone.Equipament
	phoneEntity.Model = phone.Model
	phoneEntity.IMEI = phone.IMEI
	phoneEntity.SerialNumber = phone.SerialNumber
	phoneEntity.HasChip = phone.HasChip
	phoneEntity.Situation = phone.Situation
	phoneEntity.EmployeeID = phone.EmployeeID

	if err := pr.db.Save(&phoneEntity).Error; err != nil {
		return dto.PhoneDto{}, err
	}

	return dto.PhoneDto{
		ID:           phoneEntity.ID,
		Equipament:   phoneEntity.Equipament,
		Model:        phoneEntity.Model,
		IMEI:         phoneEntity.IMEI,
		SerialNumber: phoneEntity.SerialNumber,
		HasChip:      phoneEntity.HasChip,
		Situation:    phoneEntity.Situation,
		EmployeeID:   phoneEntity.EmployeeID,
		CreatedAt:    phoneEntity.CreatedAt,
	}, nil
}

func (pr *PhoneRepository) GetAllPhones() ([]dto.PhoneDto, error) {
	var phoneEntities []entities.Phone

	if err := pr.db.Find(&phoneEntities).Error; err != nil {
		return nil, err
	}

	phones := make([]dto.PhoneDto, 0, len(phoneEntities))

	for _, phoneEntity := range phoneEntities {
		phones = append(phones, dto.PhoneDto{
			ID:           phoneEntity.ID,
			Equipament:   phoneEntity.Equipament,
			Model:        phoneEntity.Model,
			IMEI:         phoneEntity.IMEI,
			SerialNumber: phoneEntity.SerialNumber,
			HasChip:      phoneEntity.HasChip,
			Situation:    phoneEntity.Situation,
			EmployeeID:   phoneEntity.EmployeeID,
			CreatedAt:    phoneEntity.CreatedAt,
		})
	}

	return phones, nil
}

func (pr *PhoneRepository) FindByID(id uint64) (dto.PhoneDto, error) {
	var phoneEntity entities.Phone

	if err := pr.db.First(&phoneEntity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.PhoneDto{}, err
		}

		return dto.PhoneDto{}, err
	}

	return dto.PhoneDto{
		ID:           phoneEntity.ID,
		Equipament:   phoneEntity.Equipament,
		Model:        phoneEntity.Model,
		IMEI:         phoneEntity.IMEI,
		SerialNumber: phoneEntity.SerialNumber,
		HasChip:      phoneEntity.HasChip,
		Situation:    phoneEntity.Situation,
		EmployeeID:   phoneEntity.EmployeeID,
		CreatedAt:    phoneEntity.CreatedAt,
	}, nil
}

func (pr *PhoneRepository) FindByEmployeeID(employeeID uint64) ([]dto.PhoneDto, error) {
	var phoneEntities []entities.Phone

	if err := pr.db.
		Where("employee_id = ?", employeeID).
		Find(&phoneEntities).Error; err != nil {
		return nil, err
	}

	phones := make([]dto.PhoneDto, 0, len(phoneEntities))

	for _, phoneEntity := range phoneEntities {
		phones = append(phones, dto.PhoneDto{
			ID:           phoneEntity.ID,
			Equipament:   phoneEntity.Equipament,
			Model:        phoneEntity.Model,
			IMEI:         phoneEntity.IMEI,
			SerialNumber: phoneEntity.SerialNumber,
			HasChip:      phoneEntity.HasChip,
			Situation:    phoneEntity.Situation,
			EmployeeID:   phoneEntity.EmployeeID,
			CreatedAt:    phoneEntity.CreatedAt,
		})
	}

	return phones, nil
}
