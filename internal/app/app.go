// Package app is the composition root: it registers every module on one router.
package app

import (
	"asset-control/internal/core/httpx"
	"asset-control/internal/core/httpx/middleware"
	"asset-control/internal/modules/audit"
	"asset-control/internal/modules/inventory"
	"asset-control/internal/modules/organization"

	_ "asset-control/internal/app/swagger"

	auth "github.com/Brasil-Novo/asset-control-auth"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// New builds the HTTP handler. responseCache may be nil to disable caching.
func New(db *mongo.Database, authModule *auth.Module, responseCache middleware.ResponseCache) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	if err := configureClientIP(router); err != nil {
		return nil, err
	}
	router.Use(gin.Logger(), auth.RequestID(), httpx.Recovery(), authModule.BrowserSecurity())
	authModule.RegisterRoutes(router.Group("/api"))
	// Groups copy the middleware registered so far: every group created below
	// requires a verified session, including Swagger.
	router.Use(middleware.RateLimit(), authModule.RequireSession(), middleware.AuditIdentity())
	if responseCache != nil {
		router.Use(middleware.InventoryCache(responseCache))
	}

	api := router.Group("/api")
	recorder := audit.Register(api, db)
	inventory.Register(api, db, recorder)
	organization.Register(api, db)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router, nil
}

func configureClientIP(router *gin.Engine) error {
	// cloudflared connects locally and supplies the authoritative Cloudflare IP.
	// Only a local peer may supply it; direct clients cannot spoof forwarded IPs.
	router.RemoteIPHeaders = []string{gin.PlatformCloudflare}
	return router.SetTrustedProxies([]string{"127.0.0.1", "::1"})
}
