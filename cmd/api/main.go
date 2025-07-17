package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-monitoring-services/internal/database"
	"api-monitoring-services/internal/handler"
	"api-monitoring-services/internal/pkg/healthcheck"
	"api-monitoring-services/internal/repository"
	"api-monitoring-services/internal/service"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	errLoadEnv := godotenv.Load()

	if errLoadEnv != nil {
		log.Println("⚠️ .env não encontrado")
	}

	dbClient, errConnectDB := database.ConnectDB()

	if errConnectDB != nil {

		log.Fatalf("Erro ao conectar no banco: %v", errConnectDB)
	}
	defer database.CloseConnectionDB()

	checkInterval := 30 * time.Second

	checker := healthcheck.NewHTTPHealthChecker(5 * time.Second)

	repo := repository.NewServiceRepository(dbClient)
	serviceManager := service.NewServiceManager(repo, checker, checkInterval)
	serviceHandler := handler.NewServiceHandler(serviceManager)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	serviceHandler.RegisterRoutes(e)
	serviceManager.StartHealthChecks()

	go func() {
		StartServer(":3000", e)
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
