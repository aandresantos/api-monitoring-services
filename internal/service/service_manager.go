package service

import (
	"api-monitoring-services/internal/domain"
	"api-monitoring-services/internal/repository"
	"context"
	"fmt"
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

func (sm *ServiceManager) CreateService(name, urlAddress string) (*domain.Service, error) {
	if name == "" {
		return nil, fmt.Errorf("nome do Serviço é obrigatório")
	}

	if urlAddress == "" {
		return nil, fmt.Errorf("endereço de URL é obrigatório para registrar o serviço")
	}

	service := domain.NewService(domain.NewServiceBody{
		Name:       name,
		URLAddress: urlAddress,
	})

	if err := sm.repo.Create(service); err != nil {
		return nil, fmt.Errorf("erro ao criar o serviço %v", err)
	}

	return service, nil
}
