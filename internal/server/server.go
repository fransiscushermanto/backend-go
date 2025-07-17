package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fransiscushermanto/backend/internal/config"
	"github.com/fransiscushermanto/backend/internal/repositories"
	"github.com/fransiscushermanto/backend/internal/server/routes"
	"github.com/fransiscushermanto/backend/internal/services"
	"github.com/fransiscushermanto/backend/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type APIServer struct {
	cfg        *config.AppConfig
	db         *repositories.Database
	httpServer *http.Server
}

func NewAPIServer(cfg *config.AppConfig, db *repositories.Database) *APIServer {
	return &APIServer{
		cfg: cfg,
		db:  db,
	}
}

func (s *APIServer) Run() error {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	appService := services.NewAppService(repositories.NewAppRepository(s.db), s.cfg.SecretKey)
	userService := services.NewUserService(repositories.NewUserRepository(s.db), appService)
	authService := services.NewAuthService(repositories.NewAuthRepository(s.db), userService, &s.cfg.PrivateKey, &s.cfg.PublicKey)

	router.Route("/api", func(r chi.Router) {
		routes.RoutesV1(r, &routes.RoutesOptions{
			Cfg: s.cfg,
			Services: &routes.Services{
				AppService:  appService,
				UserService: userService,
				AuthService: authService,
			},
		})
	})

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// use go keyword to run this function concurently (in background)
	// go keyword creates a goroutine
	go func() {
		utils.Log().Info().Int("port", s.cfg.Port).Msg("API server is starting...")

		certPath := "certs/cert.pem"
		keyPath := "certs/dev.pem"

		if err := s.httpServer.ListenAndServeTLS(certPath, keyPath); err != nil && err != http.ErrServerClosed {
			utils.Log().Fatal().Err(err).Msg("HTTP server faild to start")
		}
	}()

	quit := make(chan os.Signal, 1)
	// waiting for signal user Ctrl+C or SIGTERM
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.Log().Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	utils.Log().Info().Msg("Server gracefully stopped.")
	return nil
}
