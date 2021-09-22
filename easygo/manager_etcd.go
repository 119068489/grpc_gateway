package easygo

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	clientv3 "go.etcd.io/etcd/client/v3"
)

//etcd连接管理
type Client3KVManager struct {
	ServerType string      //服务器类型
	ServerInfo *ServerInfo //服务器信息
	PClient    *clientv3.Client
	PClient3KV clientv3.KV
	Lease      clientv3.Lease     //租约
	LeaseId    clientv3.LeaseID   //租约id
	CancleFun  context.CancelFunc //取消租约
	Mutex      RLock
}

func NewClient3KVManager(serverType string, serverInfo *ServerInfo) *Client3KVManager { // services map[string]interface{},
	p := &Client3KVManager{}
	p.Init(serverType, serverInfo)
	return p
}

//初始化
func (c *Client3KVManager) Init(serverType string, serverInfo *ServerInfo) {
	c.ServerType = serverType
	c.ServerInfo = serverInfo
}

//连接ETCD服务器
func (c *Client3KVManager) StartClintTV3() {
	//创建连接
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	var err error
	c.PClient, err = clientv3.New(clientv3.Config{
		Endpoints:   ETCD_CENTER_ADDR,
		DialTimeout: 5 * time.Second,
	})
	PanicError(err)
	logs.Info("Client3KV 连接etcd服务器...")
	c.CreateLease()
	c.SetLeaseTime(10)
	c.UpdateLeaseTime()
	//创建KV
	c.PClient3KV = clientv3.NewKV(c.PClient)
	serverInfo, err1 := json.Marshal(c.ServerInfo)
	PanicError(err1)
	//通过租约put
	_, err = c.PClient3KV.Put(context.TODO(), ETCD_SERVER_PATH+c.ServerType+"/"+AnytoA(c.ServerInfo.Sid), string(serverInfo), clientv3.WithLease(c.LeaseId))
	if err != nil {
		logs.Info("put 失败：%s", err.Error())
		PanicError(err)
	}
	logs.Info("Client3KV put服务器信息:", c.ServerInfo)
}

//创建租约
func (c *Client3KVManager) CreateLease() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Lease = clientv3.NewLease(c.PClient)
}

//设置租约时间
func (c *Client3KVManager) SetLeaseTime(t int64) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	leaseResp, err := c.Lease.Grant(context.TODO(), t)
	if err != nil {
		logs.Info("设置租约失败")
		panic(err)
	}
	c.LeaseId = leaseResp.ID
}

//设置续租
func (c *Client3KVManager) UpdateLeaseTime() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	var ctx context.Context
	ctx, c.CancleFun = context.WithCancel(context.TODO())
	leaseRespChan, err := c.Lease.KeepAlive(ctx, c.LeaseId)
	if err != nil {
		logs.Info("续租失败:", err.Error())
		panic(err)
	}

	//监听租约
	Spawn(func() {
		for {
			select {
			case <-context.TODO().Done():
				logs.Info("已经关闭")
				return
			case leaseKeepResp, ok := <-leaseRespChan:
				if !ok {
					logs.Info("已经关闭续租功能", leaseKeepResp)
					return
				} else {
					goto END
				}
			}
		END:
			time.Sleep(500 * time.Millisecond)
		}
	})
}

//撤销租约
func (c *Client3KVManager) CancleLease() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.CancleFun()
	_, err := c.Lease.Revoke(context.TODO(), c.LeaseId)
	if err != nil {
		logs.Info("撤销租约失败:%s", err.Error())
		panic(err)
	}
	logs.Info("撤销租约成功")
}

//监听某个key值变化
func (c *Client3KVManager) WatchClientKV(key string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	Spawn(func() {
		wc := c.PClient.Watch(context.TODO(), key, clientv3.WithPrevKV())
		for v := range wc {
			for _, e := range v.Events {
				logs.Info("type:%v kv:%v  prevKey:%v \n ", e.Type, string(e.Kv.Key), e.PrevKv)
			}
		}
	})
}

//关闭Client
func (c *Client3KVManager) Close() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.PClient.Close()
}
func (c *Client3KVManager) GetClient() *clientv3.Client {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.PClient
}
func (c *Client3KVManager) GetClientKV() clientv3.KV {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.PClient3KV
}

//初始化etcd已存在的服务器
func InitExistServer(pClient3KVMgr *Client3KVManager, pServerInfoMgr *ServerInfoManager, server *ServerInfo) {
	client := pClient3KVMgr.GetClient()
	kv := pClient3KVMgr.GetClientKV()
	kvs, err := kv.Get(context.TODO(), ETCD_SERVER_PATH, clientv3.WithPrefix())
	PanicError(err)
	logs.Info("已经启动的服务器:", kvs.Kvs)
	for _, srv := range kvs.Kvs {
		s := &ServerInfo{}
		err1 := json.Unmarshal(srv.Value, s)
		PanicError(err1)
		if s.Sid == server.Sid {
			continue
		}
		// logs.Info("服务器:", *s.Sid, s)
		pServerInfoMgr.AddServerInfo(s)
	}
	//监视login服务器变化
	WatchToServer(client, ETCD_SERVER_PATH, pServerInfoMgr)
}

//监听服务器的变化
func WatchToServer(clt *clientv3.Client, key string, pServerInfoMgr *ServerInfoManager) {
	Spawn(func() {
		wc := clt.Watch(context.TODO(), key, clientv3.WithPrefix())
		for v := range wc {
			for _, e := range v.Events {
				// logs.Info("type:%v kv:%v  prevKey:%v  value:%v\n ", e.Type, string(e.Kv.Key), e.PrevKv, e.Kv.Value)
				switch e.Type {
				case 1: //删除 mvccpb.DELETE
					//关闭无效连接
					params := strings.Split(string(e.Kv.Key), "/")
					sid := AtoInt32(params[3])
					pServerInfoMgr.DelServerInfo(sid)
					logs.Info("remove ServerInfo:id=", sid)
				case 0: //增加mvccpb.PUT
					//如果已经连接
					ss := &ServerInfo{}
					if err := json.Unmarshal(e.Kv.Value, ss); err != nil {
						logs.Info("WatchToLogin err", err)
						continue
					}
					pServerInfoMgr.AddServerInfo(ss)
					logs.Info("add ServerInfo:", ss)
				}
			}
		}
	})
}
