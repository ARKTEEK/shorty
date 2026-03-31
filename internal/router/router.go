package router

import (
	"database/sql"
	"net/http"

	"github.com/ARKTEEK/shorty/internal/config"
	"github.com/ARKTEEK/shorty/internal/handlers"
	"github.com/ARKTEEK/shorty/internal/services"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func New(db *sql.DB, cfg *config.Config) *http.Server {
	userSvc := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userSvc)

	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			// r.Get("/", userHandler.ListUsers)
			r.Post("/", userHandler.CreateUser)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", userHandler.GetUser)
				// r.Delete("/", userHandler.DeleteUser)
			})
		})
	})

	return &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}
}
