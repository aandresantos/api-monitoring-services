package service

import (
	"api-monitoring-services/internal/domain"
	"api-monitoring-services/internal/pkg/healthcheck"
	"api-monitoring-services/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ServiceManager struct {
	repo          repository.IServiceRepository
	checkInterval time.Duration
	ctx           context.Context
	cancelFunc    context.CancelFunc
	healthChecker healthcheck.IHealthChecker
	wg            sync.WaitGroup
}

func NewServiceManager(repo repository.IServiceRepository, checker healthcheck.IHealthChecker, checkInterval time.Duration) *ServiceManager {
	ctx, cancel := context.WithCancel(context.Background())

	return &ServiceManager{
		repo:          repo,
		checkInterval: checkInterval,
		ctx:           ctx,
		cancelFunc:    cancel,
		healthChecker: checker,
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
	exist, err := sm.repo.Exists(id)

	if !exist {
		return fmt.Errorf("serviço com ID %s não encontrado", id)
	}

	if err != nil {
		return fmt.Errorf("erro ao encontrar serviço: %v", err)
	}

	if err := sm.repo.Delete(id); err != nil {
		return fmt.Errorf("erro ao deletar serviço: %v", err)
	}

	return nil
}

func (sm *ServiceManager) StopHealthChecks() {
	sm.cancelFunc()
	sm.wg.Wait()
}

func (sm *ServiceManager) StartHealthChecks() {
	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()

		ticker := time.NewTicker(sm.checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-sm.ctx.Done():
				return
			case <-ticker.C:
				sm.performHealthChecks()
			}
		}
	}()
}

func (sm *ServiceManager) performHealthChecks() {
	fmt.Println("[Health Check] -> Executando ciclo de verificação")

	services := sm.repo.GetAll()

	var wg sync.WaitGroup

	for _, service := range services {
		wg.Add(1)
		go func(svc domain.Service) {
			defer wg.Done()
			sm.checkSingleService(&svc)
		}(service)
	}

	wg.Wait()
	fmt.Printf("[Health Check] -> Verificação concluída para %d serviços ---\n", len(services))
}

func (sm *ServiceManager) checkSingleService(service *domain.Service) {
	checkResult := sm.healthChecker.Check(service.URLAddress)

	newCheck := &domain.HealthCheck{
		ID:             uuid.New().String(),
		ServiceID:      service.ID,
		Status:         checkResult.Status,
		ResponseTimeMs: int(checkResult.ResponseTime),
	}

	if checkResult.HTTPStatusCode != 0 {
		newCheck.HTTPStatusCode = sql.NullInt32{Int32: int32(checkResult.HTTPStatusCode), Valid: true}
	}

	if checkResult.Error != nil {
		newCheck.ErrorMessage = sql.NullString{String: string(checkResult.Error.Error()), Valid: true}
	}

	err := sm.repo.SaveCheckResult(context.Background(), newCheck)

	if err != nil {
		fmt.Printf("[Health Check] Erro ao salvar resultado para o serviço %s: %v", service.Name, err)
	} else {
		fmt.Printf("[Health Check] Serviço '%s' verificado. Status: %s.", service.Name, newCheck.Status, newCheck.ResponseTimeMs)
	}

}
