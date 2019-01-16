package bossd

import (
	"fmt"
	"strings"
	"time"
	"zwczou/boss/model"
	em "zwczou/gobase/middleware"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/json-iterator/go/extra"
	log "github.com/sirupsen/logrus"
)

func (boss *bossServer) init() error {
	boss.initRedis()
	err := boss.initDatabase()
	if err != nil {
		return err
	}
	return nil
}

// 实现gorm logger接口
type logger struct {
}

func (l logger) Print(vs ...interface{}) {
	var out []string
	for _, v := range vs {
		out = append(out, fmt.Sprint(v))
	}
	if len(vs) > 0 && vs[0] == "log" {
		log.Errorf(strings.Join(out, " "))
		return
	}
	log.Debugf(strings.Join(out, " "))
}

// 初始化数据库
func (boss *bossServer) initDatabase() error {
	dbopts := boss.opts.Database
	db, err := gorm.Open(dbopts.DataType, dbopts.DataSource)
	if err != nil {
		return err
	}
	db.DB().SetMaxOpenConns(dbopts.MaxOpenConns)
	db.DB().SetMaxIdleConns(dbopts.MaxIdleConns)
	db.LogMode(boss.opts.Verbose)
	db.SetLogger(logger{})
	boss.Lock()
	boss.db = db
	boss.Unlock()

	db.AutoMigrate(&model.Administrator{}, &model.Role{}, &model.Permission{}, &model.Menu{})
	db.AutoMigrate(&model.DivideMode{})
	db.AutoMigrate(&model.Card{}, &model.User{})
	db.AutoMigrate(&model.Auth{}, &model.AuthTrack{})
	db.AutoMigrate(&model.Plan{}, &model.CardPlan{}, &model.PartnerPlan{})
	db.AutoMigrate(&model.Order{}, &model.Comission{})
	db.AutoMigrate(&model.Task{}, &model.Remark{})
	return nil
}

func (boss *bossServer) initRedis() {
	redis := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", boss.opts.Redis.Addr, redis.DialDatabase(boss.opts.Redis.DB))
		},
		MaxIdle:     boss.opts.Redis.MaxIdle,
		IdleTimeout: time.Duration(boss.opts.Redis.IdleTimeout) * time.Second,
	}
	boss.Lock()
	boss.redis = redis
	boss.Unlock()
}

func init() {
	extra.RegisterFuzzyDecoders()
	em.SetNamingStrategy(em.LowerCaseWithUnderscores)
	em.RegisterTimeAsFormatCodec("2006-01-02 15:04:05")
}
