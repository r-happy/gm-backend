package handler

import (
	"back/model"
	"net/http"

	"github.com/labstack/echo"
)

func GetProfile(c echo.Context) error {
	email := userEmailFromToken(c)
	u := model.FindUser(&model.User{Email: email})
	u.Password = ""
	return c.JSON(http.StatusOK, u)
}
