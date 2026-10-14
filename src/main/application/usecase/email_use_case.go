package usecase

import (
	"context"

	"asset-control/src/main/application/ports"
	"asset-control/src/main/domain/dto"
)

type EmailUseCase struct {
	repository ports.EmailsPorts
	changes    ports.ChangePorts
}

func NewEmailUseCase(repository ports.EmailsPorts, changePorts ...ports.ChangePorts) *EmailUseCase {
	var changes ports.ChangePorts
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
	return updateWithChange(
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
