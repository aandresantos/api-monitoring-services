package repository

import (
	"api-monitoring-services/internal/database"
	"api-monitoring-services/internal/domain"
	"sync"
)

type ServiceRepository interface {
	GetAll() []domain.Service
	Create(svc *domain.Service) error
	GetByID(id string) (*domain.Service, bool)
	Update(service *domain.Service) error
	Delete(id string) error
	Exists(id string) bool
}

type InMemoryServiceRepository struct {
	store map[string]*domain.Service
	mutex *sync.RWMutex
	db    *database.DBClient
}

func NewInMemoryServiceRepository(dbClient *database.DBClient) *InMemoryServiceRepository {

	return &InMemoryServiceRepository{
		db:    dbClient,
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

func (r *InMemoryServiceRepository) GetByID(id string) (*domain.Service, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	service, found := r.store[id]

	if !found {
		return nil, false
	}

	serviceCopy := *service

	return &serviceCopy, true
}

func (r *InMemoryServiceRepository) Create(svc *domain.Service) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.store[svc.ID] = svc

	return nil
}

func (r *InMemoryServiceRepository) Update(service *domain.Service) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, found := r.store[service.ID]; !found {
		return nil
	}

	r.store[service.ID] = service
	return nil
}

func (r *InMemoryServiceRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, found := r.store[id]; !found {
		return nil
	}

	delete(r.store, id)
	return nil
}

func (r *InMemoryServiceRepository) Exists(id string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, found := r.store[id]
	return found
}
