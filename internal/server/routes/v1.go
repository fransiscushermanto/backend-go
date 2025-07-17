package routes

import (
	"fmt"
	"net/http"

	"github.com/fransiscushermanto/backend/internal/config"
	v1 "github.com/fransiscushermanto/backend/internal/controllers/v1"
	"github.com/fransiscushermanto/backend/internal/controllers/v1/app"
	"github.com/fransiscushermanto/backend/internal/middlewares"
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/rs/cors"

	"github.com/go-chi/chi/v5"
)

type Services struct {
	AppService  *services.AppService
	AuthService *services.AuthService
	UserService *services.UserService
}

type RoutesOptions struct {
	Cfg *config.AppConfig
	*Services
}

func RoutesV1(router chi.Router, options *RoutesOptions) {
	config := options.Cfg
	services := options.Services

	authMiddleware := middlewares.NewAuthMiddleware(services.AuthService)

	router.Route("/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Use(cors.Default().Handler)
			r.Get("/health", v1.HealthCheck)
		})

		corsCfg := cors.Options{
			AllowedOrigins:   config.AllowedOrigins,
			AllowCredentials: true,
			Debug:            config.Env == "development",
		}

		// Protected routes
		r.Group(func(r chi.Router) {
			fmt.Printf("Allowed Origins: %v\n", corsCfg.AllowedOrigins)
			r.Use(middlewares.NewCorsMiddleware(corsCfg).Handler)

			userController := v1.NewUserController(services.UserService)
			appController := v1.NewAppController(services.AppService, app.ControllerOptions{
				SecretKey: config.SecretKey,
			})
			authController := v1.NewAuthController(services.AuthService)

			r.Options("/register", func(w http.ResponseWriter, r *http.Request) {})
			r.Post("/register", authController.Register)

			r.With(authMiddleware.RequireAuth).Group(func(r chi.Router) {
				r.Get("/users", userController.GetUsers)
				r.Route("/apps", func(r chi.Router) {
					r.Get("/", appController.GetApps)
					r.Post("/register", appController.RegisterApp)
				})

				r.Get("/users/{id}", userController.GetUser)
			})

		})
	})
}
