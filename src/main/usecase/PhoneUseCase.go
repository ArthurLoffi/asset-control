package usecase

import (
	"context"

	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"
)

type PhoneUseCase struct{ repository ports.PhonePorts }

func NewPhoneUseCase(repository ports.PhonePorts) *PhoneUseCase {
	return &PhoneUseCase{repository: repository}
}

func (uc *PhoneUseCase) CreatePhone(ctx context.Context, input dto.CreatePhoneDto) (dto.PhoneDto, error) {
	if err := input.Validate(); err != nil {
		return dto.PhoneDto{}, err
	}
	return uc.repository.CreatePhone(ctx, input)
}

func (uc *PhoneUseCase) UpdatePhone(ctx context.Context, id uint64, input dto.CreatePhoneDto) (dto.PhoneDto, error) {
	if err := input.Validate(); err != nil {
		return dto.PhoneDto{}, err
	}
	return uc.repository.UpdatePhone(ctx, id, input)
}

func (uc *PhoneUseCase) GetAllPhones(ctx context.Context) ([]dto.PhoneDto, error) {
	return uc.repository.GetAllPhones(ctx)
}

func (uc *PhoneUseCase) FindByID(ctx context.Context, id uint64) (dto.PhoneDto, error) {
	return uc.repository.FindByID(ctx, id)
}

func (uc *PhoneUseCase) FindByEmployeeID(ctx context.Context, employeeID uint64) ([]dto.PhoneDto, error) {
	return uc.repository.FindByEmployeeID(ctx, employeeID)
}
