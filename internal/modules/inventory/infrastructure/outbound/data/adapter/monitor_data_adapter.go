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

func NewMonitorDataAdapter(db *mongo.Database) port.MonitorDataPort {
	collection := database.NewCollection[schema.MonitorDocument](db, schema.MonitorCollection, exception.NewMonitorNotFoundByID)
	return database.NewResourceAdapter(collection, database.ResourceConfig[dto.MonitorRequestDTO, dto.MonitorInputDTO, schema.MonitorDocument, dto.MonitorResponseDTO]{
		ToDocument:  mapper.ToMonitorDocument,
		ToResponse:  mapper.ToMonitorResponseDTO,
		SortColumns: schema.MonitorSortColumns,
		Filter: func(request dto.MonitorRequestDTO) bson.M {
			return assignmentFilter(request.AssignmentFilterDTO)
		},
	})
}
