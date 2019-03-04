package admin

import (
	"sync"
	"zwczou/boss/pkg/def"
	"zwczou/gobase/container"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type adminServer struct {
	sync.RWMutex
	name        string
	db          *gorm.DB
	redis       *redis.Pool
	echo        *echo.Echo
	redirectUrl func() string
	exitChan    chan struct{}
}

func newAdminServer() *adminServer {
	s := &adminServer{
		name:     "admin.system",
		exitChan: make(chan struct{}),
	}
	return s
}

func (as *adminServer) Name() string {
	return as.name
}

// 通过反射赋值需要的数据
func (as *adminServer) assign(db *gorm.DB, redis *redis.Pool, echo *echo.Echo, redirectUrl func() string) {
	as.db = db
	as.redis = redis
	as.echo = echo
	as.redirectUrl = redirectUrl
}

func (as *adminServer) Load(app *container.Container) error {
	app.Map(def.CheckLoginFunc(as.CheckLogin))
	_, err := app.Invoke(as.assign)
	if err != nil {
		return err
	}

	admin := as.echo.Group("/admin")
	{
		admin.GET("/login", as.loginView)
		admin.POST("/login", as.loginView)
	}

	return nil
}

func (as *adminServer) Exit() {
	close(as.exitChan)
}

func init() {
	container.Pre(newAdminServer())
}
