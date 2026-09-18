package dto

type BranchDto struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
}

func (BranchDto) TableName() string {
	return "branches"
}

type DepartmentDto struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
}

func (DepartmentDto) TableName() string {
	return "departments"
}
