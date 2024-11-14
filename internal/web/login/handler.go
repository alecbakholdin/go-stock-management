package login

import (
	"crypto/subtle"
	"errors"
	"net/http"
	"stock-management/internal/web"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Handler(secret, username, password string) echo.HandlerFunc {
	return func(c echo.Context) error {
		formUsername := c.FormValue("username")
		formPassword := c.FormValue("password")
		c.Logger().Info(username, password, formUsername, formPassword)
		usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(formUsername)) == 1
		passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(formPassword)) == 1

		claims := jwt.MapClaims{
			"user": "admin",
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		}
		if !usernameMatch || !passwordMatch {
			return web.RenderOk(c, LoginForm(formUsername, formPassword, errors.New("Invalid credentials")))
		} else if jwt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret)); err != nil {
			c.Logger().Error(err)
			return web.RenderOk(c, LoginForm(formUsername, formPassword, errors.New("Unexpected error")))
		} else {
			c.Response().Header().Add("Hx-Refresh", "true")
			http.SetCookie(c.Response(), &http.Cookie{
				Name:     "Authorization",
				MaxAge:   7 * 24 * 60 * 60,
				Path:     "/",
				HttpOnly: true,
				Value:    jwt,
			})
			return nil
		}
	}
}
