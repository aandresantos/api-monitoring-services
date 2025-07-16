package main

import (
	"fmt"
	"sync"
	"time"

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

	startHealthChecks()

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
	serv.Status = "pending"

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

func checkService(svc Service) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(svc.URLAddress)

	var newStatus string

	if err != nil || resp.StatusCode >= 400 {
		newStatus = "offline"
	} else {
		newStatus = "online"

		defer resp.Body.Close()
	}

	storeMutex.Lock()
	defer storeMutex.Unlock()

	if service, found := serviceStore[svc.ID]; found {
		service.Status = newStatus
		serviceStore[svc.ID] = service
	}
}

func startHealthChecks() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			<-ticker.C

			fmt.Println("--- [Health Check] Executando ciclo de verificação ---")

			storeMutex.Lock()

			servicesToCheck := make([]Service, 0, len(serviceStore))

			for _, service := range serviceStore {
				servicesToCheck = append(servicesToCheck, service)
			}

			storeMutex.Unlock()

			for _, service := range servicesToCheck {
				go checkService(service)
			}
		}
	}()
}
