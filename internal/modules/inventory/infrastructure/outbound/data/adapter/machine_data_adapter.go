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

func NewMachineDataAdapter(db *mongo.Database) port.MachineDataPort {
	collection := database.NewCollection[schema.MachineDocument](db, schema.MachineCollection, exception.NewMachineNotFoundByID)
	return database.NewResourceAdapter(collection, database.ResourceConfig[dto.MachineRequestDTO, dto.MachineInputDTO, schema.MachineDocument, dto.MachineResponseDTO]{
		ToDocument:  mapper.ToMachineDocument,
		ToResponse:  mapper.ToMachineResponseDTO,
		SortColumns: schema.MachineSortColumns,
		Filter: func(request dto.MachineRequestDTO) bson.M {
			return assignmentFilter(request.AssignmentFilterDTO)
		},
	})
}
