package root

import (
	"stock-management/internal/web"

	"github.com/labstack/echo/v4"
)

func Handler(c echo.Context) error {
	return web.RenderOk(c, Root())
}
