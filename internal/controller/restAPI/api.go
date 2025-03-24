package restapi

import (
	"context"
	"fmt"
	"net/http"
	"users-app/internal/controller/restAPI/handler"
	"users-app/internal/controller/restAPI/middlewares"
	"users-app/internal/controller/restAPI/router"
	"users-app/pkg/config"
	"users-app/pkg/logger"
)

type Controller struct {
	cfg         *config.Config
	log         logger.Logger
	userService handler.UserService
	srv         *http.Server
}

func New(cfg *config.Config, log logger.Logger, userService handler.UserService) *Controller {
	mw := middlewares.New(log)
	h := handler.New(log, userService)
	r := router.New(mw, h)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      r,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	return &Controller{
		cfg:         cfg,
		log:         log,
		userService: userService,
		srv:         server,
	}
}

func (c *Controller) Run() error {
	return c.srv.ListenAndServe()
}

func (c *Controller) Stop(ctx context.Context) error {
	return c.srv.Shutdown(ctx)
}
