package handler

import (
	"back/model"
	"net/http"


	"github.com/labstack/echo"
)

// GetMembers returns all members of a space.
func GetMembers(c echo.Context) error {
    email := userEmailFromToken(c)
    id := c.Param("id")

    user := model.FindUser(&model.User{Email: email})
    if user.ID == 0 {
        return echo.ErrNotFound
    }

    if !IsUserMemberOfSpace (email, id) {
        return echo.ErrNotFound
    }

    
    members := model.FindMembers(&model.Member{Space: id})

    if len(members) == 0 {
        return echo.ErrNotFound
    }

    return c.JSON(http.StatusOK, members)
}

// AddMembers adds a member to a space.
func AddMembers(c echo.Context) error {
	id := c.Param("id")
	addemail := c.FormValue("email")
	adminStr := c.FormValue("admin")
	admin := adminStr == "true"

	if adminStr == "" {
		return &echo.HTTPError{
			Code:    400,
			Message: "invalid admin",
		}
	}

	if addemail == "" {
		return &echo.HTTPError{
			Code:    400,
			Message: "invalid email",
		}
	}

	if addmember := model.FindUser(&model.User{Email: addemail}); addmember.ID == 0 {
		return echo.ErrNotFound
	}

	if space := model.FindSpace(&model.Space{ID: id}); space.ID == "" {
		return echo.ErrNotFound
	}

	email := userEmailFromToken(c)
	user := model.FindUser(&model.User{Email: email})
	if user.ID == 0 {
		return echo.ErrNotFound
	}

	member := &model.Member{
		Space: id,
		Email: addemail,
		Admin: admin,
	}

	model.CreateMember(member)

	return c.JSON(http.StatusCreated, member)
}
