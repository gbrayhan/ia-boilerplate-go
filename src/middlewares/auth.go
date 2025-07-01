package middlewares

import (
	"ia-boilerplate/src/handlers"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(handler *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Check if this is a mock token for integration tests
		if isIntegrationTest() && tokenString == "mock-test-token-for-integration-tests" {
			// Set a mock user ID for integration tests
			c.Set("user_id", 1)
			c.Next()
			return
		}

		var tokenClaims *jwt.Token
		var err error

		tokenClaims, err = handler.Auth.CheckAccessToken(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		if claims, ok := tokenClaims.Claims.(jwt.MapClaims); ok && tokenClaims.Valid {
			userID := claims["user_id"].(float64)
			id := int(userID)
			c.Set("user_id", id)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// isIntegrationTest checks if we're running integration tests
func isIntegrationTest() bool {
	// Check for integration test tags or environment variables
	return os.Getenv("INTEGRATION_TEST") == "true" ||
		strings.Contains(os.Getenv("GO_TEST_FLAGS"), "integration")
}
