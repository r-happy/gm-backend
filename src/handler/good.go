package handler

import (
	"back/model"
	"fmt"
	"net/http"
	"time"
	"strconv"

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

	if user := model.FindUser(&model.User{Email: email}); user.ID == 0 {
		return echo.ErrNotFound
	}

	good := model.FindGood(&model.Good{GoodID: gid, SpaceID: sid})
	
    if (viewedStatusBool != good.Status) {
        return echo.ErrForbidden
    }

	if good.GoodID == "" {
		return echo.ErrNotFound
	}

	if !good.CanBorrow {
		return echo.ErrForbidden
	}
	fmt.Println(good.Status)

	// 状態を切り替える
	if good.Status {
		good.Status = false
		good.WhoBorrowUid = email

		userData := model.FindUser(&model.User{Email: email})

		good.WhoBorrowName = userData.Name
		good.WhenBorrow = time.Now()
	} else {
		good.Status = true
		good.WhoReturnUid = email

		userData := model.FindUser(&model.User{Email: email})
		good.WhoReturnName = userData.Name
	}

	model.SaveGood(&good)

	return c.JSON(http.StatusOK, good)
}
