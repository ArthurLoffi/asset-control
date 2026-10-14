package database

import (
	"context"

	"asset-control/internal/core/dto"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ResourceConfig describes one entity for ResourceAdapter.
type ResourceConfig[TRequest dto.Paginated, TInput any, TDocument any, TResponse any] struct {
	ToDocument func(TInput) TDocument
	ToResponse func(TDocument) TResponse
	// Filter turns the entity's optional filters into a query; nil means none.
	Filter func(TRequest) bson.M
	// SortColumns maps public sort_by values to stored fields.
	SortColumns map[string]string
	// AfterDeactivate runs in the deactivation transaction when set.
	AfterDeactivate func(ctx context.Context, id uint64) error
}

// ResourceAdapter implements resource.DataPort over a Collection, so each
// entity only declares its mapping, filters, and sortable columns.
type ResourceAdapter[TRequest dto.Paginated, TInput any, TDocument any, TResponse any] struct {
	collection Collection[TDocument]
	config     ResourceConfig[TRequest, TInput, TDocument, TResponse]
}

func NewResourceAdapter[TRequest dto.Paginated, TInput any, TDocument any, TResponse any](
	collection Collection[TDocument], config ResourceConfig[TRequest, TInput, TDocument, TResponse],
) *ResourceAdapter[TRequest, TInput, TDocument, TResponse] {
	return &ResourceAdapter[TRequest, TInput, TDocument, TResponse]{collection: collection, config: config}
}

func (a *ResourceAdapter[TRequest, TInput, TDocument, TResponse]) FindAll(ctx context.Context, request TRequest) (dto.PaginatedResponseDTO[TResponse], error) {
	filter := bson.M{}
	if a.config.Filter != nil {
		filter = a.config.Filter(request)
	}
	pagination := request.Pagination()
	page, err := a.collection.FindPage(ctx, filter, Sort(pagination, a.config.SortColumns), pagination)
	if err != nil {
		return dto.PaginatedResponseDTO[TResponse]{}, err
	}
	return dto.MapPage(page, a.config.ToResponse), nil
}

func (a *ResourceAdapter[TRequest, TInput, TDocument, TResponse]) FindByID(ctx context.Context, id uint64) (TResponse, error) {
	return a.response(a.collection.FindByID(ctx, id))
}

func (a *ResourceAdapter[TRequest, TInput, TDocument, TResponse]) Create(ctx context.Context, input TInput) (TResponse, error) {
	return a.response(a.collection.Create(ctx, a.config.ToDocument(input)))
}

func (a *ResourceAdapter[TRequest, TInput, TDocument, TResponse]) Update(ctx context.Context, id uint64, input TInput) (TResponse, error) {
	return a.response(a.collection.Update(ctx, id, a.config.ToDocument(input)))
}

func (a *ResourceAdapter[TRequest, TInput, TDocument, TResponse]) Deactivate(ctx context.Context, id uint64) (dto.DeactivationResponseDTO, error) {
	if a.config.AfterDeactivate == nil {
		return a.collection.Deactivate(ctx, id)
	}
	var deactivation dto.DeactivationResponseDTO
	err := WithinTransaction(ctx, a.collection.Database(), func(txCtx context.Context) error {
		var err error
		if deactivation, err = a.collection.Deactivate(txCtx, id); err != nil {
			return err
		}
		return a.config.AfterDeactivate(txCtx, id)
	})
	return deactivation, err
}

func (a *ResourceAdapter[TRequest, TInput, TDocument, TResponse]) FindInactive(ctx context.Context, request dto.PaginatedRequestDTO) (dto.PaginatedResponseDTO[dto.DeactivationResponseDTO], error) {
	return a.collection.FindInactivePage(ctx, request)
}

func (a *ResourceAdapter[TRequest, TInput, TDocument, TResponse]) response(document TDocument, err error) (TResponse, error) {
	if err != nil {
		var zero TResponse
		return zero, err
	}
	return a.config.ToResponse(document), nil
}
