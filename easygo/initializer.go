package easygo

/*
公共初始化
*/

type IInitializer interface {
	CreateRedisMgr() IRedisManager
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

func (initSelf *Initializer) Execute() {
	//redis数据库
	RedisMgr = initSelf.Me.CreateRedisMgr()
	ETCD_CENTER_ADDR = []string{"127.0.0.1:2379"}
	ETCD_SERVER_PATH = "/grpc_gateway/"
}

func (initSelf *Initializer) CreateRedisMgr() IRedisManager {
	return NewRedisManager()
}

var (
	RedisMgr IRedisManager //redis存储
)
