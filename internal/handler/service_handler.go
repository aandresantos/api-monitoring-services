package handler

import (
	"api-monitoring-services/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ServiceHandler struct {
	serviceManager *service.ServiceManager
}

func NewServiceHandler(sm *service.ServiceManager) *ServiceHandler {
	return &ServiceHandler{
		serviceManager: sm,
	}
}

func (h *ServiceHandler) ListServices(ctx echo.Context) error {
	services := h.serviceManager.GetAllServices()

	return ctx.JSON(http.StatusOK, services)
}

func (h *ServiceHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/services", h.ListServices)
}
