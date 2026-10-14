package exception

import "asset-control/internal/core/exception"

const (
	BranchNotFoundByIDCode     = "BRANCH_001"
	DepartmentNotFoundByIDCode = "DEPARTMENT_001"
)

func NewBranchNotFoundByID(id uint64) error {
	return exception.NotFoundByID("Branch", BranchNotFoundByIDCode, id)
}

func NewDepartmentNotFoundByID(id uint64) error {
	return exception.NotFoundByID("Department", DepartmentNotFoundByIDCode, id)
}
