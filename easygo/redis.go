package easygo

import (
	"github.com/astaxie/beego/logs"
)

const (
	REDIS_SERVER_ADDR     = "127.0.0.1:6379" //redis服务器地址
	REDIS_SERVER_PASS     = ""               //redis服务器密码
	REDIS_SERVER_DATABASE = 0                //redis数据库序号
)

type IRedisManager interface {
	GetC() RedisCache
}

type RedisManager struct {
	Me IRedisManager
	RedisCache
	Mutex RLock
}

func NewRedisManager() *RedisManager { // services map[string]interface{},
	p := &RedisManager{}
	p.Init(p)
	return p
}

//初始化
func (rmSelf *RedisManager) Init(me IRedisManager) {
	rmSelf.Me = me
	host := REDIS_SERVER_ADDR
	pass := REDIS_SERVER_PASS
	db := REDIS_SERVER_DATABASE
	rmSelf.RedisCache = NewRedisCache(db, host, REDIS_DEFAULT, pass)
	logs.Info("Successfully connected to the Redis server %s", host)

}

func (rmSelf *RedisManager) GetC() RedisCache {
	rmSelf.Mutex.Lock()
	defer rmSelf.Mutex.Unlock()
	return rmSelf.RedisCache
}
