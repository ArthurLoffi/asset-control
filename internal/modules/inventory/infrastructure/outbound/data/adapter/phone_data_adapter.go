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

func NewPhoneDataAdapter(db *mongo.Database) port.PhoneDataPort {
	collection := database.NewCollection[schema.PhoneDocument](db, schema.PhoneCollection, exception.NewPhoneNotFoundByID)
	return database.NewResourceAdapter(collection, database.ResourceConfig[dto.PhoneRequestDTO, dto.PhoneInputDTO, schema.PhoneDocument, dto.PhoneResponseDTO]{
		ToDocument:  mapper.ToPhoneDocument,
		ToResponse:  mapper.ToPhoneResponseDTO,
		SortColumns: schema.PhoneSortColumns,
		Filter: func(request dto.PhoneRequestDTO) bson.M {
			return assignmentFilter(request.AssignmentFilterDTO)
		},
	})
}
