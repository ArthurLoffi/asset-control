package port

import (
	"context"

	"asset-control/internal/modules/inventory/domain/dto"
)

type PhonePorts interface {
	CreatePhone(ctx context.Context, input dto.CreatePhoneDto) (dto.PhoneDto, error)
	UpdatePhone(ctx context.Context, id uint64, input dto.CreatePhoneDto) (dto.PhoneDto, error)
	GetAllPhones(ctx context.Context) ([]dto.PhoneDto, error)
	FindByID(ctx context.Context, id uint64) (dto.PhoneDto, error)
	FindByEmployeeID(ctx context.Context, employeeID uint64) ([]dto.PhoneDto, error)
}
