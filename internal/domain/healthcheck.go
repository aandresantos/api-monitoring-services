package domain

import (
	"database/sql"
	"time"
)

type HealthCheck struct {
	ID             string         `json:"id"`
	ServiceID      string         `json:"serviceId"`
	Status         ServiceStatus  `json:"status"`
	ResponseTimeMs int            `json:"responseTimeMs"`
	HTTPStatusCode sql.NullInt32  `json:"httpStatusCode,omitempty"`
	ErrorMessage   sql.NullString `json:"errorMessage,omitempty"`
	CheckedAt      time.Time      `json:"checkedAt"`
}
