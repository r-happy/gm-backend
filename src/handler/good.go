package handler

import (
	"back/model"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func AddGoods(c echo.Context) error {
	goods := new(model.Good)
	if err := c.Bind(goods); err != nil {
		return err
	}

	id := c.Param("id")
	s := model.FindSpace(&model.Space{ID: id})
	if s.ID == "" {
		return echo.ErrNotFound
	}
	goods.SpaceID = s.ID

	if goods.GoodName == "" {
		return &echo.HTTPError{
			Code:    400,
			Message: "invalid goods name",
		}
	}

	goods.GoodID, _ = generateUniqueID()
	goods.AddEmail = userEmailFromToken(c)
	goods.CanBorrow = true

	model.CreateGood(goods)

	return c.JSON(http.StatusCreated, goods)
}

func GetGood(c echo.Context) error {
	sid := c.Param("sid")
	gid := c.Param("gid")

	space := model.FindSpace(&model.Space{ID: sid})
	if space.ID == "" {
		return echo.ErrNotFound
	}

	good := model.FindGood(&model.Good{SpaceID: space.ID, GoodID: gid})
	if good.GoodID == "" {
		return echo.ErrNotFound
	}

	email := userEmailFromToken(c)
	if user := model.FindUser(&model.User{Email: email}); user.ID == 0 {
		return echo.ErrNotFound
	}

	if !IsUserMemberOfSpace(email, sid) {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, good)
}

func GetGoods(c echo.Context) error {
	id := c.Param("id")

	space := model.FindSpace(&model.Space{ID: id})
	if space.ID == "" {
		return echo.ErrNotFound
	}

	email := userEmailFromToken(c)
	user := model.FindUser(&model.User{Email: email})
	if user.ID == 0 {
		return echo.ErrNotFound
	}

	if !IsUserMemberOfSpace(email, id) {
		return echo.ErrNotFound
	}

	goods := model.FindGoods(&model.Good{SpaceID: space.ID})

	return c.JSON(http.StatusOK, goods)
}

func ToggleGood(c echo.Context) error {
	sid := c.Param("sid")
	gid := c.Param("gid")

	email := c.FormValue("email")
	viewedStatusStr := c.FormValue("viewed_status")
	viewedStatusBool, _ := strconv.ParseBool(viewedStatusStr)

	if email == "" {
		return &echo.HTTPError{
			Code:    400,
			Message: "invalid email",
		}
	}

	good := model.FindGood(&model.Good{GoodID: gid, SpaceID: sid})

	if viewedStatusBool != good.Status {
		return echo.ErrForbidden
	}

	if good.GoodID == "" {
		return echo.ErrNotFound
	}

	if !good.CanBorrow {
		return echo.ErrForbidden
	}

	if !good.Status {
		// memberとして有効か, 管理者か確認
		member := model.FindMembers(&model.Member{Email: email, Space: sid})
		if len(member) == 0 {
			return echo.ErrNotFound
		}
		if !member[0].Admin {
			return echo.ErrForbidden
		}

		good.Status = true
		good.WhoBorrowUid = email
		good.WhoBorrowName = member[0].Name
		good.WhenBorrow = time.Now()
	} else {
		good.Status = false
	}

	model.SaveGood(&good)

	return c.JSON(http.StatusOK, good)
}
