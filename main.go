package main

import (
	"pulse-service/apps/repository/adapter"
	"pulse-service/apps/repository/instance"
	"pulse-service/apps/routes"
	"pulse-service/constants"

	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"pulse-service/logger"
	"syscall"
	"time"
)

func main() {

	logger.InitLogger()
	logger.InitEventLogger()
	// configs := config.GetConfig()
	// aws := instance.GetAwsSession()
	RedisDBConnection := instance.GetRedisConnection()
	PSqlConnection := instance.GetPSqlConnection()

	repository := adapter.RepositoryAdapter(RedisDBConnection, PSqlConnection)

	fmt.Println("Starting %s API server", "pulse-service")

	server := &http.Server{
		Addr:    constants.PORT,
		Handler: routes.NewRouter().SetRouters(repository),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Error::%v", err)
			fmt.Println("Failed to start %s service\n", "pulse-service")
		}
	}()

	fmt.Println("Listening on port %v ", server.Addr)

	// queue := svc.NewServiceRepo(repository).SQSService
	// go queue.InitSQS()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	fmt.Println("Shutting down server.")
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Server forced to shutdown: %v\n", err)
	}
}
