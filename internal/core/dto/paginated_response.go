package dto

type PagingDTO struct {
	TotalElements int64 `json:"total_elements"`
	Size          int64 `json:"size"`
	Page          int64 `json:"page"`
	TotalPages    int64 `json:"total_pages"`
	HasContent    bool  `json:"has_content"`
	HasNext       bool  `json:"has_next"`
	HasPrevious   bool  `json:"has_previous"`
}

type PaginatedResponseDTO[T any] struct {
	Content []T       `json:"content"`
	Paging  PagingDTO `json:"paging"`
}

func NewPaginatedResponseDTO[T any](content []T, page int64, pageSize int64, totalElements int64) PaginatedResponseDTO[T] {
	if content == nil {
		content = []T{}
	}
	totalPages := (totalElements + pageSize - 1) / pageSize
	return PaginatedResponseDTO[T]{
		Content: content,
		Paging: PagingDTO{
			TotalElements: totalElements,
			Size:          pageSize,
			Page:          page,
			TotalPages:    totalPages,
			HasContent:    len(content) > 0,
			HasNext:       page < totalPages-1,
			HasPrevious:   page > 0,
		},
	}
}

// MapPage converts the content and keeps the paging metadata.
func MapPage[T any, R any](page PaginatedResponseDTO[T], mapContent func(T) R) PaginatedResponseDTO[R] {
	content := make([]R, 0, len(page.Content))
	for _, item := range page.Content {
		content = append(content, mapContent(item))
	}
	return PaginatedResponseDTO[R]{Content: content, Paging: page.Paging}
}

// PageDoc only documents {content, paging} in Swagger annotations
// (`dto.PageDoc{content=[]dto.XxxResponseDTO}`); swag cannot resolve generic
// types across packages that share the name dto. Handlers return PaginatedResponseDTO.
type PageDoc struct {
	Content []any     `json:"content"`
	Paging  PagingDTO `json:"paging"`
}
