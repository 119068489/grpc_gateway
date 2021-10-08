package easygo

/*
公共包变量
*/

var (
	EDITION          string // 发行版
	IS_FORMAL_SERVER bool   //是否正式服：true正式服, false测试服
	IS_TFSERVER      bool   //是否走转发
	SERVER_ADDR      string //服务器ip
	SERVER_ID        int    //服务器id
	SERVER_NAME      string //服务器名称
)
