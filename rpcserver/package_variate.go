package rpcserver

import eg "grpc_gateway/easygo"

var PServerInfo *eg.ServerInfo
var PServerInfoMgr *eg.ServerInfoManager
var PClient3KVMgr *eg.Client3KVManager

func Initialize() {
	eg.Server_IP = "127.0.0.1:9192"
	eg.Server_ID = 101
	eg.Server_Name = "rpc_server"

	PServerInfo = &eg.ServerInfo{
		Sid:        eg.Server_ID,
		Name:       eg.Server_Name,
		Type:       eg.SERVER_TYPE_RPC,
		ExternalIp: eg.Server_IP,
		InternalIP: eg.Server_IP,
		State:      eg.SERVER_StATE_ON,
		ConNum:     0,
	}

	PClient3KVMgr = eg.NewClient3KVManager(eg.Server_Name, PServerInfo)
	PServerInfoMgr = eg.NewServerInfoManager()
}
