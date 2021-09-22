package easygo

import (
	"github.com/garyburd/redigo/redis"
)

/*
redis后台管理员管理器
*/
const (
	ADMIN_ONLINE_LIST = "admin_online_list"
)

type RedisAdmin struct {
	UserId    int64  //ID
	Role      int32  //管理员类型 0超管，1管理员
	ServerId  int32  //服务器Id
	Timestamp int64  //登录时间戳
	Token     string //登录token
}

//添加管理员到Redis
func SetRedisAdmin(obj *RedisAdmin) {
	//js, err := json.Marshal(obj)
	//PanicError(err)
	err1 := RedisMgr.GetC().HMSet(MakeRedisKey(ADMIN_ONLINE_LIST, obj.UserId), obj)
	PanicError(err1)
}

//查询管理员列表
func GetRedisAdminList() []*RedisAdmin {
	var lst []*RedisAdmin
	keys, err := RedisMgr.GetC().Scan(ADMIN_ONLINE_LIST)
	if err != nil {
		return lst
	}
	for _, key := range keys {
		value, err := RedisMgr.GetC().HGetAll(key)
		if err != nil {
			continue
		}
		var admin RedisAdmin
		err = redis.ScanStruct(value, &admin)
		if err != nil {
			continue
		}
		lst = append(lst, &admin)
	}
	return lst
}

//根据id查询
func GetRedisAdmin(UserId int64) *RedisAdmin {
	val, err := RedisMgr.GetC().HGetAll(MakeRedisKey(ADMIN_ONLINE_LIST, UserId))
	if err != nil {
		return nil
	}

	var admin RedisAdmin
	err = redis.ScanStruct(val, &admin)
	if err != nil {
		return nil
	}
	return &admin
}

func GetRedisAdmin2(UserId int64, Fild string) int64 {
	val, err := RedisMgr.GetC().HGet(MakeRedisKey(ADMIN_ONLINE_LIST, UserId), Fild)
	if err != nil {
		return 0
	}

	return AtoInt64(string(val))
}

//删除Redis中的数据
func DelRedisAdmin(id int64) {
	_, err := RedisMgr.GetC().Delete(MakeRedisKey(ADMIN_ONLINE_LIST, id))
	PanicError(err)
}

func SetOne(fild string, v int) {
	err1 := RedisMgr.GetC().Set(fild, v)
	PanicError(err1)
}
