package easygo

/*
公共包变量
*/

var (
	EDITION          string   // 发行版
	IS_FORMAL_SERVER bool     //是否正式服：true正式服, false测试服
	ETCD_CENTER_ADDR []string //etcd中心集群地址 127.0.0.1:2379,127.0.0.1:22379,127.0.0.1:32379
	ETCD_SERVER_PATH string   //etcd 项目key path
	IS_TFSERVER      bool     //是否走转发
	Server_IP        string   //服务器ip
	Server_ID        int32    //服务器id
	Server_Name      string   //服务器名称
)
