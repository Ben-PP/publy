package middleware

import (
	"github.com/gin-gonic/gin"
)

// Tries to extract the JWT from Authorization header. Returns an error status
// to the client if it fails.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")

		if bearerToken == "" {
			c.Set("x-authorized", false)
			c.Abort()
			return
		}

		token := bearerToken[len("Bearer "):]

		c.Set("x-token", token)

		// logging.LogTokenEvent(true, c.FullPath(), logging.TokenEventTypeAccess, c.RemoteIP(), token)
		c.Next()
	}
}
