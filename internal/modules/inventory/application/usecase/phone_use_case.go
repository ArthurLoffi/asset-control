package usecase

import (
	"asset-control/internal/core/audit"
	"context"

	"asset-control/internal/modules/inventory/application/port"
	"asset-control/internal/modules/inventory/domain/dto"
)

type PhoneUseCase struct {
	repository port.PhonePorts
	changes    audit.Recorder
}

func NewPhoneUseCase(repository port.PhonePorts, changePorts ...audit.Recorder) *PhoneUseCase {
	var changes audit.Recorder
	if len(changePorts) > 0 {
		changes = changePorts[0]
	}
	return &PhoneUseCase{repository: repository, changes: changes}
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
	return audit.UpdateWithChange(
		ctx,
		uc.changes,
		"phone",
		int64(id),
		func(txCtx context.Context) (dto.PhoneDto, error) { return uc.repository.FindByID(txCtx, id) },
		func(txCtx context.Context) (dto.PhoneDto, error) { return uc.repository.UpdatePhone(txCtx, id, input) },
	)
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
