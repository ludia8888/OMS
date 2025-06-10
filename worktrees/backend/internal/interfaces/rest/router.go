package rest

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openfoundry/oms/internal/config"
	"github.com/openfoundry/oms/internal/interfaces/rest/middleware"
	"go.uber.org/zap"
)

// NewRouter creates a new HTTP router
func NewRouter(cfg *config.Config, db *sql.DB, logger *zap.Logger) http.Handler {
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

	// Health check endpoints
	router.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	router.GET("/health/ready", func(c *gin.Context) {
		// Check database connection
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  "database connection failed",
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
			objectTypes.GET("", handleListObjectTypes)
			objectTypes.POST("", handleCreateObjectType)
			objectTypes.GET("/:id", handleGetObjectType)
			objectTypes.PUT("/:id", handleUpdateObjectType)
			objectTypes.DELETE("/:id", handleDeleteObjectType)
		}

		// Link types endpoints
		linkTypes := v1.Group("/link-types")
		{
			linkTypes.GET("", handleListLinkTypes)
			linkTypes.POST("", handleCreateLinkType)
			linkTypes.GET("/:id", handleGetLinkType)
			linkTypes.PUT("/:id", handleUpdateLinkType)
			linkTypes.DELETE("/:id", handleDeleteLinkType)
		}

		// Search endpoint
		v1.GET("/search", handleSearch)
	}

	// GraphQL endpoint (to be implemented)
	router.POST("/graphql", handleGraphQL)
	router.GET("/graphql", handleGraphQLPlayground)

	// Metrics endpoint
	if cfg.Metrics.Enabled {
		router.GET(cfg.Metrics.Path, handleMetrics)
	}

	return router
}

// Placeholder handlers - to be implemented
func handleListObjectTypes(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func handleCreateObjectType(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func handleGetObjectType(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func handleUpdateObjectType(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func handleDeleteObjectType(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func handleListLinkTypes(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func handleCreateLinkType(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func handleGetLinkType(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func handleUpdateLinkType(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func handleDeleteLinkType(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func handleSearch(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func handleGraphQL(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func handleGraphQLPlayground(c *gin.Context) {
	c.String(http.StatusNotImplemented, "GraphQL Playground not implemented")
}

func handleMetrics(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Metrics not implemented")
}