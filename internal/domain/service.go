package domain

import (
	"time"

	"github.com/google/uuid"
)

type ServiceStatus string

const (
	StatusPeding  ServiceStatus = "pending"
	StatusOnline  ServiceStatus = "online"
	StatusOffline ServiceStatus = "offline"
)

type Service struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	URLAddress string        `json:"urlAddress"`
	Status     ServiceStatus `json:"status"`
	CreatedAt  time.Time     `json:"createdAt"`
	UpdatedAt  time.Time     `json:"updatedAt"`
}

type NewServiceBody struct {
	Name       string `json:"name"`
	URLAddress string `json:"urlAddress"`
}

func NewService(body NewServiceBody) *Service {
	dateNow := time.Now()

	return &Service{
		ID:         uuid.New().String(),
		Name:       body.Name,
		URLAddress: body.URLAddress,
		Status:     StatusPeding,
		CreatedAt:  dateNow,
		UpdatedAt:  dateNow,
	}
}
