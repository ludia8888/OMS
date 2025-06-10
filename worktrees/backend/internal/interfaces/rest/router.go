package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openfoundry/oms/internal/config"
	"github.com/openfoundry/oms/internal/domain/repository"
	"github.com/openfoundry/oms/internal/domain/service"
	"github.com/openfoundry/oms/internal/interfaces/rest/handler"
	"github.com/openfoundry/oms/internal/interfaces/rest/middleware"
	"go.uber.org/zap"
)

// Services holds all the services needed by the router
type Services struct {
	ObjectTypeService *service.ObjectTypeService
	LinkTypeService   *service.LinkTypeService
}

// NewRouter creates a new HTTP router with all dependencies
func NewRouter(cfg *config.Config, services *Services, logger *zap.Logger) http.Handler {
	// Set Gin mode based on environment
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Cors(cfg.Security.AllowedOrigins))

	// Create handlers
	objectTypeHandler := handler.NewObjectTypeHandler(services.ObjectTypeService, logger)
	linkTypeHandler := handler.NewLinkTypeHandler(services.LinkTypeService, logger)

	// Health check endpoints
	router.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	router.GET("/health/ready", func(c *gin.Context) {
		// Check database connection via service
		if _, err := services.ObjectTypeService.List(c.Request.Context(), repository.ObjectTypeFilter{PageSize: 1}); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  "service not ready",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Authentication middleware for API routes
		v1.Use(middleware.Auth(cfg.Security.JWTSecret))

		// Object types endpoints
		objectTypes := v1.Group("/object-types")
		{
			objectTypes.GET("", objectTypeHandler.List)
			objectTypes.POST("", objectTypeHandler.Create)
			objectTypes.GET("/:id", objectTypeHandler.Get)
			objectTypes.PUT("/:id", objectTypeHandler.Update)
			objectTypes.DELETE("/:id", objectTypeHandler.Delete)
			objectTypes.GET("/:id/versions/compare", objectTypeHandler.CompareVersions)
		}

		// Link types endpoints
		linkTypes := v1.Group("/link-types")
		{
			linkTypes.GET("", linkTypeHandler.List)
			linkTypes.POST("", linkTypeHandler.Create)
			linkTypes.GET("/by-object-types", linkTypeHandler.GetByObjectTypes)
			linkTypes.POST("/validate-circular", linkTypeHandler.ValidateCircularReference)
			linkTypes.GET("/:id", linkTypeHandler.Get)
			linkTypes.PUT("/:id", linkTypeHandler.Update)
			linkTypes.DELETE("/:id", linkTypeHandler.Delete)
		}

		// Search endpoint
		v1.GET("/search", objectTypeHandler.Search)
	}

	// GraphQL endpoint (to be implemented)
	router.POST("/graphql", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "GraphQL not implemented yet"})
	})
	router.GET("/graphql", func(c *gin.Context) {
		c.String(http.StatusNotImplemented, "GraphQL Playground not implemented yet")
	})

	// Metrics endpoint
	if cfg.Metrics.Enabled {
		router.GET(cfg.Metrics.Path, func(c *gin.Context) {
			c.String(http.StatusNotImplemented, "Metrics not implemented yet")
		})
	}

	return router
}