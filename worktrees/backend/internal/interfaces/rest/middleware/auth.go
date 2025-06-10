package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Auth creates an authentication middleware
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

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Set user ID in context
			if userID, ok := claims["sub"].(string); ok {
				c.Set("user_id", userID)
			}
			
			// Set user roles in context
			if roles, ok := claims["roles"].([]interface{}); ok {
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