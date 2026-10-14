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

func NewPrinterDataAdapter(db *mongo.Database) port.PrinterDataPort {
	collection := database.NewCollection[schema.PrinterDocument](db, schema.PrinterCollection, exception.NewPrinterNotFoundByID)
	return database.NewResourceAdapter(collection, database.ResourceConfig[dto.PrinterRequestDTO, dto.PrinterInputDTO, schema.PrinterDocument, dto.PrinterResponseDTO]{
		ToDocument:  mapper.ToPrinterDocument,
		ToResponse:  mapper.ToPrinterResponseDTO,
		SortColumns: schema.PrinterSortColumns,
		Filter: func(request dto.PrinterRequestDTO) bson.M {
			filter := bson.M{}
			if request.BranchID != nil {
				filter["branch_id"] = *request.BranchID
			}
			if request.DepartmentsID != nil {
				filter["departments_id"] = *request.DepartmentsID
			}
			return filter
		},
	})
}
