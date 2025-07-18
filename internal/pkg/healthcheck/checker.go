package healthcheck

import (
	"api-monitoring-services/internal/domain"
	"fmt"
	"net/http"
	"time"
)

type CheckResult struct {
	Status         domain.ServiceStatus
	ResponseTime   time.Duration
	HTTPStatusCode int
	Error          error
}

type HTTPHealthChecker struct {
	client http.Client
}

type IHealthChecker interface {
	Check(url string) CheckResult
}

type IChecker interface {
	Check(url string) CheckResult
}

func NewHTTPHealthChecker(timeout time.Duration) IChecker {
	return &HTTPHealthChecker{
		client: http.Client{Timeout: timeout},
	}
}

func (h *HTTPHealthChecker) Check(url string) CheckResult {
	startTime := time.Now()
	resp, err := h.client.Get(url)
	responseTime := time.Since(startTime)

	if err != nil {
		return CheckResult{
			Status:       domain.StatusOffline,
			ResponseTime: responseTime,
			Error:        err,
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return CheckResult{
			Status:         domain.StatusOnline,
			ResponseTime:   responseTime,
			HTTPStatusCode: resp.StatusCode,
		}
	}

	return CheckResult{
		Status:         domain.StatusOffline,
		ResponseTime:   responseTime,
		HTTPStatusCode: resp.StatusCode,
		Error:          fmt.Errorf("status code: %d", resp.StatusCode),
	}
}
