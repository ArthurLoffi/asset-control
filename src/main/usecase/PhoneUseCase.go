package usecase

import (
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"
)

type PhoneUseCase struct {
	phonePorts ports.PhonePorts
}

func NewPhoneUseCase(phonePorts ports.PhonePorts) *PhoneUseCase {
	return &PhoneUseCase{
		phonePorts: phonePorts,
	}
}

func (uc *PhoneUseCase) CreatePhone(phone dto.CreatePhoneDto) error {
	return uc.phonePorts.CreatePhone(phone)
}

func (uc *PhoneUseCase) UpdatePhone(id uint64, phone dto.CreatePhoneDto) (dto.PhoneDto, error) {
	return uc.phonePorts.UpdatePhone(id, phone)
}

func (uc *PhoneUseCase) GetAllPhones() ([]dto.PhoneDto, error) {
	return uc.phonePorts.GetAllPhones()
}

func (uc *PhoneUseCase) FindByID(id uint64) (dto.PhoneDto, error) {
	return uc.phonePorts.FindByID(id)
}

func (uc *PhoneUseCase) FindByEmployeeID(employeeId uint64) ([]dto.PhoneDto, error) {
	return uc.phonePorts.FindByEmployeeID(employeeId)
}
