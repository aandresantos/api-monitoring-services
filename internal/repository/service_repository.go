package repository

import (
	"api-monitoring-services/internal/database"
	"api-monitoring-services/internal/domain"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type IServiceRepository interface {
	GetAll() []domain.Service
	Create(svc *domain.Service) error
	GetByID(id string) (*domain.Service, bool)
	Update(service *domain.Service) error
	Delete(id string) error
	Exists(id string) (bool, error)
	SaveCheckResult(ctx context.Context, check *domain.HealthCheck) error
}

type ServiceRepository struct {
	db *database.DBClient
}

func NewServiceRepository(dbClient *database.DBClient) *ServiceRepository {

	return &ServiceRepository{
		db: dbClient,
	}
}

func (r *ServiceRepository) GetAll() []domain.Service {
	rows, err := r.db.Conn.Query(
		context.Background(),
		`SELECT id, name, url_address, status, created_at, updated_at
		FROM services`,
	)

	if err != nil {
		fmt.Printf("Erro ao buscar serviços: %v\n", err)
		return []domain.Service{}
	}
	defer rows.Close()

	var services []domain.Service

	for rows.Next() {
		var svc domain.Service
		err := rows.Scan(
			&svc.ID,
			&svc.Name,
			&svc.URLAddress,
			&svc.Status,
			&svc.CreatedAt,
			&svc.UpdatedAt,
		)
		if err != nil {
			fmt.Printf("Erro ao ler linha: %v\n", err)
			continue
		}
		services = append(services, svc)
	}

	return services
}

func (r *ServiceRepository) GetByID(id string) (*domain.Service, bool) {
	query := `SELECT id, name, url_address, status, created_at, updated_at FROM services WHERE id = $1`

	row := r.db.Conn.QueryRow(context.Background(), query, id)

	var svc domain.Service
	err := row.Scan(&svc.ID, &svc.Name, &svc.URLAddress, &svc.Status, &svc.CreatedAt, &svc.UpdatedAt)

	if err != nil {
		return nil, false
	}

	return &svc, true
}

func (r *ServiceRepository) Create(svc *domain.Service) error {
	query := `
		INSERT INTO services (id, name, url_address, status, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Conn.Exec(
		context.Background(),
		query,
		svc.ID,
		svc.Name,
		svc.URLAddress,
		svc.Status,
		svc.UpdatedAt,
	)

	return err
}

func (r *ServiceRepository) Update(service *domain.Service) error {
	query := `
		UPDATE services
		SET name = $1, url_address = $2, status = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.Conn.Exec(
		context.Background(),
		query,
		service.Name,
		service.URLAddress,
		service.Status,
		time.Now(),
		service.ID,
	)

	return err
}

func (r *ServiceRepository) Delete(id string) error {
	res, err := r.db.Conn.Exec(context.Background(), "DELETE FROM services WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("erro ao deletar serviço com id %s: %w", id, err)
	}

	rowsAffected := res.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum serviço encontrado com o id %s", id)
	}

	return nil
}

func (r *ServiceRepository) Exists(id string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM services WHERE id = $1)`
	err := r.db.Conn.QueryRow(context.Background(), query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência do serviço com id %s: %w", id, err)
	}
	return exists, nil
}

func (r *ServiceRepository) SaveCheckResult(ctx context.Context, check *domain.HealthCheck) error {
	tx, err := r.db.Conn.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	checkQuery := `
		INSERT INTO health_checks (id, service_id, status, response_time_ms, http_status_code, error_message)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = tx.Exec(ctx, checkQuery, check.ID, check.ServiceID, check.Status, check.ResponseTimeMs, check.HTTPStatusCode, check.ErrorMessage)
	if err != nil {
		return err
	}

	updateQuery := `UPDATE services SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err = tx.Exec(ctx, updateQuery, check.Status, check.ServiceID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
