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

func (sm *ServiceManager) GetServiceByID(id string) (*domain.Service, error) {
	service, found := sm.repo.GetByID(id)

	if !found {
		return nil, fmt.Errorf("serviço com ID %s não encontrado", id)
	}

	return service, nil
}

func (sm *ServiceManager) UpdateService(id string, name, urlAddress *string) (*domain.Service, error) {
	service, found := sm.repo.GetByID(id)

	if !found {
		return nil, fmt.Errorf("serviço com ID %s não encontrado", id)
	}

	service.UpdateDetails(domain.EditServiceBody{
		Name:       name,
		URLAddress: urlAddress,
	})

	if err := sm.repo.Update(service); err != nil {
		return nil, fmt.Errorf("erro ao atualizar serviço: %v", err)
	}

	return service, nil
}

func (sm *ServiceManager) DeleteService(id string) error {
	if !sm.repo.Exists(id) {
		return fmt.Errorf("serviço com ID %s não encontrado", id)
	}

	if err := sm.repo.Delete(id); err != nil {
		return fmt.Errorf("erro ao deletar serviço: %v", err)
	}

	return nil
}
