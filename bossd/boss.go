package bossd

import (
	"sync"
	"time"

	"zwczou/gobase/container"
	"zwczou/gobase/er"
	em "zwczou/gobase/middleware"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
)

type bossServer struct {
	sync.RWMutex
	opts     *option
	db       *gorm.DB
	redis    *redis.Pool
	echo     *echo.Echo
	startAt  time.Time
	exitChan chan struct{}
}

func NewBossServer(opts *option) *bossServer {
	boss := &bossServer{
		opts:     opts,
		startAt:  time.Now(),
		exitChan: make(chan struct{}),
	}
	return boss
}

func (boss *bossServer) Main() {
	err := boss.init()
	if err != nil {
		log.WithError(err).Fatal("init error")
	}

	opts := boss.opts
	echo := echo.New()
	em.Pprof(echo)
	echo.Validator = em.NewValidator()
	echo.HTTPErrorHandler = er.HTTPErrorHandler
	echo.Use(em.Context())
	echo.Use(em.Hook())
	echo.Use(middleware.Recover())
	echo.Use(middleware.Gzip())

	// 注册静态目录，以及模板
	// 如果纯粹API服务可以注释掉下面这一块
	echo.Group(opts.Static.Path, middleware.Static(opts.Static.Dir))
	renderer, err := em.NewRenderer(opts.Template.Dir, opts.Verbose)
	if err != nil {
		log.WithError(err).Fatal("new renderer error")
	}
	renderer.TplSet.Globals.Update(opts.Template.toPongoCtx())
	echo.Renderer = renderer

	boss.echo = echo
	container.App().Map(echo).Map(renderer).Map(boss.db).Map(boss.redis)

	err = container.Load()
	if err != nil {
		log.WithError(err).Fatal("load extensions error")
	}

	go func() {
		log.WithField("http_addr", opts.HTTPAddr).Info("start http server")
		log.Fatal(echo.Start(opts.HTTPAddr))
	}()
}

func (boss *bossServer) Exit() {
	log.Infof("server exiting")
	container.Exit()

	close(boss.exitChan)
}
