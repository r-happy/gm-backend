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
	e.POST("/user", handler.CreateUser) // Body: { "username": "string", "password": "string", "email": "string" }
	e.POST("/login", handler.Login)     // Body: { "email": "string", "password": "string" }
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!!")
	})

	// JWT Middleware
	api := e.Group("/api")
	api.Use(middleware.JWTWithConfig(handler.Config))
	api.POST("/space", handler.AddSpace) // Body: { "space_name": "string", "parent_id": "string" }

	/*
		/space/:id
		スペースの詳細を取得
		Parameters:
			id: string
	*/
	api.GET("/space/:id", handler.GetSpace)

	/*
		/space/:id/member
		スペースのメンバーを取得
		Parameters:
			id: string
	*/
	api.GET("/space/:id/member", handler.GetMembers)

	/*
		/space/:id/member
		スペースのメンバーを追加
		Parameters:
			id: string
		Body: { "email": "string", "name": "string", "admin": "bool" }
	*/
	api.POST("/space/:id/member", handler.AddMembers)

	/*
		/space/:id/admin
		管理者権限の切り替え
		Parameters:
			id: string
		Body: { "emial": "string" }
	*/
	api.POST("/space/:id/admin", handler.ToggleMemberAdmin)

	/*
		/spaces
		全てのスペースを取得
		Parameters:
			なし
	*/
	api.GET("/spaces", handler.GetSpaces)

	/*
		/space/:id/children
		スペースの子スペースを取得
		Parameters:
			id: string
	*/
	api.GET("/space/:id/children", handler.GetChildrens)

	/*
		/space/:id/good
		スペースのグッズを取得
		Parameters:
			id: string
	*/
	api.GET("/space/:id/good", handler.GetGoods)
	
	/*
		/space/:id/good/:gid/borrow
		グッズの借りれる人を返す
		Parameters:
			id: string
			gid: string
	*/
	api.GET("/space/:id/good/:gid/borrow", handler.GetBorrowUser)

	/*
		/space/:sid/good/:gid
		特定のグッズを取得
		Parameters:
			sid: string
			gid: string
	*/
	api.GET("/space/:sid/good/:gid", handler.GetGood)

	/*
		/space/:sid/good/:gid
		グッズの状態をトグル
		Parameters:
			sid: string
			gid: string
		Body: { "email": "string", "viewed_status": "string (true or false)" }
	*/
	api.POST("/space/:sid/good/:gid", handler.ToggleGood)

	/*
		/space/:id/good
		スペースにグッズを追加
		Parameters:
			id: string
		Body: { "good_name": "string", "description": "string", "borrow_user_emails": "string" }
	*/
	api.POST("/space/:id/good", handler.AddGoods)

	/*
		/profile
		ユーザープロフィールを取得
		Parameters:
			なし
	*/
	api.GET("/profile", handler.GetProfile)

	return e
}
