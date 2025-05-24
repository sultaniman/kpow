package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RenderForm(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "form.html", map[string]any{
		"name": "HOME",
		"msg":  "Hello, Boatswain!",
	})
}
