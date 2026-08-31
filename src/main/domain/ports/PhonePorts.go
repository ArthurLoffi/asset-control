package ports

import "asset-control/src/main/domain/dto"

type PhonePorts interface {
	CreatePhone(phone dto.CreatePhoneDto) error

	UpdatePhone(id uint64, phone dto.CreatePhoneDto) (dto.PhoneDto, error)

	GetAllPhones() ([]dto.PhoneDto, error)

	FindByID(id uint64) (dto.PhoneDto, error)

	FindByEmployeeID(employeeId uint64) ([]dto.PhoneDto, error)
}