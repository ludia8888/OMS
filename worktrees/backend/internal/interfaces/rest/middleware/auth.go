package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Auth creates an authentication middleware with enhanced security
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header missing",
			})
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			return
		}

		// Parse and validate token with options
		parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
		token, err := parser.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
				"details": err.Error(),
			})
			return
		}

		// Validate claims
		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			return
		}

		// Validate time-based claims
		now := time.Now()
		
		// Check expiration
		if claims.ExpiresAt != nil && now.After(claims.ExpiresAt.Time) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token expired",
			})
			return
		}

		// Check not before
		if claims.NotBefore != nil && now.Before(claims.NotBefore.Time) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token not yet valid",
			})
			return
		}

		// Check issued at
		if claims.IssuedAt != nil && now.Before(claims.IssuedAt.Time) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token issued in the future",
			})
			return
		}

		// Set user ID in context
		if claims.Subject != "" {
			c.Set("user_id", claims.Subject)
		}
		
		// Extract custom claims for roles
		if customClaims, ok := token.Claims.(jwt.MapClaims); ok {
			if roles, ok := customClaims["roles"].([]interface{}); ok {
				c.Set("user_roles", roles)
			}
		}

		c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserRoles extracts user roles from context
func GetUserRoles(c *gin.Context) []string {
	if roles, exists := c.Get("user_roles"); exists {
		if r, ok := roles.([]interface{}); ok {
			result := make([]string, len(r))
			for i, role := range r {
				if s, ok := role.(string); ok {
					result[i] = s
				}
			}
			return result
		}
	}
	return []string{}
}

// HasRole checks if user has a specific role
func HasRole(c *gin.Context, role string) bool {
	roles := GetUserRoles(c)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}