package admin

import (
	"strings"
	"zwczou/boss/model"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

// 生成token
func (as *adminServer) genToken(user *model.User) (token string, err error) {
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
