package easygo

import "net"

const (
	SERVER_TYPE_RPC     = 1 //RPC服务器类型
	SERVER_TYPE_GATEWAY = 2 //gateway服务器类型
	SERVER_StATE_ON     = 1 //服务器状态正常
	SERVER_StATE_OFF    = 2 //服务器状态关闭
)

type IServerInfo interface {
	GetInfo() *ServerInfo
	GetSid() int
	GetName() string
	GetType() int
	GetAddress() string
}

//服务器表 server_info
type ServerInfo struct {
	Sid        int    //服务器编号
	Name       string //服务器名称
	Type       int    //服务器类型
	ExternalIp string //对外ip
	InternalIP string //内部ip
	Port       int    //端口
	State      int    //状态
	ConNum     int    //连接数
	Version    string //版本
}

func NewServerInfo(yaml IYamlConfig) *ServerInfo {
	p := &ServerInfo{}
	p.Init(p, yaml)
	return p
}

func (sSelf *ServerInfo) Init(me IServerInfo, yaml IYamlConfig) {
	sSelf.Sid = YamlCfg.GetValueAsInt("SERVER_ID")
	sSelf.Name = YamlCfg.GetValueAsString("SERVER_NAME")
	sSelf.Type = YamlCfg.GetValueAsInt("SERVER_TYPE")
	sSelf.ExternalIp = YamlCfg.GetValueAsString("SERVER_ADDR")
	sSelf.InternalIP = YamlCfg.GetValueAsString("SERVER_ADDR_INTERNAL")
	sSelf.Port = YamlCfg.GetValueAsInt("LISTEN_PORT_FOR_CLIENT")
	sSelf.State = SERVER_StATE_ON
	sSelf.ConNum = 0
}

func (sSelf *ServerInfo) GetInfo() *ServerInfo {
	return sSelf
}

func (sSelf *ServerInfo) GetSid() int {
	return sSelf.Sid
}

func (sSelf *ServerInfo) GetName() string {
	return sSelf.Name
}

func (sSelf *ServerInfo) GetType() int {
	return sSelf.Type
}

func (sSelf *ServerInfo) GetAddress() string {
	return net.JoinHostPort(sSelf.ExternalIp, AnytoA(sSelf.Port))
}
