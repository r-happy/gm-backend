package handler

import "github.com/labstack/echo"
import (
	"back/model"
	"net/http"
)

func GetProfile(c echo.Context) error {
	email := userEmailFromToken(c)
	u := model.FindUser(&model.User{Email: email})
	u.Password = ""
	return c.JSON(http.StatusOK, u)
}