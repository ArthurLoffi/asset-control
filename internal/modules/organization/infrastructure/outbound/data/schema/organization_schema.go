package schema

const (
	BranchCollection     = "branches"
	DepartmentCollection = "departments"
)

// OrganizationSortColumns maps the public sort_by value to the stored field.
var OrganizationSortColumns = map[string]string{"id": "id", "name": "name"}

type BranchDocument struct {
	ID   uint64 `bson:"id"`
	Name string `bson:"name"`
}

type DepartmentDocument struct {
	ID   uint64 `bson:"id"`
	Name string `bson:"name"`
}
