package handler

import (
	"back/model"
	"net/http"

	"github.com/labstack/echo"
)

func GetHistories(c echo.Context) error {
	sid := c.Param("sid")
	myemail := userEmailFromToken(c)

	if !IsUserMemberOfSpace(myemail, sid) {
		return echo.ErrNotFound
	}

	histories := model.FindHistories(&model.History{SpaceID: sid})
	return c.JSON(http.StatusOK, histories)
}
