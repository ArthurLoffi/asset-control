package usecase

import (
	"asset-control/internal/core/audit"
	"context"

	"asset-control/internal/modules/inventory/application/port"
	"asset-control/internal/modules/inventory/domain/dto"
)

type EmailUseCase struct {
	repository port.EmailsPorts
	changes    audit.Recorder
}

func NewEmailUseCase(repository port.EmailsPorts, changePorts ...audit.Recorder) *EmailUseCase {
	var changes audit.Recorder
	if len(changePorts) > 0 {
		changes = changePorts[0]
	}
	return &EmailUseCase{repository: repository, changes: changes}
}

func (uc *EmailUseCase) GetAllEmails(ctx context.Context) ([]dto.EmailDto, error) {
	return uc.repository.GetAllEmails(ctx)
}

func (uc *EmailUseCase) CreateEmail(ctx context.Context, input dto.CreateEmailDto) (dto.EmailDto, error) {
	if err := input.Validate(); err != nil {
		return dto.EmailDto{}, err
	}
	return uc.repository.CreateEmail(ctx, input)
}

func (uc *EmailUseCase) UpdateEmail(ctx context.Context, id int32, input dto.CreateEmailDto) (dto.EmailDto, error) {
	if err := input.Validate(); err != nil {
		return dto.EmailDto{}, err
	}
	return audit.UpdateWithChange(
		ctx,
		uc.changes,
		"email",
		int64(id),
		func(txCtx context.Context) (dto.EmailDto, error) { return uc.repository.FindByID(txCtx, id) },
		func(txCtx context.Context) (dto.EmailDto, error) { return uc.repository.UpdateEmail(txCtx, id, input) },
	)
}

func (uc *EmailUseCase) FindByID(ctx context.Context, id int32) (dto.EmailDto, error) {
	return uc.repository.FindByID(ctx, id)
}

func (uc *EmailUseCase) FindByUser(ctx context.Context, user string) (dto.EmailDto, error) {
	return uc.repository.FindByUser(ctx, user)
}
