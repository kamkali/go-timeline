package processservice

import (
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/logger"
	"golang.org/x/net/context"
)

type ProcessService struct {
	log logger.Logger

	repo domain.ProcessRepository
}

func (t ProcessService) GetProcess(ctx context.Context, id uint) (domain.Process, error) {
	return t.repo.GetProcess(ctx, id)
}

func (t ProcessService) UpdateProcess(ctx context.Context, id uint, process *domain.Process) error {
	return t.repo.UpdateProcess(ctx, id, process)
}

func (t ProcessService) DeleteProcess(ctx context.Context, id uint) error {
	return t.repo.DeleteProcess(ctx, id)
}

func (t ProcessService) CreateProcess(ctx context.Context, process *domain.Process) (uint, error) {
	return t.repo.CreateProcess(ctx, process)
}

func (t ProcessService) ListProcesses(ctx context.Context) ([]domain.Process, error) {
	return t.repo.ListProcesses(ctx)
}

func New(log logger.Logger, repo domain.ProcessRepository) *ProcessService {
	return &ProcessService{log: log, repo: repo}
}
