package middleware

import (
	"errors"
	"fmt"

	"publy/util/config"
	"publy/util/logging"
	"publy/util/passwords"

	"github.com/gin-gonic/gin"
)

// Compares the provided token with the stored hash for the given pub.
// If the token is valid, the request proceeds, otherwise, it is aborted with a 401 status.
func AuthMiddleware(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")

		if bearerToken == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "token missing"})
			return
		}

		token := bearerToken[len("Bearer "):]

		pub := c.Query("pub")
		pubConfig, exists := config.Pubs[pub]
		if !exists {
			logging.LogError(errors.New("Pub does not exist"), fmt.Sprintf("Pub '%s' does not exist in config.", pub))
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		if !passwords.CompareToHash(token, pubConfig.TokenHash) {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		c.Next()
	}
}
