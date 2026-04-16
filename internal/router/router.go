package router

import (
	"database/sql"
	"net/http"

	"github.com/ARKTEEK/shorty/internal/config"
	"github.com/ARKTEEK/shorty/internal/handlers"
	"github.com/ARKTEEK/shorty/internal/middleware"
	"github.com/ARKTEEK/shorty/internal/services"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func New(db *sql.DB, cfg *config.Config) *http.Server {
	userSvc := services.NewUserService(db)
	authSvc := services.NewAuthService(db, userSvc)
	linkSvc := services.NewLinkService(db, userSvc)
	userHandler := handlers.NewUserHandler(userSvc)
	authHandler := handlers.NewAuthHandler(authSvc)
	linkHandler := handlers.NewLinkHandler(linkSvc)

	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.Recoverer)
	r.Use(chimw.CleanPath)

	r.Route("/api/v1", func(r chi.Router) {

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.With(middleware.JWTAuth).Post("/deactivate", authHandler.Deactivate)
		})

		r.Route("/users", func(r chi.Router) {
			r.With(middleware.JWTAuth).Patch("/update", userHandler.UpdateUser)

			r.Route("/{id}", func(r chi.Router) {
				r.With(middleware.JWTAuth).Get("/", userHandler.GetUser)
			})
		})

		r.Route("/links", func(r chi.Router) {
			r.Post("/create", linkHandler.CreateShortLink)
			r.Get("/{shortCode}", linkHandler.Redirect)

			r.With(middleware.JWTAuth).Get("/mine", linkHandler.List)
			r.With(middleware.JWTAuth).Delete("/{shortCode}", linkHandler.Delete)
		})
	})

	return &http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}
}
