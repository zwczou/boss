package admin

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"zwczou/boss/model"
	"zwczou/gobase/er"
	md "zwczou/gobase/middleware"

	"github.com/flosch/pongo2"
	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

// 生成token
func (as *adminServer) genToken(user *model.Administrator) (token string, err error) {
	rc := as.redis.Get()
	defer rc.Close()
	ut := newUserToken(rc, TokenExpires)
	val := newUserTokenValue(user.Id)
	token = strings.NewReplacer("-", "").Replace(uuid.New().String())
	err = ut.Set(token, val)
	return
}

// 解析token
func (as *adminServer) parseToken(ctx echo.Context) (userId int, err error) {
	tokenCookie, err := ctx.Request().Cookie(AdminToken)
	if err != nil {
		return 0, err
	}
	token := tokenCookie.Value
	rc := as.redis.Get()
	defer rc.Close()
	ut := newUserToken(rc, TokenExpires)
	val, ok := ut.Check(token)
	ctx.Set(ContextToken, token)
	if !ok {
		return 0, ErrUnauthenticated
	}
	return val.UserId, nil
}

func (as *adminServer) CheckLogin(isJson bool) echo.MiddlewareFunc {
	renderer := as.echo.Renderer.(*md.Renderer)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			redirectUrl := fmt.Sprintf("/admin/login?next=%s", url.QueryEscape(ctx.Request().URL.String()))
			userId, err := as.parseToken(ctx)
			if err != nil {
				if isJson {
					return er.ErrUnauthorized
				}
				return ctx.Redirect(http.StatusFound, redirectUrl)
			}

			var user model.Administrator
			fields := log.Fields{
				"user_id": userId,
			}
			user.Id = userId
			if err := as.db.First(&user).Error; err != nil {
				fields["error"] = err
				log.WithFields(fields).Warn("find user error")
				if isJson {
					return er.ErrUnauthorized
				}
				return ctx.Redirect(http.StatusFound, redirectUrl)
			}

			// 校验权限
			if !user.Check(ctx.Request()) {
				if isJson {
					return er.ErrForbidden
				}
				return ctx.HTML(http.StatusForbidden, "<h1>403 forbidden</h1>")
			}

			ctx.Set(ContextUserId, userId)
			ctx.Set(ContextUser, &user)
			data := pongo2.Context{"_user": &user, "ctx": ctx, "req": ctx.Request()}
			renderer.TplSet.Globals.Update(data)
			return next(ctx)
		}
	}
}
