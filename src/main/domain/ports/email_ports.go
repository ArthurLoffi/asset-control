package ports

import (
	"context"

	"asset-control/src/main/domain/dto"
)

type EmailsPorts interface {
	GetAllEmails(ctx context.Context) ([]dto.EmailDto, error)
	CreateEmail(ctx context.Context, input dto.CreateEmailDto) (dto.EmailDto, error)
	UpdateEmail(ctx context.Context, id int32, input dto.CreateEmailDto) (dto.EmailDto, error)
	FindByID(ctx context.Context, id int32) (dto.EmailDto, error)
	FindByUser(ctx context.Context, user string) (dto.EmailDto, error)
}
