package admin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

type userTokenValue struct {
	UserId    int
	ExpiredAt time.Time
}

func newUserTokenValue(userId int) *userTokenValue {
	return &userTokenValue{
		UserId: userId,
	}
}

// 用户令牌
type userToken struct {
	conn         redis.Conn
	expires      int
	tokenListKey string
	tokenKey     string
}

func newUserToken(conn redis.Conn, expires int) *userToken {
	return &userToken{
		conn:         conn,
		expires:      expires,
		tokenListKey: "admin:user:token:list:%d",
		tokenKey:     "admin:user:token",
	}
}

func (ut *userToken) listKey(userId int) string {
	return fmt.Sprintf(ut.tokenListKey, userId)
}

func (ut *userToken) key(token string) string {
	return ut.tokenKey
}

func (ut *userToken) Set(token string, val *userTokenValue) error {
	listKey := ut.listKey(val.UserId)
	key := ut.key(token)

	if ut.expires > 0 {
		val.ExpiredAt = time.Now().Add(time.Second * time.Duration(ut.expires))
	}
	body, err := json.Marshal(val)
	if err != nil {
		return err
	}

	_, err = ut.conn.Do("LPUSH", listKey, token)
	if err != nil {
		return err
	}
	_, err = ut.conn.Do("HSET", key, token, body)
	return err
}

// 删除用户的历史token
// 比如当修改密码的时候，所有设备重新登录
func (ut *userToken) Delete(userId int) error {
	listKey := ut.listKey(userId)

	tokens, err := redis.Strings(ut.conn.Do("LRANGE", listKey, 0, -1))
	if err != nil {
		return err
	}
	ut.conn.Send("MULTI")
	ut.conn.Send("DEL", listKey)
	for _, token := range tokens {
		key := ut.key(token)
		ut.conn.Send("HDEL", key, token)
	}
	_, err = ut.conn.Do("EXEC")
	return err
}

func (ut *userToken) Get(token string) (val *userTokenValue, err error) {
	val = &userTokenValue{}
	key := ut.key(token)
	body, err := redis.Bytes(ut.conn.Do("HGET", key, token))
	if err != nil {
		return val, err
	}
	err = json.Unmarshal(body, val)
	return
}

// 获取token并且检测token是否有效
func (ut *userToken) Check(token string) (val *userTokenValue, ok bool) {
	value, err := ut.Get(token)
	if err != nil {
		log.WithError(err).WithField("token", token).Warnf("get token error")
		return nil, false
	}
	if time.Now().Sub(value.ExpiredAt) > 0 {
		return nil, false
	}
	return value, true
}

// 删除单个token
func (ut *userToken) DeleteToken(userId int, token string) error {
	listKey := ut.listKey(userId)
	key := ut.key(token)
	_, err := ut.conn.Do("LREM", listKey, 0, token)
	if err != nil {
		return err
	}
	_, err = ut.conn.Do("HDEL", key, token)
	return err
}
