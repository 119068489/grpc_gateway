package easygo

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/sasha-s/go-deadlock"
)

/*
公共初始化
*/

type IInitializer interface {
	CreateYamlConfig(dict KWAT) IYamlConfig
	SetDeadLockOptions()
	SetBeegoLogs(dict KWAT)
	GetBeeLogger() *logs.BeeLogger

	CreateServerInfo(yaml IYamlConfig) IServerInfo
	CreateServerInfoMgr() IServerManager
	CreateRedisMgr(yaml IYamlConfig) IRedisManager
	CreateEtcd(yaml IYamlConfig) IClient3KVManager
}

type Initializer struct {
	Me IInitializer
	// 此类不要带状态
}

func NewInitializer() *Initializer {
	p := &Initializer{}
	p.Init(p)
	return p
}

func (initSelf *Initializer) Init(me IInitializer) {
	initSelf.Me = me
}

//执行公共配置初始化
func (initSelf *Initializer) Execute(dict KWAT) {
	initSelf.Me.SetBeegoLogs(dict)
	initSelf.Me.SetDeadLockOptions()

	YamlCfg = initSelf.Me.CreateYamlConfig(dict)
	EDITION = YamlCfg.GetValueAsString("EDITION")
	IS_FORMAL_SERVER = YamlCfg.GetValueAsBool("IS_FORMAL_SERVER")
	IS_FORMAL_SERVER = YamlCfg.GetValueAsBool("IS_FORMAL_SERVER")
	IS_TFSERVER = YamlCfg.GetValueAsBool("IS_TFSERVER")

	//服务器管理器
	ServerMgr = initSelf.Me.CreateServerInfoMgr()
	//服务器信息
	PServer = initSelf.Me.CreateServerInfo(YamlCfg)
	SERVER_ADDR = PServer.GetAddress()
	SERVER_NAME = PServer.GetName()
	SERVER_ID = PServer.GetSid()

	//redis数据库
	RedisMgr = initSelf.Me.CreateRedisMgr(YamlCfg)

	//etcd
	EtcdMgr = initSelf.Me.CreateEtcd(YamlCfg)
	EtcdMgr.StartClintTV3()
}

func (initSelf *Initializer) CreateYamlConfig(dict KWAT) IYamlConfig {
	path, ok := dict["yamlPath"]
	if !ok {
		panic("我需要一个 yamlPath 参数")
	}
	return NewYamlConfig(path.(string))
}

func (initSelf *Initializer) CreateServerInfo(yaml IYamlConfig) IServerInfo {
	return NewServerInfo(yaml)
}

func (initSelf *Initializer) CreateServerInfoMgr() IServerManager {
	return NewServerInfoManager()
}

func (initSelf *Initializer) CreateRedisMgr(yaml IYamlConfig) IRedisManager {
	return NewRedisManager(yaml)
}

func (initSelf *Initializer) CreateEtcd(yaml IYamlConfig) IClient3KVManager {
	return NewClient3KVManager(yaml)
}

func (initSelf *Initializer) SetBeegoLogs(dict KWAT) {
	logName, ok := dict["logName"]
	if !ok {
		panic("我需要一个 logName 参数")
	}

	var err error
	err = logs.SetLogger(logs.AdapterConsole, fmt.Sprintf(`{"level":%d}`, logs.LevelDebug))
	PanicError(err)

	config := `{"filename":"logs/%s.log", "separate":["error", "warning", "info", "debug"], "level":%d, "daily":true, "rotate":true, "maxdays":15,"perm":"777"}`
	config = fmt.Sprintf(config, logName, logs.LevelDebug)
	err = logs.SetLogger(logs.AdapterMultiFile, config)
	PanicError(err)

	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
	logs.Async()
}

func (initSelf *Initializer) GetBeeLogger() *logs.BeeLogger {
	// TODO: 不要使用 logs 的包变量 beeLogger，要 new 一个新的 loger 出来
	return logs.GetBeeLogger()
}

func (initSelf *Initializer) SetDeadLockOptions() {
	// deadlock.Opts.Disable = true // 死锁检查 + 超时检查
	// deadlock.Opts.DisableLockOrderDetection = true // 死锁检查
	deadlock.Opts.DeadlockTimeout = 0             // 默认 30 秒,0 表示不检查 // 4 * time.Second
	deadlock.Opts.OnPotentialDeadlock = func() {} // 默认 exit()
}

var (
	YamlCfg   IYamlConfig
	RedisMgr  IRedisManager     //redis存储
	EtcdMgr   IClient3KVManager //etcd存储
	ServerMgr IServerManager    //serverMgr存储
	PServer   IServerInfo       //server存储
)
