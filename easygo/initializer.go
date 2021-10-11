package easygo

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/sasha-s/go-deadlock"
	"go.uber.org/zap"
)

/*
公共初始化
*/

type IInitializer interface {
	CreateYamlConfig(dict KWAT) IYamlConfig
	SetDeadLockOptions()
	SetBeegoLogs(dict KWAT)
	GetBeeLogger() *logs.BeeLogger
	SetZapLogger()

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
	initSelf.Me.SetZapLogger()

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

//beego日志设置
func (initSelf *Initializer) SetBeegoLogs(dict KWAT) {
	logName, ok := dict["logName"]
	if !ok {
		panic("我需要一个 logName 参数")
	}

	var err error
	err = logs.SetLogger(logs.AdapterConsole, fmt.Sprintf(`{"level":%d}`, logs.LevelDebug))
	PanicError(err)

	/*
		filename 保存的文件名
		maxlines 每个文件保存的最大行数，默认值 1000000
		maxsize 每个文件保存的最大尺寸，默认值是 1 << 28, //256 MB
		daily 是否按照每天 logrotate，默认是 true
		maxdays 文件最多保存多少天，默认保存 7 天
		rotate 是否开启 logrotate，默认是 true
		level 日志保存的时候的级别，默认是 Trace 级别
		perm 日志文件权限
		separate 需要单独写入文件的日志级别,设置后命名类似 test.error.log
	*/
	config := `{"filename":"logs/%s.log", "separate":["error", "warning", "info", "debug"], "level":%d, "daily":true, "rotate":true, "maxdays":15,"maxlines":1000000,"perm":"777"}`
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

//zap日志设置
func (initSelf *Initializer) SetZapLogger() {
	// file, _ := os.Create("logs/test.log")
	// writeSyncer := zapcore.AddSync(file)

	// encoderConfig := zap.NewProductionEncoderConfig()
	// encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// encoder := zapcore.NewJSONEncoder(encoderConfig)
	// core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	// logger := zap.New(core, zap.AddCaller())
	logger, _ := zap.NewProduction(zap.AddCaller())
	Logs = logger.Sugar()
}

var (
	YamlCfg   IYamlConfig        //配置
	RedisMgr  IRedisManager      //redis存储
	EtcdMgr   IClient3KVManager  //etcd存储
	ServerMgr IServerManager     //serverMgr存储
	PServer   IServerInfo        //server存储
	Logs      *zap.SugaredLogger //zap日志
)
