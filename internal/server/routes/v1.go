package routes

import (
	"net/http"

	"github.com/fransiscushermanto/backend/internal/config"
	v1 "github.com/fransiscushermanto/backend/internal/controllers/v1"
	"github.com/fransiscushermanto/backend/internal/controllers/v1/app"
	"github.com/fransiscushermanto/backend/internal/middlewares"
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
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
		// This middleware will run for every request to /api/v1/*
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				utils.Log().Info().Str("path", req.URL.Path).Msg("Request path entering /v1 routes")
				next.ServeHTTP(w, req)
			})
		})

		// Public routes
		r.Group(func(r chi.Router) {
			r.Use(cors.Default().Handler)
			r.Get("/health", v1.HealthCheck)
		})

		corsCfg := cors.Options{
			AllowedOrigins:   config.AllowedOrigins,
			AllowCredentials: true,
			Debug:            utils.IsDevelopment(),
		}

		// Protected routes
		r.Group(func(rProtected chi.Router) {
			rProtected.Use(middlewares.NewCorsMiddleware(corsCfg).Handler)
			rProtected.Options("/*", func(w http.ResponseWriter, r *http.Request) {})

			userController := v1.NewUserController(services.UserService)
			appController := v1.NewAppController(services.AppService, app.ControllerOptions{
				SecretKey: config.SecretKey,
			})
			authController := v1.NewAuthController(services.AuthService)

			r.Group(func(rAuthGroup chi.Router) {
				rAuthGroup.Post("/register", authController.Register)
				rAuthGroup.Post("/refresh", authController.Refresh)
				rAuthGroup.Post("/login", authController.Login)
			})

			rProtected.Route("/apps", func(rApps chi.Router) {
				rApps.Get("/", appController.GetApps)
				rApps.Post("/register", appController.RegisterApp)
			})

			rProtected.With(authMiddleware.RequireAuth).Group(func(rAuthed chi.Router) {
				rAuthed.Get("/users", userController.GetUsers)
				rAuthed.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
					if !utils.IsDevelopment() {
						http.NotFound(w, r)
						return
					}
					userController.GetUser(w, r)
				})

				rAuthed.Get("/profile", userController.Profile)
			})

		})
	})
}
