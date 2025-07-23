package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"api-monitoring-services/internal/service"
)

type ServiceHandler struct {
	serviceManager *service.ServiceManager
}

type CreateServiceRequest struct {
	Name       string `json:"name" binding:"required"`
	URLAddress string `json:"urlAddress" binding:"required"`
}

type UpdateServiceRequest struct {
	Name       *string `json:"name"`
	URLAddress *string `json:"urlAddress"`
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

func (h *ServiceHandler) GetServiceByID(ctx echo.Context) error {
	id := ctx.Param("id")

	if id == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Defina um ID para buscar o Serviço"})
	}

	service, err := h.serviceManager.GetServiceByID(id)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(http.StatusOK, service)
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

func (h *ServiceHandler) UpdateService(c echo.Context) error {
	id := c.Param("id")

	var req UpdateServiceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Dados inválidos"})
	}

	service, err := h.serviceManager.UpdateService(id, req.Name, req.URLAddress)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, service)
}

func (h *ServiceHandler) DeleteService(c echo.Context) error {
	id := c.Param("id")

	if err := h.serviceManager.DeleteService(id); err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ServiceHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/api/services", h.ListServices)
	e.POST("/api/services", h.CreateService)
	e.GET("/api/services/:id", h.GetServiceByID)
	e.PATCH("/api/services/:id", h.UpdateService)
	e.DELETE("/api/services/:id", h.DeleteService)
}
