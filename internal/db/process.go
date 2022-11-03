package db

import (
	"errors"
	"fmt"
	"github.com/kamkali/go-timeline/internal/db/schema/models"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/logger"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type ProcessRepository struct {
	log logger.Logger

	db *gorm.DB
}

func NewProcessRepository(log logger.Logger, db *gorm.DB) *ProcessRepository {
	return &ProcessRepository{log: log, db: db}
}

func toDBProcess(de *domain.Process) (*models.Process, error) {
	return &models.Process{
		Name:                de.Name,
		StartTime:           de.StartTime,
		EndTime:             de.EndTime,
		ShortDescription:    de.ShortDescription,
		DetailedDescription: de.DetailedDescription,
		Graphic:             de.Graphic,
		TypeID:              de.TypeID,
	}, nil
}

func (t ProcessRepository) GetProcess(ctx context.Context, id uint) (domain.Process, error) {
	var process models.Process
	if err := t.db.WithContext(ctx).First(&process, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Process{}, domain.ErrNotFound
		}
		return domain.Process{}, fmt.Errorf("db error on select query: %w", err)
	}
	domainProcess, err := toDomainProcess(process)
	if err != nil {
		return domain.Process{}, fmt.Errorf("cannot translate db model to domain")
	}
	return domainProcess, nil
}

func (t ProcessRepository) UpdateProcess(ctx context.Context, id uint, process *domain.Process) error {
	var e models.Process
	r := t.db.WithContext(ctx).Find(&e, id)
	if r.Error != nil {
		return fmt.Errorf("db error on select query: %w", r.Error)
	}

	e.Name = process.Name
	e.StartTime = process.StartTime
	e.EndTime = process.EndTime
	e.ShortDescription = process.ShortDescription
	e.DetailedDescription = process.DetailedDescription
	e.Graphic = process.Graphic
	e.TypeID = process.TypeID

	if err := t.db.WithContext(ctx).Save(&e).Error; err != nil {
		return fmt.Errorf("db error on update query: %w", r.Error)
	}

	return nil
}

func (t ProcessRepository) DeleteProcess(ctx context.Context, id uint) error {
	if err := t.db.WithContext(ctx).Delete(&models.Process{}, id).Error; err != nil {
		return fmt.Errorf("error while deleting: %w", err)
	}
	return nil
}

func (t ProcessRepository) CreateProcess(ctx context.Context, process *domain.Process) (uint, error) {
	dbProcess, err := toDBProcess(process)
	if err != nil {
		return 0, err
	}

	typ := models.Type{}
	result := t.db.Table("types").Find(&typ, process.TypeID)
	if result.Error != nil || result.RowsAffected == 0 {
		return 0, fmt.Errorf("cannot find type of name %s", process.TypeID)
	}
	dbProcess.TypeID = typ.ID

	result = t.db.WithContext(ctx).Create(dbProcess)
	if result.Error != nil {
		return 0, fmt.Errorf("cannot create process: %w", result.Error)
	}
	return dbProcess.ID, nil
}

func (t ProcessRepository) ListProcesses(ctx context.Context) ([]domain.Process, error) {
	var processes []models.Process
	r := t.db.WithContext(ctx).Find(&processes)
	if r.Error != nil {
		return nil, fmt.Errorf("db error on select query: %w", r.Error)
	}

	domainProcesses := []domain.Process{}
	for _, e := range processes {
		domainProcess, err := toDomainProcess(e)
		if err != nil {
			return nil, fmt.Errorf("cannot translate db model to domain")
		}
		domainProcesses = append(domainProcesses, domainProcess)
	}

	return domainProcesses, nil
}

func toDomainProcess(e models.Process) (domain.Process, error) {
	domainProcess := domain.Process{
		ID:                  e.ID,
		Name:                e.Name,
		StartTime:           e.StartTime,
		EndTime:             e.EndTime,
		ShortDescription:    e.ShortDescription,
		DetailedDescription: e.DetailedDescription,
		Graphic:             e.Graphic,
		TypeID:              e.TypeID,
	}
	return domainProcess, nil
}
