package statistics

import (
	"github.com/Luna-devv/nekostic/config"
	"github.com/labstack/echo"
)

type route config.ApiContext

func NewRouter(context config.ApiContext, c *echo.Group) {
	route := route(context)
	c.GET("/statistics", route.Statistics)
}
