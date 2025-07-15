package main

import (
	"sync"

	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Service struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	URLAddress string `json:"urlAddress"`
	Status     string `json:"status"`
}

var serviceStore = make(map[string]Service)
var storeMutex = &sync.Mutex{}

func main() {
	Echo := echo.New()

	Echo.GET("/services", listServices)
	Echo.POST("/services", addService)

	Echo.GET("/services/:id", getServiceByID)
	Echo.PATCH("/services/:id", patchService)
	Echo.DELETE("/services/:id", deleteService)

	Echo.Logger.Fatal(Echo.Start(":3000"))
}

func listServices(ctx echo.Context) error {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	var services []Service

	for _, service := range serviceStore {

		services = append(services, service)
	}

	return ctx.JSON(http.StatusOK, services)
}

func addService(ctx echo.Context) error {
	var serv Service

	if err := ctx.Bind(&serv); err != nil {
		return ctx.String(http.StatusBadRequest, "Erro ao adicionar Serviço")
	}

	serv.ID = uuid.New().String()
	serv.Status = "Pending"

	storeMutex.Lock()
	serviceStore[serv.ID] = serv
	storeMutex.Unlock()

	return ctx.JSON(http.StatusCreated, serv)
}

func getServiceByID(ctx echo.Context) error {
	id := ctx.Param("id")

	storeMutex.Lock()
	defer storeMutex.Unlock()

	service, found := serviceStore[id]

	if !found {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Serviço não encontrado"})
	}

	return ctx.JSON(http.StatusOK, service)
}

func patchService(ctx echo.Context) error {
	id := ctx.Param("id")

	type PatchServiceRequest struct {
		Name       *string `string:"name"`
		URLAddress *string `json:"urlAddress"`
	}

	var requestData PatchServiceRequest

	if err := ctx.Bind(&requestData); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Body inválido"})
	}

	storeMutex.Lock()
	defer storeMutex.Unlock()

	service, found := serviceStore[id]

	if !found {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Serviço não encontrado"})
	}

	if requestData.Name != nil {
		service.Name = *requestData.Name
	}

	if requestData.URLAddress != nil {
		service.URLAddress = *requestData.URLAddress
	}

	serviceStore[id] = service

	return ctx.JSON(http.StatusOK, service)
}

func deleteService(ctx echo.Context) error {
	id := ctx.Param("id")

	storeMutex.Lock()
	defer storeMutex.Unlock()

	_, found := serviceStore[id]

	if !found {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Serviço não encontrado"})
	}

	delete(serviceStore, id)

	return ctx.NoContent(http.StatusNoContent)
}
