package httpx

import (
	"context"
	"net/http"
	"strconv"

	"asset-control/internal/core/dto"
	"asset-control/internal/core/exception"
	"asset-control/internal/core/validation"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// The handlers below only do transport: read parameters, call one use case,
// and write the response or the error envelope.

// List serves GET /<resources>: binds and validates the query, checks sort_by
// against the entity's allow list, and answers {content, paging}.
func List[TRequest dto.Paginated, TResponse any](c *gin.Context, sortableFields []string,
	execute func(context.Context, TRequest) (dto.PaginatedResponseDTO[TResponse], error)) {
	var request TRequest
	if err := BindQuery(c, &request); err != nil {
		WriteError(c, err)
		return
	}
	if err := request.Pagination().ValidateSortBy(sortableFields); err != nil {
		WriteError(c, err)
		return
	}
	body, err := execute(c.Request.Context(), request)
	respond(c, http.StatusOK, body, err)
}

// Get serves GET /<resources>/{id}.
func Get[TResponse any](c *gin.Context, execute func(context.Context, uint64) (TResponse, error)) {
	id, err := PathID(c)
	if err != nil {
		WriteError(c, err)
		return
	}
	body, err := execute(c.Request.Context(), id)
	respond(c, http.StatusOK, body, err)
}

// Create serves POST /<resources> and answers 201 with the saved record.
func Create[TInput any, TResponse any](c *gin.Context, execute func(context.Context, TInput) (TResponse, error)) {
	var input TInput
	if err := BindJSON(c, &input); err != nil {
		WriteError(c, err)
		return
	}
	body, err := execute(c.Request.Context(), input)
	respond(c, http.StatusCreated, body, err)
}

// Update serves PUT /<resources>/{id}.
func Update[TInput any, TResponse any](c *gin.Context, execute func(context.Context, uint64, TInput) (TResponse, error)) {
	id, err := PathID(c)
	if err != nil {
		WriteError(c, err)
		return
	}
	var input TInput
	if err := BindJSON(c, &input); err != nil {
		WriteError(c, err)
		return
	}
	body, err := execute(c.Request.Context(), id, input)
	respond(c, http.StatusOK, body, err)
}

func respond[T any](c *gin.Context, status int, body T, err error) {
	if err != nil {
		WriteError(c, err)
		return
	}
	c.JSON(status, body)
}

// BindQuery maps ?key=value (repeated keys fill slices) by `form` tags, applies
// UnmarshalParam (for example the direction upper-casing), ignores unknown
// parameters, and validates `binding` tags.
func BindQuery(c *gin.Context, target any) error {
	if err := binding.MapFormWithTag(target, c.Request.URL.Query(), "form"); err != nil {
		return exception.NewValidationError("query parameters are invalid: " + err.Error())
	}
	return validation.Struct(target)
}

// PathID reads a positive {id} that fits BSON's signed 64-bit integers.
func PathID(c *gin.Context) (uint64, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 63)
	if err != nil || id == 0 {
		return 0, exception.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	return id, nil
}
