package handler

import (
	"back/model"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

// AddSpace creates a new space.
func AddSpace(c echo.Context) error {
	space := new(model.Space)
	if err := c.Bind(space); err != nil {
		return err
	}

	if space.SpaceName == "" {
		return &echo.HTTPError{
			Code:    400,
			Message: "invalid space name",
		}
	}

	email := userEmailFromToken(c)
	if user := model.FindUser(&model.User{Email: email}); user.ID == 0 {
		return echo.ErrNotFound
	}

	space.ID, _ = generateUniqueID()
	space.TimeOfBorn = time.Now()
	model.CreateSpace(space)

	member := &model.Member{
		Space: space.ID,
		Email: email,
		Admin: true,
	}
	model.CreateMember(member)

	return c.JSON(http.StatusCreated, space)
}

// GetSpaces returns all spaces that the user is a member of.
func GetSpaces(c echo.Context) error {
	email := userEmailFromToken(c)
	user := model.FindUser(&model.User{Email: email})
	if user.ID == 0 {
		return echo.ErrNotFound
	}

	members := model.FindMembers(&model.Member{Email: email})
	spaces := make([]model.Space, len(members))
	for i, m := range members {
		spaces[i] = model.FindSpace(&model.Space{ID: m.Space})
	}
	return c.JSON(http.StatusOK, spaces)
}

// GetSpace returns a space.
func GetSpace(c echo.Context) error {

	id := c.Param("id")

	s := model.FindSpace(&model.Space{ID: id})
	if s.ID == "" {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, s)
}
