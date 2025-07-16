package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-monitoring-services/internal/handler"
	"api-monitoring-services/internal/repository"
	"api-monitoring-services/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	checkInterval := 30 * time.Second

	repo := repository.NewInMemoryServiceRepository()
	serviceManager := service.NewServiceManager(repo, checkInterval)
	serviceHandler := handler.NewServiceHandler(serviceManager)

	Echo := echo.New()
	Echo.Use(middleware.Logger())
	Echo.Use(middleware.Recover())

	serviceHandler.RegisterRoutes(Echo)

	go func() {
		StartServer(":3000", Echo)
	}()

	// TODO: implementar o resto do graceful
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

}

func StartServer(address string, e *echo.Echo) {
	if err := e.Start(address); err != nil {
		log.Printf("Erro ao iniciar servidor: %v", err)
	}
}
