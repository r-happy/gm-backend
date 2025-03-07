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

	if !IsUserMemberOfSpace(email, id) {
		return echo.ErrNotFound
	}

	members := model.FindMembers(&model.Member{Space: id})

	if len(members) == 0 {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, members)
}

type AddMemberRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
}

func AddMembers(c echo.Context) error {
	id := c.Param("id")

	// JSONリクエストをバインド
	var req AddMemberRequest
	if err := c.Bind(&req); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		}
	}

	// リクエストのバリデーション
	if req.Email == "" {
		return &echo.HTTPError{
			Code:    400,
			Message: "invalid email",
		}
	}

	// ユーザー存在確認
	if addmember := model.FindUser(&model.User{Email: req.Email}); addmember.ID == 0 {
		return echo.ErrNotFound
	}

	// スペース存在確認
	if space := model.FindSpace(&model.Space{ID: id}); space.ID == "" {
		return echo.ErrNotFound
	}

	email := userEmailFromToken(c)
	user := model.FindUser(&model.User{Email: email})
	if user.ID == 0 {
		return echo.ErrNotFound
	}

	// メンバー作成
	member := &model.Member{
		Space: id,
		Email: req.Email,
		Name:  user.Name,
		Admin: req.Admin,
	}

	model.CreateMember(member)

	return c.JSON(http.StatusCreated, member)
}

func ToggleMemberAdmin(c echo.Context) error {
	id := c.Param("id")

	email := c.FormValue("email")
	myemail := userEmailFromToken(c)

	if email == "" {
		return &echo.HTTPError{
			Code:    400,
			Message: "invalid email",
		}
	}

	// ismeadmin?
	mymemberinfo := model.FindMembers(&model.Member{Email: myemail, Space: id})
	if len(mymemberinfo) == 0 {
		return echo.ErrNotFound
	}
	if !mymemberinfo[0].Admin {
		return echo.ErrForbidden
	}

	if email == myemail {
		return echo.ErrForbidden
	}

	// memberとして有効か, 管理者か確認
	members := model.FindMembers(&model.Member{Email: email, Space: id})
	if len(members) == 0 {
		return echo.ErrNotFound
	}
	member := members[0]
	member.Admin = !member.Admin

	model.SaveMember(&member)

	return c.JSON(http.StatusOK, member)
}
