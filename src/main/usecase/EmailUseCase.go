package usecase

import (
	"context"

	"asset-control/src/main/domain/dto"
	"asset-control/src/main/domain/ports"
)

type EmailUseCase struct{ repository ports.EmailsPorts }

func NewEmailUseCase(repository ports.EmailsPorts) *EmailUseCase {
	return &EmailUseCase{repository: repository}
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
	return uc.repository.UpdateEmail(ctx, id, input)
}

func (uc *EmailUseCase) FindByID(ctx context.Context, id int32) (dto.EmailDto, error) {
	return uc.repository.FindByID(ctx, id)
}

func (uc *EmailUseCase) FindByUser(ctx context.Context, user string) (dto.EmailDto, error) {
	return uc.repository.FindByUser(ctx, user)
}
