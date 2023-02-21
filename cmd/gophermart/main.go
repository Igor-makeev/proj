package main

import (
	"fmt"
	"os"
	"os/signal"
	"proj/config"
	"proj/internal/handler"
	"proj/internal/repository"
	"proj/internal/repository/postgress"
	"proj/internal/server"
	"proj/internal/service"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	// ctx := context.Background()
	cfg := config.NewConfig()

	pc, err := postgress.NewPostgresClient(cfg)
	if err != nil {
		logrus.Fatalf("error:%v", err)
	}

	repo := repository.NewRepository(pc)
	service := service.NewService(repo, cfg)
	handler := handler.NewHandler(service)

	srv := new(server.Server)

	serverErrChan := srv.Run(cfg, handler)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-signals:

		fmt.Println("main: got terminate signal. Shutting down...")

		if err := srv.Shutdown(); err != nil {
			fmt.Printf("main: received an error while shutting down the server: %v", err)
		}

	case <-serverErrChan:
		fmt.Println("main: got server err signal. Shutting down...")

	}
}
