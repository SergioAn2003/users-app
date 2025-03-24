package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"users-app/internal/controller/restAPI"
	"users-app/internal/repository"
	"users-app/internal/service"
	"users-app/pkg/config"
	"users-app/pkg/logger"
	"users-app/pkg/postgres"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New(".env")
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		return
	}

	log, err := logger.New(cfg.Mode)
	if err != nil {
		fmt.Printf("failed to create logger: %v\n", err)
		return
	}

	pool, err := postgres.Connect(ctx, cfg.Postgres.DSN, cfg.Postgres.MaxConns)
	if err != nil {
		log.ErrorF("failed to connect to database: %s", err.Error())
		return
	}

	if err := postgres.UpMigrations(pool); err != nil {
		log.ErrorF("failed to run migrations: %s", err.Error())
		return
	}

	userRepo := repository.New(pool)
	userService := service.New(userRepo)
	restController := restapi.New(cfg, log, userService)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := restController.Run(); err != nil {
			log.ErrorF("failed to run server: %s", err.Error())
			return
		}
	}()

	log.InfoF("server started on port %d", cfg.HTTP.Port)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	sig := <-ch

	log.InfoF("got OS signal: %s\n", sig)

	if err := restController.Stop(ctx); err != nil {
		log.ErrorF("failed to stop server: %s", err.Error())
		return
	}

	wg.Wait()
}
