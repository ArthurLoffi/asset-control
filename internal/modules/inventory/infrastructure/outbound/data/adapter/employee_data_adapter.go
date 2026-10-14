package adapter

import (
	"asset-control/internal/core/database"
	"asset-control/internal/modules/inventory/application/port"
	"asset-control/internal/modules/inventory/domain/dto"
	"asset-control/internal/modules/inventory/domain/exception"
	"asset-control/internal/modules/inventory/infrastructure/outbound/data/mapper"
	"asset-control/internal/modules/inventory/infrastructure/outbound/data/schema"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewEmployeeDataAdapter(db *mongo.Database) port.EmployeeDataPort {
	collection := database.NewCollection[schema.EmployeeDocument](db, schema.EmployeeCollection, exception.NewEmployeeNotFoundByID)
	return database.NewResourceAdapter(collection, database.ResourceConfig[dto.EmployeeRequestDTO, dto.EmployeeInputDTO, schema.EmployeeDocument, dto.EmployeeResponseDTO]{
		ToDocument:  mapper.ToEmployeeDocument,
		ToResponse:  mapper.ToEmployeeResponseDTO,
		SortColumns: schema.EmployeeSortColumns,
		Filter: func(request dto.EmployeeRequestDTO) bson.M {
			filter := bson.M{}
			if request.Name != nil {
				filter["name"] = database.Contains(*request.Name)
			}
			if request.BranchID != nil {
				filter["branch_id"] = *request.BranchID
			}
			if request.DepartmentID != nil {
				filter["department_id"] = *request.DepartmentID
			}
			return filter
		},
		AfterDeactivate: unassignEmployee(db),
	})
}
