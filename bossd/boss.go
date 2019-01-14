package bossd

import (
	"net/http"
	"sync"
	"time"

	"zwczou/gobase/container"
	"zwczou/gobase/er"
	em "zwczou/gobase/middleware"

	"github.com/facebookgo/grace/gracehttp"
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

	echo.NotFoundHandler = er.NotFoundHandler
	echo.MethodNotAllowedHandler = er.MethodNotAllowedHandler

	opts := boss.opts
	e := echo.New()
	e.Debug = opts.Verbose
	em.Pprof(e)
	e.Validator = em.NewValidator()
	e.HTTPErrorHandler = er.HTTPErrorHandler
	e.Use(em.Context())
	e.Use(em.Hook())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	// 注册静态目录，以及模板
	// 如果纯粹API服务可以注释掉下面这一块
	e.Group(opts.Static.Path, middleware.Static(opts.Static.Dir))
	renderer, err := em.NewRenderer(opts.Template.Dir, opts.Verbose)
	if err != nil {
		log.WithError(err).Fatal("new renderer error")
	}
	renderer.TplSet.Globals.Update(opts.Template.toPongoCtx())
	e.Renderer = renderer

	boss.echo = e
	container.App().Map(e).Map(renderer).Map(boss.db).Map(boss.redis)

	err = container.Load()
	if err != nil {
		log.WithError(err).Fatal("load extensions error")
	}

	go func() {
		log.WithField("http_addr", opts.HTTPAddr).Info("start http server")
		gracehttp.Serve(&http.Server{Addr: opts.HTTPAddr, Handler: e})
	}()
}

func (boss *bossServer) Exit() {
	log.Infof("server exiting")
	close(boss.exitChan)
	container.Exit()
}
