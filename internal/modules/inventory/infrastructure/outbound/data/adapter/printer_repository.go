package adapter

import (
	"context"

	"asset-control/internal/modules/inventory/application/port"
	"asset-control/internal/modules/inventory/domain/dto"
	"asset-control/internal/modules/inventory/infrastructure/outbound/data/schema"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PrinterRepository struct {
	documents documentRepository[schema.Printer]
}

var _ port.PrinterPorts = (*PrinterRepository)(nil)

func NewPrinterRepository(db *mongo.Database) *PrinterRepository {
	return &PrinterRepository{documents: newDocumentRepository[schema.Printer](db, (schema.Printer{}).CollectionName())}
}

func printerDocument(input dto.CreatePrinterDto) schema.Printer {
	return schema.Printer{
		AssetNumber:      input.AssetNumber,
		Model:            input.Model,
		IP:               input.IP,
		MachineCondition: input.MachineCondition,
		BranchID:         input.BranchID,
		DepartmentsID:    input.DepartmentsID,
	}
}

func (r *PrinterRepository) CreatePrinter(ctx context.Context, input dto.CreatePrinterDto) (dto.PrinterDto, error) {
	document, err := r.documents.create(ctx, printerDocument(input))
	return dto.PrinterDto(document), err
}

func (r *PrinterRepository) UpdatePrinter(ctx context.Context, id uint64, input dto.CreatePrinterDto) (dto.PrinterDto, error) {
	document, err := r.documents.update(ctx, id, printerDocument(input))
	return dto.PrinterDto(document), err
}

func (r *PrinterRepository) GetAllPrinters(ctx context.Context) ([]dto.PrinterDto, error) {
	documents, err := r.documents.find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	result := make([]dto.PrinterDto, 0, len(documents))
	for _, document := range documents {
		result = append(result, dto.PrinterDto(document))
	}
	return result, nil
}

func (r *PrinterRepository) FindByID(ctx context.Context, id uint64) (dto.PrinterDto, error) {
	document, err := r.documents.findOne(ctx, bson.M{"id": id})
	return dto.PrinterDto(document), err
}
