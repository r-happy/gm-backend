package handler

import (
	"back/model"
	"net/http"

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
            Code:   400,
            Message: "invalid goods name",
        }
    }

    goods.GoodID, _ = generateUniqueID()
    goods.AddEmail = userEmailFromToken(c)

    model.CreateGood(goods)

    return c.JSON(http.StatusCreated, goods)
}

func GetGood(c echo.Context) error {
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

    if !IsUserMemberOfSpace(email, id){
        return echo.ErrNotFound
    }

    good := model.FindGood(&model.Good{SpaceID: space.ID})

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

    if !IsUserMemberOfSpace(email, id){
        return echo.ErrNotFound
    }

	goods := model.FindGoods(&model.Good{SpaceID: space.ID})

    return c.JSON(http.StatusOK, goods)
}