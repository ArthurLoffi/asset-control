package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"asset-control/src/main/adapter/inbound/http/controller"
	"asset-control/src/main/adapter/inbound/http/routes"
	_ "asset-control/src/main/adapter/inbound/http/swagger"
	"asset-control/src/main/adapter/outbound/data/database"
	"asset-control/src/main/adapter/outbound/data/repository"
	"asset-control/src/main/usecase"

	authController "github.com/Brasil-Novo/asset-control-auth/src/main/adapter/inbound/http/controller"
	authMiddleware "github.com/Brasil-Novo/asset-control-auth/src/main/adapter/inbound/http/middleware"
	authRoutes "github.com/Brasil-Novo/asset-control-auth/src/main/adapter/inbound/http/routes"
	authRepository "github.com/Brasil-Novo/asset-control-auth/src/main/adapter/outbound/data/repository"
	"github.com/Brasil-Novo/asset-control-auth/src/main/adapter/security"
	authUsecase "github.com/Brasil-Novo/asset-control-auth/src/main/usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// @title           Asset control API
// @version         v0.3.0
// @description     Está é uma API utilizada para controlar os equipamentos de TI.

// @contact.name   Suporte BNTG
// @contact.email  suporte@brasilnovo.com

// @host      localhost:8080
// @BasePath  /api

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	authCfg, err := security.ConfigFromEnv()
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	cfg, err := database.ConfigFromEnv()
	if err != nil {
		return err
	}
	db, err := database.Connect(ctx, cfg)
	if err != nil {
		return err
	}
	defer func() {
		if err := database.Disconnect(db.Client()); err != nil {
			log.Printf("disconnect MongoDB: %v", err)
		}
	}()
	schemaCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	err = database.EnsureSchema(schemaCtx, db)
	cancel()
	if err != nil {
		return err
	}

	users := authRepository.NewUserRepository(db)
	schemaCtx, cancel = context.WithTimeout(ctx, cfg.Timeout)
	err = users.EnsureSchema(schemaCtx)
	cancel()
	if err != nil {
		return errors.New("could not create unique users.name index; check permissions and duplicate names")
	}
	r, err := newRouter(db, authCfg)
	if err != nil {
		return err
	}

	server := &http.Server{Addr: ":8080", Handler: r, ReadHeaderTimeout: 5 * time.Second}
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

func newRouter(db *mongo.Database, authCfg security.Config) (*gin.Engine, error) {
	tokens, err := security.NewTokens(authCfg.Secret, authCfg.TTL)
	if err != nil {
		return nil, err
	}
	passwords, err := security.NewPasswords()
	if err != nil {
		return nil, err
	}
	login := authUsecase.NewLoginUseCase(authRepository.NewUserRepository(db), passwords, tokens, security.NewAttempts())
	employeeRepository := repository.NewEmployeeRepository(db)
	employeeUseCase := usecase.NewEmployeeUseCase(employeeRepository)
	employeeController := controller.NewEmployeeController(employeeUseCase)

	machineRepository := repository.NewMachineRepository(db)
	machineUseCase := usecase.NewMachineUseCase(machineRepository)
	machineController := controller.NewMachineController(machineUseCase)

	phoneRepository := repository.NewPhoneRepository(db)
	phoneUseCase := usecase.NewPhoneUseCase(phoneRepository)
	phoneController := controller.NewPhoneController(phoneUseCase)

	monitorRepository := repository.NewMonitorRepository(db)
	monitorUseCase := usecase.NewMonitorUseCase(monitorRepository)
	monitorController := controller.NewMonitorController(monitorUseCase)

	printerRepository := repository.NewPrinterRepository(db)
	printerUseCase := usecase.NewPrinterUseCase(printerRepository)
	printerController := controller.NewPrinterController(printerUseCase)

	utilsRepository := repository.NewUtilsRepository(db)
	utilsUseCase := usecase.NewUtilsUseCase(utilsRepository)
	utilsController := controller.NewUtilsController(utilsUseCase)

	emailRepository := repository.NewEmailRepository(db)
	emailUseCase := usecase.NewEmailUseCase(emailRepository)
	emailController := controller.NewEmailController(emailUseCase)

	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.Use(gin.Logger(), authMiddleware.Recovery(), authMiddleware.RequestID(), authMiddleware.BrowserSecurity(authCfg.Origins))
	authRoutes.Register(r.Group("/api"), authController.NewAuthController(login, authCfg.Secure), tokens)
	// All inventory routes and Swagger require a verified session.
	r.Use(authMiddleware.RequireSession(tokens))

	routes.Routes(
		r,
		employeeController,
		machineController,
		phoneController,
		monitorController,
		printerController,
		utilsController,
		emailController,
	)

	return r, nil
}
