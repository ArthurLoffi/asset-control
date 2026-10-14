package dto

type BranchResponseDTO struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type DepartmentResponseDTO struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}
