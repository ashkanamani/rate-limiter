package main

import (
	"fmt"
	"github.com/ashkanamani/rate-limiter/config"
	"github.com/ashkanamani/rate-limiter/internal/limiter"
	"github.com/ashkanamani/rate-limiter/internal/middleware"
	"github.com/ashkanamani/rate-limiter/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/rueidis"
	"github.com/sirupsen/logrus"
	"time"
)

type App struct {
	Config      *config.Config
	Logger      *logrus.Logger
	RedisClient rueidis.Client
}

func main() {
	app := App{
		Config: config.LoadConfig(),
		Logger: utils.InitLogger(),
	}
	var err error
	app.RedisClient, err = limiter.InitRedis(app.Config)
	if err != nil {
		app.Logger.WithError(err).Fatalln("could not connect to redis")
	}
	app.Logger.Infoln("starting server on port", app.Config.ServerPort)

	r := gin.Default()

	// Initilize rate limiter
	rateLimiter := limiter.NewLimiter(app.RedisClient, 5, 1*time.Second, "ratelimiter:", limiter.FixedWindow)

	r.Use(middleware.NewRateLimiterMiddleware(rateLimiter))

	// Sample route for testing
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Start the server
	if err := r.Run(fmt.Sprintf(":%s", app.Config.ServerPort)); err != nil {
		app.Logger.WithError(err).Fatalln("error starting server")
	}

}
