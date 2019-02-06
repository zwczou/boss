package admin

import (
	"net/http"

	"github.com/labstack/echo"
)

func (as *adminServer) loginView(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "admin/login.html", nil)
}
