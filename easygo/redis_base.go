package easygo

import (
	"encoding/json"

	"time"

	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

/*
redis数据管理基类模板
*/

const REDIS_SAVE_TIME = 600 * time.Second //保存时间

const REDIS_SAVE_KEY_LIST = "save_keys"   //需要保存的keys集合
const REDIS_EXIST_KEY_LIST = "exist_keys" //需要保存的keys集合

type IRedisBase interface {
	GetId() interface{}
	GetKeyId() string
	UpdateData()
	GetRedisSaveData() interface{}
	InitRedis()
	SaveOtherData() //保存其他数据
}
type RedisBase struct {
	Me         IRedisBase
	TBName     string //表名 redis key名并且数据库表名
	ExistKey   string //reidis现存key值d
	SaveKeys   string //存储列表
	CreateTime int64  //对象创建时间
	//SaveStatus bool         //存储状态:true，需要存储，false不需要存储
	Mutex   RLock //局部数据锁,只对同一服务器内有用
	IncrId  int64
	Sid     int32
	IsCheck bool
}

func (rbSelf *RedisBase) Init(me IRedisBase, id interface{}, tbName string) {
	rbSelf.TBName = tbName
	rbSelf.CreateTime = GetMillSecond()
	//rbSelf.SaveStatus = false //默认不需要存储
	rbSelf.IsCheck = true
	rbSelf.Me = me
	rbSelf.SaveKeys = MakeRedisKey(REDIS_SAVE_KEY_LIST, tbName)
	rbSelf.ExistKey = MakeRedisKey(REDIS_EXIST_KEY_LIST, tbName)
}

//组装redis Key
func MakeRedisKey(keys ...interface{}) string {
	res := ""
	for k, v := range keys {
		res += AnytoA(v)
		if k < (len(keys) - 1) {
			res += ":"
		}
	}
	return res
}

//增加现存keys
func (rbSelf *RedisBase) AddToExistList(id interface{}) {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()
	err := RedisMgr.GetC().HSet(rbSelf.ExistKey, AnytoA(id), rbSelf.Sid)
	PanicError(err)
}

//删除key
func (rbSelf *RedisBase) DelToExistList(id interface{}) {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()
	_, err := RedisMgr.GetC().Hdel(rbSelf.ExistKey, AnytoA(id))
	PanicError(err)
}

//增加到需要保存列表
func (rbSelf *RedisBase) AddToSaveList() {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()
	err := RedisMgr.GetC().SAdd(rbSelf.SaveKeys, rbSelf.Me.GetId())
	PanicError(err)
}

//从需保存列表删除key
func (rbSelf *RedisBase) DelToSaveList() {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()
	err := RedisMgr.GetC().SRem(rbSelf.SaveKeys, rbSelf.Me.GetId())
	PanicError(err)
}

//获取所有要保存的Keys
func GetAllRedisSaveList(key string, data interface{}) {
	val, err := RedisMgr.GetC().Smembers(MakeRedisKey(REDIS_SAVE_KEY_LIST, key))
	PanicError(err)
	switch value := data.(type) {
	case *[]string:
		InterfersToStrings(val, value)
		data = value
	case *[]int64:
		InterfersToInt64s(val, value)
		data = value
	case *[]int32:
		InterfersToInt32s(val, value)
		data = value
	default:
		panic("找不到类型，请自行定义")
	}
}
func GetAllRedisExistList(key string, data interface{}) {
	val, err := RedisMgr.GetC().HKeys(MakeRedisKey(REDIS_EXIST_KEY_LIST, key))
	PanicError(err)
	switch value := data.(type) {
	case *[]string:
		InterfersToStrings(val, value)
		data = value
	case *[]int64:
		InterfersToInt64s(val, value)
		data = value
	case *[]int32:
		InterfersToInt32s(val, value)
		data = value
	default:
		panic("找不到类型，请自行定义")
	}
}

func (rbSelf *RedisBase) SetSaveSid() {
	rbSelf.AddToExistList(rbSelf.Me.GetId())
}
func (rbSelf *RedisBase) CheckIsDelRedisKey() bool {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()
	if !rbSelf.IsExistKey() {
		return false
	}
	sid, err := redis.Int(RedisMgr.GetC().HGet(rbSelf.ExistKey, AnytoA(rbSelf.Me.GetId())))
	if err != nil {
		logs.Error("CheckIsDelRedisKey 报错:", err.Error())
		return false
	}

	//logs.Info("key:", sid, rbSelf.Sid)
	return int32(sid) == rbSelf.Sid
}

//设置数据是否需要存储
func (rbSelf *RedisBase) SetSaveStatus(b bool) {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()
	if b {
		exist := RedisMgr.GetC().SIsMember(rbSelf.SaveKeys, rbSelf.Me.GetId())
		if !exist {
			rbSelf.AddToSaveList()
		}
	} else {
		rbSelf.DelToSaveList()
	}
}

//存储
func (rbSelf *RedisBase) GetSaveStatus() bool {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()
	b := RedisMgr.GetC().SIsMember(rbSelf.SaveKeys, rbSelf.Me.GetId())
	return b
}

//指定值增加
func (rbSelf *RedisBase) IncrOneValue(key string, val int64) int64 {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()
	if !rbSelf.IsExistKey() {
		//如果key值不存在，先获取
		rbSelf.Me.InitRedis()
	}
	res := RedisMgr.GetC().HIncrBy(rbSelf.Me.GetKeyId(), key, val)
	rbSelf.SetSaveStatus(true)
	rbSelf.CreateTime = GetMillSecond()
	rbSelf.SetSaveSid()
	return res

}

//修改指定值
func (rbSelf *RedisBase) SetOneValue(key string, val interface{}) {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()
	if !rbSelf.IsExistKey() {
		//如果key值不存在，先获取
		rbSelf.Me.InitRedis()
	}
	err := RedisMgr.GetC().HSet(rbSelf.Me.GetKeyId(), key, val)
	PanicError(err)
	rbSelf.SetSaveStatus(true)
	rbSelf.CreateTime = GetMillSecond()
	rbSelf.SetSaveSid()
}

//获取指定值
func (rbSelf *RedisBase) GetOneValue(key string, val interface{}) {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()
	if !rbSelf.IsExistKey() {
		//如果key值不存在，先获取
		rbSelf.Me.InitRedis()
	}
	res, err := RedisMgr.GetC().HMGet(rbSelf.Me.GetKeyId(), key)
	PanicError(err)
	_, err = redis.Scan(res, val)
	PanicError(err)
	rbSelf.CreateTime = GetMillSecond()
	rbSelf.SetSaveSid()

}

//ridis存储指定hash值
func (rbSelf *RedisBase) SetStringValueToRedis(name string, data interface{}, save ...bool) {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()
	if !rbSelf.IsExistKey() {
		//如果key值不存在，先获取
		rbSelf.Me.InitRedis()
	}
	//logs.Info("设置玩家数据:", name, data)
	val, err := json.Marshal(data)
	PanicError(err)
	err = RedisMgr.GetC().HSet(rbSelf.Me.GetKeyId(), name, string(val))
	PanicError(err)
	isSave := append(save, true)[0]
	if isSave {
		rbSelf.SetSaveStatus(true)
	}
	rbSelf.CreateTime = GetMillSecond()
	rbSelf.SetSaveSid()
}

//redis获取指定hash值
func (rbSelf *RedisBase) GetStringValueToRedis(name string, data interface{}) {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()
	if !rbSelf.IsExistKey() {
		//如果key值不存在，先获取
		rbSelf.Me.InitRedis()
	}
	val, err := RedisMgr.GetC().HGet(rbSelf.Me.GetKeyId(), name)
	if len(val) == 0 {
		return
	}
	PanicError(err)
	err = json.Unmarshal(val, &data)
	PanicError(err)
	rbSelf.CreateTime = GetMillSecond()
	rbSelf.SetSaveSid()
}

//从mongo中读取数据
func (rbSelf *RedisBase) QueryDBData(id interface{}) interface{} {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()

	return nil
}

//检测是否存在reids数据
func (rbSelf *RedisBase) IsExistKey() bool {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()
	if rbSelf.IsCheck {
		res, err := RedisMgr.GetC().Exist(rbSelf.Me.GetKeyId())
		PanicError(err)
		return res
	}
	return true
}

//全部数据写到DB
func (rbSelf *RedisBase) SaveToDB() {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()

	if rbSelf.IsExistKey() { //只有redis 存在key才进行存储
		data := rbSelf.Me.GetRedisSaveData()
		if data == nil {
			logs.Error("SaveToMongo 存储时获取到数据为空:", rbSelf.Me.GetKeyId())
			return
		}
		rbSelf.Me.SaveOtherData()
		//logs.Info("redis 存储 ：", rbSelf.Me.GetKeyId(), data)

		go func() {

			//save to db
		}()

		rbSelf.SetSaveStatus(false)
	}

}

//保存单个字段
func (rbSelf *RedisBase) SaveOneRedisDataToDB(file string, value interface{}) {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()
	//save to db
}

//把key从redis删除
func (rbSelf *RedisBase) DelRedisKey() {
	rbSelf.Mutex.Lock()
	defer rbSelf.Mutex.Unlock()
	_, err := RedisMgr.GetC().Delete(rbSelf.Me.GetKeyId())
	PanicError(err)
}

//对外通用方法
func DelAllKeyFromRedis(dbName, existName string, val interface{}) {
	var delKeys []interface{}
	switch ids := val.(type) {
	case []int64:
		for _, id := range ids {
			delKeys = append(delKeys, MakeRedisKey(dbName, id))
		}
	case []string:
		for _, id := range ids {
			delKeys = append(delKeys, MakeRedisKey(dbName, id))
		}
	}
	delKeys = append(delKeys, existName)
	if len(delKeys) > 0 {
		_, err := RedisMgr.GetC().Delete(delKeys...)
		PanicError(err)
	}
}

//:TODO 后续优化分布式锁
func (rbSelf *RedisBase) Lock() {
	_, err := RedisMgr.GetC().Delete(rbSelf.Me.GetKeyId())
	PanicError(err)
	logs.Info("删除redis key:", rbSelf.Me.GetKeyId())
}
func (rbSelf *RedisBase) UnLock() {
	_, err := RedisMgr.GetC().Delete(rbSelf.Me.GetKeyId())
	PanicError(err)
	logs.Info("删除redis key:", rbSelf.Me.GetKeyId())
}
func (rbSelf *RedisBase) CheckLock() {
	_, err := RedisMgr.GetC().Delete(rbSelf.Me.GetKeyId())
	PanicError(err)
	logs.Info("删除redis key:", rbSelf.Me.GetKeyId())
}
