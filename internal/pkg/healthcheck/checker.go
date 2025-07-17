package healthcheck

import (
	"api-monitoring-services/internal/domain"
	"context"
	"net/http"
	"time"
)

type HealthChecker interface {
	Check(svc domain.Service) domain.ServiceStatus
	CheckWithMetrics(svc *domain.Service) *CheckResult
}

type HTTPHealthChecker struct {
	client  *http.Client
	timeout time.Duration
}

type CheckResult struct {
	ServiceID string
	Status    domain.ServiceStatus
	CheckedAt time.Time
	Duration  time.Duration
}

func NewHTTPHealthChecker(timeout time.Duration) *HTTPHealthChecker {
	return &HTTPHealthChecker{
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (h HTTPHealthChecker) Check(svc *domain.Service) domain.ServiceStatus {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, svc.URLAddress, nil)

	if err != nil {
		return domain.StatusOffline
	}

	resp, err := h.client.Do(req)

	if err != nil {
		return domain.StatusOffline

	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return domain.StatusOffline
	}

	return domain.StatusOffline
}

func (h *HTTPHealthChecker) CheckWithMetrics(svc *domain.Service) *CheckResult {

	start := time.Now()
	status := h.Check(svc)
	duration := time.Since(start)

	return &CheckResult{
		ServiceID: svc.ID,
		Status:    status,
		CheckedAt: start,
		Duration:  duration,
	}
}
