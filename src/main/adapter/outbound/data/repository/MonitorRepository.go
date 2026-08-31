package repository

import (
	"asset-control/src/main/adapter/outbound/data/entities"
	"asset-control/src/main/domain/dto"
	"errors"

	"gorm.io/gorm"
)

type MonitorRepository struct {
	db *gorm.DB
}

func NewMonitorRepository(db *gorm.DB) *MonitorRepository {
	return &MonitorRepository{db: db}
}

func (mr *MonitorRepository) CreateMonitor(monitor dto.CreateMonitorDto) error {
	monitorEntity := entities.Monitor{
		HasMonitor:       monitor.HasMonitor,
		AssetNumber:      monitor.AssetNumber,
		Brand:            monitor.Brand,
		Model:            monitor.Model,
		SerialNumber:     monitor.SerialNumber,
		MonitorCondition: monitor.MonitorCondition,
		Situation:        monitor.Situation,
		EmployeeID:       monitor.EmployeeID,
	}

	if err := mr.db.Create(&monitorEntity).Error; err != nil {
		return err
	}

	return nil
}

func (mr *MonitorRepository) UpdateMonitor(
	id uint64,
	monitor dto.CreateMonitorDto,
) (dto.MonitorDto, error) {
	var monitorEntity entities.Monitor

	if err := mr.db.First(&monitorEntity, id).Error; err != nil {
		return dto.MonitorDto{}, err
	}

	monitorEntity.HasMonitor = monitor.HasMonitor
	monitorEntity.AssetNumber = monitor.AssetNumber
	monitorEntity.Brand = monitor.Brand
	monitorEntity.Model = monitor.Model
	monitorEntity.SerialNumber = monitor.SerialNumber
	monitorEntity.MonitorCondition = monitor.MonitorCondition
	monitorEntity.Situation = monitor.Situation
	monitorEntity.EmployeeID = monitor.EmployeeID

	if err := mr.db.Save(&monitorEntity).Error; err != nil {
		return dto.MonitorDto{}, err
	}

	return dto.MonitorDto{
		ID:               monitorEntity.ID,
		HasMonitor:       monitorEntity.HasMonitor,
		AssetNumber:      monitorEntity.AssetNumber,
		Brand:            monitorEntity.Brand,
		Model:            monitorEntity.Model,
		SerialNumber:     monitorEntity.SerialNumber,
		MonitorCondition: monitorEntity.MonitorCondition,
		Situation:        monitorEntity.Situation,
		EmployeeID:       monitorEntity.EmployeeID,
		CreatedAt:        monitorEntity.CreatedAt,
	}, nil
}

func (mr *MonitorRepository) GetAllMonitors() ([]dto.MonitorDto, error) {
	var monitorEntities []entities.Monitor

	if err := mr.db.Find(&monitorEntities).Error; err != nil {
		return nil, err
	}

	monitors := make([]dto.MonitorDto, 0, len(monitorEntities))

	for _, monitorEntity := range monitorEntities {
		monitors = append(monitors, dto.MonitorDto{
			ID:               monitorEntity.ID,
			HasMonitor:       monitorEntity.HasMonitor,
			AssetNumber:      monitorEntity.AssetNumber,
			Brand:            monitorEntity.Brand,
			Model:            monitorEntity.Model,
			SerialNumber:     monitorEntity.SerialNumber,
			MonitorCondition: monitorEntity.MonitorCondition,
			Situation:        monitorEntity.Situation,
			EmployeeID:       monitorEntity.EmployeeID,
			CreatedAt:        monitorEntity.CreatedAt,
		})
	}

	return monitors, nil
}

func (mr *MonitorRepository) FindByID(id uint64) (dto.MonitorDto, error) {
	var monitorEntity entities.Monitor

	if err := mr.db.First(&monitorEntity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.MonitorDto{}, err
		}

		return dto.MonitorDto{}, err
	}

	return dto.MonitorDto{
		ID:               monitorEntity.ID,
		HasMonitor:       monitorEntity.HasMonitor,
		AssetNumber:      monitorEntity.AssetNumber,
		Brand:            monitorEntity.Brand,
		Model:            monitorEntity.Model,
		SerialNumber:     monitorEntity.SerialNumber,
		MonitorCondition: monitorEntity.MonitorCondition,
		Situation:        monitorEntity.Situation,
		EmployeeID:       monitorEntity.EmployeeID,
		CreatedAt:        monitorEntity.CreatedAt,
	}, nil
}

func (mr *MonitorRepository) FindByEmployeeID(employeeID uint64) ([]dto.MonitorDto, error) {
	var monitorEntities []entities.Monitor

	if err := mr.db.
		Where("employee_id = ?", employeeID).
		Find(&monitorEntities).Error; err != nil {
		return nil, err
	}

	monitors := make([]dto.MonitorDto, 0, len(monitorEntities))

	for _, monitorEntity := range monitorEntities {
		monitors = append(monitors, dto.MonitorDto{
			ID:               monitorEntity.ID,
			HasMonitor:       monitorEntity.HasMonitor,
			AssetNumber:      monitorEntity.AssetNumber,
			Brand:            monitorEntity.Brand,
			Model:            monitorEntity.Model,
			SerialNumber:     monitorEntity.SerialNumber,
			MonitorCondition: monitorEntity.MonitorCondition,
			Situation:        monitorEntity.Situation,
			EmployeeID:       monitorEntity.EmployeeID,
			CreatedAt:        monitorEntity.CreatedAt,
		})
	}

	return monitors, nil
}
