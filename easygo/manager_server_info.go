package easygo

import (
	"sort"
)

type (
	server_id = int32
)

//服务器管理器，管理连接服务器信息
type ServerInfoManager struct {
	ServerInfo map[server_id]*ServerInfo
	Mutex      RLock
}

func NewServerInfoManager() *ServerInfoManager { // services map[string]interface{},
	p := &ServerInfoManager{}
	p.Init()
	return p
}

func (serverSelf *ServerInfoManager) Init() {
	serverSelf.ServerInfo = make(map[server_id]*ServerInfo)
}
func (serverSelf *ServerInfoManager) AddServerInfo(srv *ServerInfo) {
	serverSelf.Mutex.Lock()
	defer serverSelf.Mutex.Unlock()
	serverSelf.ServerInfo[srv.Sid] = srv
}
func (serverSelf *ServerInfoManager) DelServerInfo(id server_id) {
	serverSelf.Mutex.Lock()
	defer serverSelf.Mutex.Unlock()
	delete(serverSelf.ServerInfo, id)
}
func (serverSelf *ServerInfoManager) GetServerInfo(serverId server_id) *ServerInfo {
	serverSelf.Mutex.Lock()
	defer serverSelf.Mutex.Unlock()

	srvInfo, ok := serverSelf.ServerInfo[serverId]
	if !ok {
		return nil
	}
	return srvInfo
}
func (serverSelf *ServerInfoManager) ChangeServerState(serverId server_id, st int32) {
	serverSelf.Mutex.Lock()
	defer serverSelf.Mutex.Unlock()
	srv := serverSelf.GetServerInfo(serverId)
	srv.State = st
}

//负载均衡，分配一台服务器
func (serverSelf *ServerInfoManager) GetIdelServer(t int32) *ServerInfo {
	serverSelf.Mutex.Lock()
	defer serverSelf.Mutex.Unlock()
	temMap := make(map[int32]int32, len(serverSelf.ServerInfo))
	for k, v := range serverSelf.ServerInfo {
		if v.Type == t {
			temMap[k] = v.ConNum
		}
	}
	if len(temMap) > 0 {
		sid := serverSelf.SortMapByValue(temMap)
		return serverSelf.ServerInfo[sid]
	}
	return serverSelf.ServerInfo[0]
}

func (serverSelf *ServerInfoManager) GetAllServers(t int32) []*ServerInfo {
	serverSelf.Mutex.Lock()
	defer serverSelf.Mutex.Unlock()
	var temMap []*ServerInfo
	temMap = make([]*ServerInfo, 0, len(serverSelf.ServerInfo))
	for _, v := range serverSelf.ServerInfo {
		if v.Type == t {
			temMap = append(temMap, v)
		}
	}
	return temMap
}

func (serverSelf *ServerInfoManager) GetAll() []*ServerInfo {
	serverSelf.Mutex.Lock()
	defer serverSelf.Mutex.Unlock()
	var temMap []*ServerInfo
	temMap = make([]*ServerInfo, 0, len(serverSelf.ServerInfo))
	for _, v := range serverSelf.ServerInfo {
		temMap = append(temMap, v)
	}
	return temMap
}

//连接数增加
func (serverSelf *ServerInfoManager) AddConNum(sid server_id) {
	serverSelf.Mutex.Lock()
	defer serverSelf.Mutex.Unlock()
	for _, v := range serverSelf.ServerInfo {
		if v.Sid == sid {
			v.ConNum = v.ConNum + 1
			break
		}
	}
}

//连接数减少
func (serverSelf *ServerInfoManager) DelConNum(sid server_id) {
	serverSelf.Mutex.Lock()
	defer serverSelf.Mutex.Unlock()
	for _, v := range serverSelf.ServerInfo {
		if v.Sid == sid {
			v.ConNum = v.ConNum - 1
			break
		}
	}
}

//------------------------------
// A data structure to hold a key/value pair.
type EGOPair struct {
	Key   int32
	Value int32
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type EGOPairList []EGOPair

func (p EGOPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p EGOPairList) Len() int           { return len(p) }
func (p EGOPairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// 返回人数最少的服务器ID
func (serverSelf *ServerInfoManager) SortMapByValue(m map[int32]int32) int32 {
	p := make(EGOPairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = EGOPair{k, v}
		i++
	}
	sort.Sort(p)
	return p[0].Key
}
