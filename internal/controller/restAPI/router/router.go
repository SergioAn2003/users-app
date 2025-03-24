package router

import (
	"users-app/internal/controller/restAPI/handler"
	"users-app/internal/controller/restAPI/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(mw *middlewares.Middleware, h *handler.Handler) chi.Router {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.RealIP, middleware.Recoverer, mw.Log)

		r.Get("/user", h.GetUserByID)
		r.Post("/user", h.CreateUser)
		r.Put("/user", h.UpdateUser)
		r.Delete("/user", h.DeleteUser)
	})

	return r
}
