package adapter

import (
	"asset-control/internal/core/database"
	"asset-control/internal/modules/organization/application/port"
	"asset-control/internal/modules/organization/domain/dto"
	"asset-control/internal/modules/organization/domain/exception"
	"asset-control/internal/modules/organization/infrastructure/outbound/data/mapper"
	"asset-control/internal/modules/organization/infrastructure/outbound/data/schema"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewBranchDataAdapter(db *mongo.Database) port.BranchDataPort {
	collection := database.NewCollection[schema.BranchDocument](db, schema.BranchCollection, exception.NewBranchNotFoundByID)
	return database.NewResourceAdapter(collection, database.ResourceConfig[dto.BranchRequestDTO, struct{}, schema.BranchDocument, dto.BranchResponseDTO]{
		ToResponse:  mapper.ToBranchResponseDTO,
		SortColumns: schema.OrganizationSortColumns,
		Filter:      func(request dto.BranchRequestDTO) bson.M { return nameFilter(request.Name) },
	})
}

func NewDepartmentDataAdapter(db *mongo.Database) port.DepartmentDataPort {
	collection := database.NewCollection[schema.DepartmentDocument](db, schema.DepartmentCollection, exception.NewDepartmentNotFoundByID)
	return database.NewResourceAdapter(collection, database.ResourceConfig[dto.DepartmentRequestDTO, struct{}, schema.DepartmentDocument, dto.DepartmentResponseDTO]{
		ToResponse:  mapper.ToDepartmentResponseDTO,
		SortColumns: schema.OrganizationSortColumns,
		Filter:      func(request dto.DepartmentRequestDTO) bson.M { return nameFilter(request.Name) },
	})
}

func nameFilter(name *string) bson.M {
	if name == nil {
		return bson.M{}
	}
	return bson.M{"name": database.Contains(*name)}
}
