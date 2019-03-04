package admin

import (
	"net/http"
	"strconv"
	"zwczou/boss/model"
	"zwczou/gobase/er"
	"zwczou/gobase/tools"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

type loginViewForm struct {
	Username string `form:"username" validate:"min=5,max=30"`
	Password string `form:"password" validate:"min=6"`
	Remember bool   `form:"remember"`
}

func (as *adminServer) loginView(ctx echo.Context) error {
	template := "admin/login.html"
	isGet := ctx.Request().Method == echo.GET
	if isGet {
		return ctx.Render(http.StatusOK, template, nil)
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
		return ctx.Render(http.StatusOK, template, data)
	}

	next := ctx.QueryParam("next")
	if next == "" {
		next = RedirectUrl
	}
	if as.redirectUrl != nil {
		next = as.redirectUrl()
	}
	token, err := as.genToken(&user)
	if err != nil {
		data = tools.Flash(data, "warning", "服务器内部错误")
		return ctx.Render(http.StatusOK, template, data)
	}
	ctx.SetCookie(&http.Cookie{Name: AdminToken, Value: token, Path: "/"})
	return ctx.Redirect(http.StatusFound, next)
}

type updateUserPasswordViewForm struct {
	UserId  int    `query:"id" form:"id" validate:"min=1"`
	OldPass string `form:"old_pass" validate:"min=6"`
	NewPass string `form:"new_pass" validate:"min=6"`
}

func (as *adminServer) updateUserPasswordView(ctx echo.Context) error {
	template := "admin/users_update_password.html"
	isGet := ctx.Request().Method == echo.GET

	var form updateUserPasswordViewForm
	if err := tools.Bind(ctx, &form); err != nil {
		data := tools.Flash(nil, "warning", "参数错误")
		return ctx.Render(http.StatusOK, template, data)
	}

	var user model.Administrator
	if user.Id == form.UserId || form.UserId == 0 {
		user = *ctx.Get(ContextUser).(*model.Administrator)
	} else {
		err := as.db.Scopes(model.QueryAdministratorScope).First(&user).Error
		if err != nil {
			log.WithError(err).WithField("form", form).Error("query users error")
			data := tools.Flash(nil, "warning", "服务器内部错误")
			return ctx.Render(http.StatusOK, template, data)
		}
	}

	data := pongo2.Context{"user": user, "form": form}
	if isGet {
		return ctx.Render(http.StatusOK, template, data)
	}

	if err := ctx.Validate(&form); err != nil {
		log.WithError(err).Warn("validate error")
		data := tools.Flash(data, "warning", "参数错误")
		return ctx.Render(http.StatusOK, template, data)
	}
	if !user.CheckPassword(form.OldPass) {
		data = tools.Flash(data, "warning", "老密码输入有误")
		return ctx.Render(http.StatusOK, template, data)
	}
	if form.OldPass == form.NewPass {
		data = tools.Flash(data, "success", "密码修改成功!")
		return ctx.Render(http.StatusOK, template, data)
	}

	user.SetPassword(form.NewPass)
	as.db.Select("password", "updated_at").Save(&user)
	data = tools.Flash(data, "success", "密码修改成功!")
	return ctx.Render(http.StatusOK, template, data)
}
