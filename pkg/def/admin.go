package def

import (
	"github.com/labstack/echo"
)

type CheckLoginFunc func(bool) echo.MiddlewareFunc

func (cl CheckLoginFunc) CheckLogin(isJson bool) echo.MiddlewareFunc {
	return cl(isJson)
}
