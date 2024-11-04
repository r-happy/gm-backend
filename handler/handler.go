package handler

import (
	"back/model"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

// generateUniqueID generates a unique ID for a space.
func generateUniqueID() (string, error) {
	for {
		tmpUUID := uuid.NewString()
		s := model.FindSpace(&model.Space{ID: tmpUUID})
		if s.ID == "" {
			return tmpUUID, nil
		}
	}
}

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

// GetMembers returns all members of a space.
func GetMembers(c echo.Context) error {
	email := userEmailFromToken(c)
	id := c.Param("id")
	user := model.FindUser(&model.User{Email: email})
	if user.ID == 0 {
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

