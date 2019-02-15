package admin

import (
	"net/http"
	"strconv"
	"zwczou/boss/model"
	"zwczou/gobase/er"
	"zwczou/gobase/tools"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
)

type loginViewForm struct {
	Username string `form:"username" validate:"min=5,max=30"`
	Password string `form:"password" validate:"min=6"`
	Remember bool   `form:"remember"`
}

func (as *adminServer) loginView(ctx echo.Context) error {
	if ctx.Request().Method == echo.GET {
		return ctx.Render(http.StatusOK, "admin/login.html", nil)
	}

	var form loginViewForm
	if err := tools.Validate(ctx, &form); err != nil {
		return er.ErrInvalidArgs
	}
	db := as.db

	var user model.Administrator
	if _, err := strconv.Atoi(form.Username); err == nil && len(form.Username) == 11 {
		db.Where("mobile = ?", form.Username).First(&user)
	} else {
		db.Where("nickname = ?", form.Username).First(&user)
	}

	data := pongo2.Context{"form": form}
	if user.Id == 0 || !user.CheckPassword(form.Password) {
		data = tools.Flash(data, "warning", "用户名或密码错误")
		return ctx.Render(http.StatusOK, "admin/login.html", data)
	}

	next := ctx.QueryParam("next")
	if next == "" {
		next = RedirectUrl
	}
	token, err := as.genToken(&user)
	if err != nil {
		data = tools.Flash(data, "warning", "服务器内部错误")
		return ctx.Render(http.StatusOK, "admin/login.html", data)
	}
	ctx.SetCookie(&http.Cookie{Name: AdminToken, Value: token, Path: "/"})
	return ctx.Redirect(http.StatusFound, next)
}
