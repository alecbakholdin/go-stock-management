package root

import (
	"stock-management/internal/task"
	"stock-management/internal/web"

	"github.com/labstack/echo/v4"
)

func Handler(tasks []task.TaskStatus) echo.HandlerFunc {
	return func(c echo.Context) error {
		return web.RenderOk(c, Root(tasks))
	}
}
