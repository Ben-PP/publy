package main

import (
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"publy/middleware"
	"publy/util/config"
	"publy/util/logging"
	"publy/util/passwords"

	"github.com/gin-gonic/gin"
)

func generateHash(ctx *gin.Context) {
	// Allow only local requests
	re := regexp.MustCompile(`^(10|127|192\.168|172\.(1[6-9]|2[0-9]|3[0-1]))\.|^localhost`)
	if !re.MatchString(ctx.Request.Host) {
		logging.LogReq(
			ctx.Request.Host,
			ctx.ClientIP(),
			ctx.Request.Method,
			ctx.FullPath(),
			ctx.Request.UserAgent(),
			403,
		)
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}

	var req struct {
		String string `json:"string" binding:"required"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	hash, err := passwords.GenerateHash(req.String)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed to generate hash"})
		return
	}

	ctx.JSON(200, gin.H{"hash": hash})
}

func publish(ctx *gin.Context) {
	pub := ctx.Query("pub")

	logging.LogPublish(
		ctx.ClientIP(),
		ctx.Request.Method,
		ctx.FullPath(),
		ctx.Request.UserAgent(),
		true,
		pub,
	)

	ctx.JSON(200, gin.H{"message": "Publish triggered on pub: " + pub})
}

func main() {
	appLogger := logging.GetLogger()
	slog.SetDefault(appLogger)

	config, err := config.Get()
	if err != nil {
		logging.LogError(err, "Failed to initialize config.")
		return
	}

	router := gin.Default()

	router.Use(middleware.Logger())
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	{
		v1 := router.Group("/api/v1")
		v1.GET("/status", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"status": "ok"})
		})

		v1.GET("/publish", middleware.AuthMiddleware(*config), publish)
		v1.POST("/generate-hash", generateHash)

	}

	slog.Info("Starting server.")
	router.Run(fmt.Sprintf("%v:%v", config.Host, config.Port))
}
