package main

import (
	"back/handler"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func newRouter() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.POST("/user", handler.CreateUser)
	e.POST("/login", handler.Login)
	e.GET("/", func (c echo.Context) error {
		return c.String(200, "Hello, World!!")
	})

	api := e.Group("/api")
	api.Use(middleware.JWTWithConfig(handler.Config))
	api.POST("/space", handler.AddSpace)
	api.GET("/space/:id", handler.GetSpace)
	api.GET("/space/:id/member", handler.GetMembers)
	api.POST("/space/:id/member", handler.AddMembers)
	api.GET("/spaces", handler.GetSpaces)
	// api.GET("/space/:id/good", handler.GetGoods)
	// api.GET("/space/:sid/good/:gid", handler.GetGood)
	api.POST("/space/:id/good", handler.AddGoods)
	api.GET("/profile", handler.GetProfile)

	return e
}
