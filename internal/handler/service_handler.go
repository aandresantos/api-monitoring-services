package handler

import (
	"api-monitoring-services/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ServiceHandler struct {
	serviceManager *service.ServiceManager
}

type CreateServiceRequest struct {
	Name       string `json:"name" binding:"required"`
	URLAddress string `json:"urlAddress" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
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

func (h *ServiceHandler) CreateService(ctx echo.Context) error {
	var req CreateServiceRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "dados inválidos"})
	}

	if req.Name == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Nome é obrigatório"})
	}

	if req.URLAddress == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "URL é obrigatória"})
	}

	service, err := h.serviceManager.CreateService(req.Name, req.URLAddress)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, service)
}

func (h *ServiceHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/services", h.ListServices)
	e.POST("/services", h.CreateService)
}
