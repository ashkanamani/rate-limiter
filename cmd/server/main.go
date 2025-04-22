package main

import (
	"github.com/ashkanamani/rate-limiter/config"
	"github.com/ashkanamani/rate-limiter/internal/limiter"
	"github.com/ashkanamani/rate-limiter/internal/utils"
	"github.com/redis/rueidis"
	"github.com/sirupsen/logrus"
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
}
