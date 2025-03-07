package handler

import (
	"back/model"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

func AddGoods(c echo.Context) error {
	good := new(model.Good)
	if err := c.Bind(good); err != nil {
		return err
	}

	borrowUserEmailsString := c.FormValue("borrowUserEmails")
	if borrowUserEmailsString == "" {
		return &echo.HTTPError{
			Code:    400,
			Message: "invalid borrow user emails",
		}
	}

	id := c.Param("id")
	s := model.FindSpace(&model.Space{ID: id})
	if s.ID == "" {
		return echo.ErrNotFound
	}
	good.SpaceID = s.ID

	if good.GoodName == "" {
		return &echo.HTTPError{
			Code:    400,
			Message: "invalid goods name",
		}
	}

	good.GoodID, _ = generateUniqueID()
	good.AddEmail = userEmailFromToken(c)
	good.CanBorrow = true

	model.CreateGood(good)

	borrowUserEmails := strings.Split(borrowUserEmailsString, ",")
	for i := range borrowUserEmails {
		var borrowUser model.BorrowUser
		tmp := strings.TrimSpace(borrowUserEmails[i])
		isMember := model.FindMembers(&model.Member{Email: tmp, Space: good.SpaceID})

		if len(isMember) == 0 {
			return echo.ErrNotFound
		}

		borrowUser.Email = tmp
		borrowUser.GoodID = good.GoodID
		borrowUser.Name = isMember[0].Name
		model.CreateBorrowUser(&borrowUser)
	}

	return c.JSON(http.StatusCreated, good)
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
	myemail := userEmailFromToken(c)
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

	mymemberinfo := model.FindMembers(&model.Member{Email: myemail, Space: sid})
	if len(mymemberinfo) == 0 {
		return echo.ErrNotFound
	}
	if !mymemberinfo[0].Admin {
		return echo.ErrForbidden
	}

	if !good.Status {
		member := model.FindBorrowUser(&model.BorrowUser{GoodID: gid, Email: email})
		fmt.Println(member)
		if member.Email == "" {
			return echo.ErrNotFound
		}

		good.Status = true
		good.WhoBorrowUid = email
		good.WhoBorrowName = member.Name
		good.WhenBorrow = time.Now()
	} else {
		hisory := model.History{
			SpaceID:    sid,
			GoodID:     gid,
			GoodName:   good.GoodName,
			BorrowUser: good.WhoBorrowName,
			WhenBorrow: good.WhenBorrow,
			WhenReturn: time.Now(),
		}
		model.CreateHistory(&hisory)
		good.Status = false
	}

	model.SaveGood(&good)

	return c.JSON(http.StatusOK, good)
}

func GetBorrowUser(c echo.Context) error {
	sid := c.Param("id")
	gid := c.Param("gid")

	email := userEmailFromToken(c)
	if user := model.FindUser(&model.User{Email: email}); user.ID == 0 {
		return echo.ErrNotFound
	}

	if !IsUserMemberOfSpace(email, sid) {
		return echo.ErrNotFound
	}

	good := model.FindGood(&model.Good{GoodID: gid, SpaceID: sid})
	if good.GoodID == "" {
		return echo.ErrNotFound
	}

	borrowUser := model.FindBorrowUsers(&model.BorrowUser{GoodID: gid})

	return c.JSON(http.StatusOK, borrowUser)
}

func UpdateGood(c echo.Context) error {
	sid := c.Param("sid")
	gid := c.Param("gid")

	borrowUserEmailsString := c.FormValue("borrowUserEmails")
	if borrowUserEmailsString == "" {
		return &echo.HTTPError{
			Code:    400,
			Message: "invalid borrow user emails",
		}
	}

	good := model.FindGood(&model.Good{GoodID: gid, SpaceID: sid})

	email := userEmailFromToken(c)

	user := model.FindMembers(&model.Member{Email: email, Space: sid})
	if len(user) == 0 {
		return echo.ErrNotFound
	}

	if !user[0].Admin {
		return echo.ErrForbidden
	}

	borrowUserEmails := strings.Split(borrowUserEmailsString, ",")
	model.RemoveBorrowUser(gid)

	for i := range borrowUserEmails {
		var borrowUser model.BorrowUser
		tmp := strings.TrimSpace(borrowUserEmails[i])
		isMember := model.FindMembers(&model.Member{Email: tmp, Space: good.SpaceID})

		if len(isMember) == 0 {
			return echo.ErrNotFound
		}

		borrowUser.Email = tmp
		borrowUser.GoodID = good.GoodID
		borrowUser.Name = isMember[0].Name
		model.CreateBorrowUser(&borrowUser)
	}

	good.SpaceID = sid
	good.GoodID = gid
	good.GoodName = c.FormValue("goodName")
	good.Description = c.FormValue("description")

	model.SaveGood(&good)

	return c.JSON(http.StatusOK, good)
}
