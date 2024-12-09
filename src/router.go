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
	e.POST("/user", handler.CreateUser)
	e.POST("/login", handler.Login)

	api := e.Group("/api")
	api.Use(middleware.JWTWithConfig(handler.Config))
	api.POST("/space", handler.AddSpace)
	api.GET("/space/:id", handler.GetSpace)
	api.GET("/space/:id/member", handler.GetMembers)
	api.POST("/space/:id/member", handler.AddMembers)
	api.GET("/spaces", handler.GetSpaces)
	// api.GET("/space/:id/good", handler.GetGoods)
	// api.GET("/space/:sid/good/:gid", handler.GetGood)
	// api.POST("/space/:id/good", handler.AddGoods)

	return e
}
