package ports

type UtilsPort interface {
	BranchIdToDescription(id uint64) (string, error)

	DepartmentIdToDescription(id uint64) (string, error)
}
