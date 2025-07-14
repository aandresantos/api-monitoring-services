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
