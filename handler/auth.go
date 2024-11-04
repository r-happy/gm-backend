package handler

import (
	"back/model"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/badoux/checkmail"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type jwtCustomClaims struct {
	UID   int    `json:"uid"`
	Email string `json:"email"`
	jwt.StandardClaims
}

var Config middleware.JWTConfig

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load signing key from .env file
	signingKey := os.Getenv("JWT_KEY")
	if signingKey == "" {
		log.Fatal("JWT_KEY is not set in .env file")
	}

	// Set up JWT middleware config
	Config = middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte(signingKey),
	}
}

func CreateUser(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	if user.Name == "" || user.Password == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid name or password",
		}
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid email",
		}
	}

	if u := model.FindUser(&model.User{Email: user.Email}); u.ID != 0 {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "this email is already used",
		}
	}

	hashedPassword := sha256.Sum256([]byte(user.Password))
	user.Password = fmt.Sprintf("%x", hashedPassword)

	model.CreateUser(user)
	user.Password = ""

	return c.JSON(http.StatusCreated, user)
}

func Login(c echo.Context) error {
	u := new(model.User)
	if err := c.Bind(u); err != nil {
		return err
	}

	user := model.FindUser(&model.User{Email: u.Email})

	hashedPassword := sha256.Sum256([]byte(u.Password))
	u.Password = fmt.Sprintf("%x", hashedPassword)

	if user.Password != u.Password {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "invalid name or password",
		}
	}

	claims := &jwtCustomClaims{
		UID:   user.ID,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func userEmailFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	return claims.Email
}
