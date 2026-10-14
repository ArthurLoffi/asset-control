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

func NewEmailDataAdapter(db *mongo.Database) port.EmailDataPort {
	collection := database.NewCollection[schema.EmailDocument](db, schema.EmailCollection, exception.NewEmailNotFoundByID).WithInt32IDs()
	return database.NewResourceAdapter(collection, database.ResourceConfig[dto.EmailRequestDTO, dto.EmailInputDTO, schema.EmailDocument, dto.EmailResponseDTO]{
		ToDocument:  mapper.ToEmailDocument,
		ToResponse:  mapper.ToEmailResponseDTO,
		SortColumns: schema.EmailSortColumns,
		Filter: func(request dto.EmailRequestDTO) bson.M {
			filter := bson.M{}
			if request.User != nil {
				filter["user"] = database.Contains(*request.User)
			}
			return filter
		},
	})
}
