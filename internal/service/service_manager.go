package service

import (
	"api-monitoring-services/internal/domain"
	"api-monitoring-services/internal/repository"
	"context"
	"time"
)

type ServiceManager struct {
	repo          repository.ServiceRepository
	checkInterval time.Duration
	ctx           context.Context
	cancelFunc    context.CancelFunc
}

func NewServiceManager(repo repository.ServiceRepository, checkInterval time.Duration) *ServiceManager {
	ctx, cancel := context.WithCancel(context.Background())

	return &ServiceManager{
		repo:          repo,
		checkInterval: checkInterval,
		ctx:           ctx,
		cancelFunc:    cancel,
	}
}

func (sm *ServiceManager) GetAllServices() []domain.Service {
	return sm.repo.GetAll()
}
