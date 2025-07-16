package repository

import (
	"api-monitoring-services/internal/domain"
	"sync"
)

type ServiceRepository interface {
	GetAll() []domain.Service
}

type InMemoryServiceRepository struct {
	store map[string]*domain.Service
	mutex *sync.RWMutex
}

func NewInMemoryServiceRepository() *InMemoryServiceRepository {

	return &InMemoryServiceRepository{
		store: make(map[string]*domain.Service),
		mutex: &sync.RWMutex{},
	}
}

func (r *InMemoryServiceRepository) GetAll() []domain.Service {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	services := make([]domain.Service, 0, len(r.store))

	for _, service := range r.store {
		services = append(services, *service)
	}

	return services
}
