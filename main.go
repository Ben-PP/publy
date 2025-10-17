package main

import (
	"fmt"
	"log/slog"
	"time"

	"publy/middleware"
	"publy/util/config"
	"publy/util/logging"

	"github.com/gin-gonic/gin"
)

func authorizePub(pub string, token string) bool {
	config, err := config.Get()
	if err != nil {
		logging.LogError(err, "Failed to load config.")
		return false
	}

	pubConfig, exists := config.Pubs[pub]
	if !exists {
		logging.LogError(err, fmt.Sprintf("Pub '%s' does not exist in config.", pub))
		return false
	}

	if pubConfig.Token != token {
		return false
	}

	return true
}

func publish(ctx *gin.Context) {
	pub := ctx.Query("pub")
	println("Publishing to pub:", pub)
	authorized := authorizePub(pub, ctx.GetString("x-token"))
	println("Authorized:", authorized)

	if !authorized {
		logging.LogPublish(
			ctx.ClientIP(),
			ctx.Request.Method,
			ctx.FullPath(),
			ctx.Request.UserAgent(),
			false,
			pub,
		)
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

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

		v1.GET("/publish", middleware.AuthMiddleware(), publish)

	}

	slog.Info("Starting server.")
	router.Run(fmt.Sprintf("%v:%v", config.Host, config.Port))
}
