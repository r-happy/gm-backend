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

	if space.ParentID != "" {
		// 	parent spaceが有効かどうか確認
		parent_space := model.FindSpace(&model.Space{ID: space.ParentID})
		if parent_space.ID == "" {
			return &echo.HTTPError{
				Code:	400,
				Message: "parent space not found",
			}
		}

		// parent spaceがすでにparentを持ってないか確認
		if parent_space.ParentID != "" {
			return &echo.HTTPError{
				Code:	400,
				Message: "parent space already has a parent",
			}
		}
			// parent spaceに所属するか
		if !IsUserMemberOfSpace(email, space.ParentID) {
			return &echo.HTTPError{
				Code:	400,
				Message: "not a member of parent space",
			}
		}
	}

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

func GetSpace(c echo.Context) error {

    email := userEmailFromToken(c)
    id := c.Param("id")

    // Find the user based on their email
    user := model.FindUser(&model.User{Email: email})
    if user.ID == 0 {
        return echo.ErrNotFound
    }

    // Check if the user is a member of the space
    if !IsUserMemberOfSpace(email, id) {
        // Do nothing if the user is not a member of the space
        return echo.ErrNotFound
    }

    // Find the space
    s := model.FindSpace(&model.Space{ID: id})
    if s.ID == "" {
        return echo.ErrNotFound
    }
    return c.JSON(http.StatusOK, s)
}

func GetChildrens(c echo.Context) error {
	email := userEmailFromToken(c)
    id := c.Param("id")

    // Find the user based on their email
    user := model.FindUser(&model.User{Email: email})
    if user.ID == 0 {
        return echo.ErrNotFound
    }

    // Check if the user is a member of the space
    if !IsUserMemberOfSpace(email, id) {
        // Do nothing if the user is not a member of the space
        return echo.ErrNotFound
    }

	s := model.FindSpace(&model.Space{ID: id})
    if s.ID == "" {
        return echo.ErrNotFound
    }

	spaces := model.FindSpaces(&model.Space{ParentID: s.ID})

	return c.JSON(http.StatusOK, spaces)
}
