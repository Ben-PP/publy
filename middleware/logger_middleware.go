package middleware

import (
	"publy/util/logging"

	"github.com/gin-gonic/gin"
)

// Logger is a Gin middleware that logs each request using the logging package.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		logging.LogReq(
			c.Request.Host,
			c.ClientIP(),
			c.Request.Method,
			c.FullPath(),
			c.Request.UserAgent(),
			c.Writer.Status(),
		)
	}
}
