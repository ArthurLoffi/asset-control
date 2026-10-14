package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"asset-control/internal/app"
	"asset-control/internal/core/cache"
	"asset-control/internal/core/config"
	"asset-control/internal/core/database"
	"asset-control/internal/core/httpx/middleware"

	auth "github.com/Brasil-Novo/asset-control-auth"
)

// @title           Asset control API
// @version         v0.3.0
// @description     Está é uma API utilizada para controlar os equipamentos de TI.

// @contact.name   Suporte BNTG
// @contact.email  suporte@brasilnovo.com

// @host      localhost:8060
// @BasePath  /api

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	environment, err := config.Load()
	if err != nil {
		return err
	}
	authConfig, err := authConfigFromEnv(environment)
	if err != nil {
		return err
	}
	databaseConfig, err := database.ConfigFromEnv(environment)
	if err != nil {
		return err
	}
	cacheConfig, err := cache.ConfigFromEnv(environment, databaseConfig.Database)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := database.Connect(ctx, databaseConfig)
	if err != nil {
		return err
	}
	defer func() {
		if err := database.Disconnect(db.Client()); err != nil {
			log.Printf("disconnect MongoDB: %v", err)
		}
	}()
	authModule, err := auth.New(db, authConfig)
	if err != nil {
		return err
	}
	schemaCtx, cancel := context.WithTimeout(ctx, databaseConfig.Timeout)
	defer cancel()
	if err := app.EnsureSchema(schemaCtx, db); err != nil {
		return err
	}
	if err := authModule.EnsureSchema(schemaCtx); err != nil {
		return errors.New("could not create unique users.name index; check permissions and duplicate names")
	}

	responseCache, err := cache.Connect(ctx, cacheConfig)
	if err != nil {
		return err
	}
	var routerCache middleware.ResponseCache
	if responseCache != nil {
		routerCache = responseCache
		defer func() {
			if err := responseCache.Close(); err != nil {
				log.Printf("disconnect Redis: %v", err)
			}
		}()
	}
	router, err := app.New(db, authModule, routerCache)
	if err != nil {
		return err
	}
	return serve(ctx, &http.Server{Addr: ":8060", Handler: router, ReadHeaderTimeout: 5 * time.Second})
}

func authConfigFromEnv(environment config.Environment) (auth.Config, error) {
	secret, err := environment.Require("AUTH_JWT_SECRET")
	if err != nil {
		return auth.Config{}, err
	}
	ttl, err := environment.ReadDuration("AUTH_TOKEN_TTL", 8*time.Hour)
	if err != nil {
		return auth.Config{}, err
	}
	isSecure, err := environment.ReadBool("AUTH_COOKIE_SECURE", true)
	if err != nil {
		return auth.Config{}, err
	}
	origins := environment.ReadList("AUTH_ALLOWED_ORIGINS", "http://localhost:5173")
	return auth.Config{Secret: secret, TTL: ttl, IsSecure: isSecure, Origins: origins}, nil
}

// serve drains in-flight requests on SIGINT/SIGTERM.
func serve(ctx context.Context, server *http.Server) error {
	serverErrors := make(chan error, 1)
	go func() { serverErrors <- server.ListenAndServe() }()
	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return err
		}
		return nil
	}
}
