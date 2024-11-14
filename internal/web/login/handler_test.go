package login

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestInvalidUsername(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("username=invalid&password=validpassword"))
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, Handler("secret", "validusername", "validpassword")(c)) {
		assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
		bytes, err := io.ReadAll(rec.Result().Body)
		assert.NoError(t, err)
		assert.Contains(t, string(bytes), "Invalid credentials")
	}
}

func TestInvalidPassword(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("username=validusername&password=invalid"))
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, Handler("secret", "validusername", "validpassword")(c)) {
		assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
		bytes, err := io.ReadAll(rec.Result().Body)
		assert.NoError(t, err)
		assert.Contains(t, string(bytes), "Invalid credentials")
	}
}

func TestValid(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("username=validusername&password=validpassword"))
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, Handler("secret", "validusername", "validpassword")(c)) {
		assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
		assert.Equal(t, "true", rec.Header().Get("Hx-Refresh"))
		cookies := rec.Result().Cookies()
		if assert.Equal(t, 1, len(cookies)) {
			cookie := cookies[0]
			assert.Equal(t, "Authorization", cookie.Name)
			_, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
				return []byte("secret"), nil
			})
			assert.NoError(t, err)
		}
	}
}
