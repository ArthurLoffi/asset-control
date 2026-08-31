package repository

import (
	"asset-control/src/main/adapter/outbound/data/entities"
	"asset-control/src/main/domain/dto"
	"errors"

	"gorm.io/gorm"
)

type PrinterRepository struct {
	db *gorm.DB
}

func NewPrinterRepository(db *gorm.DB) *PrinterRepository {
	return &PrinterRepository{db: db}
}

func (pr *PrinterRepository) CreatePrinter(printer dto.CreatePrinterDto) error {
	printerEntity := entities.Printer{
		AssetNumber:      printer.AssetNumber,
		Model:            printer.Model,
		IP:               printer.IP,
		MachineCondition: printer.MachineCondition,
		BranchID:         printer.BranchID,
		DepartmentsID:    printer.DepartmentsID,
	}

	if err := pr.db.Create(&printerEntity).Error; err != nil {
		return err
	}

	return nil
}

func (pr *PrinterRepository) UpdatePrinter(
	id uint64,
	printer dto.CreatePrinterDto,
) (dto.PrinterDto, error) {
	var printerEntity entities.Printer

	if err := pr.db.First(&printerEntity, id).Error; err != nil {
		return dto.PrinterDto{}, err
	}

	printerEntity.AssetNumber = printer.AssetNumber
	printerEntity.Model = printer.Model
	printerEntity.IP = printer.IP
	printerEntity.MachineCondition = printer.MachineCondition
	printerEntity.BranchID = printer.BranchID
	printerEntity.DepartmentsID = printer.DepartmentsID

	if err := pr.db.Save(&printerEntity).Error; err != nil {
		return dto.PrinterDto{}, err
	}

	return dto.PrinterDto{
		ID:               printerEntity.ID,
		AssetNumber:      printerEntity.AssetNumber,
		Model:            printerEntity.Model,
		IP:               printerEntity.IP,
		MachineCondition: printerEntity.MachineCondition,
		BranchID:         printerEntity.BranchID,
		DepartmentsID:    printerEntity.DepartmentsID,
		CreatedAt:        printerEntity.CreatedAt,
	}, nil
}

func (pr *PrinterRepository) GetAllPrinters() ([]dto.PrinterDto, error) {
	var printerEntities []entities.Printer

	if err := pr.db.Find(&printerEntities).Error; err != nil {
		return nil, err
	}

	printers := make([]dto.PrinterDto, 0, len(printerEntities))

	for _, printerEntity := range printerEntities {
		printers = append(printers, dto.PrinterDto{
			ID:               printerEntity.ID,
			AssetNumber:      printerEntity.AssetNumber,
			Model:            printerEntity.Model,
			IP:               printerEntity.IP,
			MachineCondition: printerEntity.MachineCondition,
			BranchID:         printerEntity.BranchID,
			DepartmentsID:    printerEntity.DepartmentsID,
			CreatedAt:        printerEntity.CreatedAt,
		})
	}

	return printers, nil
}

func (pr *PrinterRepository) FindByID(id uint64) (dto.PrinterDto, error) {
	var printerEntity entities.Printer

	if err := pr.db.First(&printerEntity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.PrinterDto{}, err
		}

		return dto.PrinterDto{}, err
	}

	return dto.PrinterDto{
		ID:               printerEntity.ID,
		AssetNumber:      printerEntity.AssetNumber,
		Model:            printerEntity.Model,
		IP:               printerEntity.IP,
		MachineCondition: printerEntity.MachineCondition,
		BranchID:         printerEntity.BranchID,
		DepartmentsID:    printerEntity.DepartmentsID,
		CreatedAt:        printerEntity.CreatedAt,
	}, nil
}
