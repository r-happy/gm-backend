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