package api

import (
	"net/http"

	"github.com/Luna-devv/nekostic/api/routes/statistics"
	"github.com/Luna-devv/nekostic/config"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/redis/go-redis/v9"
)

func New(conf config.Config) {
	e := echo.New()

	api := e.Group("")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"wamellow.com,local.wamellow.com"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodOptions},
		AllowCredentials: true,
	}))

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		Username: conf.Redis.Username,
		DB:       conf.Redis.Db,
	})

	context := config.ApiContext{Config: conf, Client: client}
	statistics.NewRouter(context, api)

	e.Logger.Fatal(e.Start(":" + conf.Port))
}
