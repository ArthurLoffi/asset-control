package repository

import (
	"context"

	"asset-control/src/main/application/ports"
	"asset-control/src/main/domain/dto"
	"asset-control/src/main/infrastructure/outbound/data/entities"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PrinterRepository struct {
	documents documentRepository[entities.Printer]
}

var _ ports.PrinterPorts = (*PrinterRepository)(nil)

func NewPrinterRepository(db *mongo.Database) *PrinterRepository {
	return &PrinterRepository{documents: newDocumentRepository[entities.Printer](db, (entities.Printer{}).CollectionName())}
}

func printerDocument(input dto.CreatePrinterDto) entities.Printer {
	return entities.Printer{
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
