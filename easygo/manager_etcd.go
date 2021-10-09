package easygo

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type IClient3KVManager interface {
	StartClintTV3(isWatch ...bool) //连接ETCD服务器
	CancleLease()                  //撤销租约
	Close()                        //关闭Client
}

//etcd连接管理
type Client3KVManager struct {
	ServerInfoManager IServerManager //服务器管理器
	ServerType        string         //服务器类型
	ServerInfo        IServerInfo    //服务器信息
	EtcdCenterAddr    []string       //etcd中心集群地址 127.0.0.1:2379,127.0.0.1:22379,127.0.0.1:32379
	EtcdServerPath    string         //etcd 项目key path
	PClient           *clientv3.Client
	PClient3KV        clientv3.KV
	Lease             clientv3.Lease     //租约
	LeaseId           clientv3.LeaseID   //租约id
	CancleFun         context.CancelFunc //取消租约
	Mutex             RLock
}

func NewClient3KVManager(yaml IYamlConfig) *Client3KVManager { // services map[string]interface{},
	p := &Client3KVManager{}
	p.Init(yaml)
	return p
}

//初始化
func (c *Client3KVManager) Init(yaml IYamlConfig) {
	c.EtcdCenterAddr = yaml.GetValueAsArrayString("ETCD_CENTER_ADDR")
	c.EtcdServerPath = yaml.GetValueAsString("ETCD_SERVER_PATH")
	c.ServerInfoManager = ServerMgr
	c.ServerInfo = PServer.GetInfo()
}

//连接ETCD服务器
func (c *Client3KVManager) StartClintTV3(isWatch ...bool) {
	//创建连接
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	if c.ServerInfo == nil {
		return
	} else {
		c.ServerType = c.ServerInfo.GetName()
	}
	var err error
	c.PClient, err = clientv3.New(clientv3.Config{
		Endpoints:   c.EtcdCenterAddr,
		DialTimeout: 5 * time.Second,
	})
	PanicError(err)
	logs.Info("Client3KV Connecting to ETCD server", c.EtcdCenterAddr)

	c.CreateLease()
	c.SetLeaseTime(10)
	c.UpdateLeaseTime()
	//创建KV
	c.PClient3KV = clientv3.NewKV(c.PClient)
	serverInfo, err1 := json.Marshal(c.ServerInfo.GetInfo())
	PanicError(err1)
	logs.Info("Successfully connected to the ETCD server")
	//通过租约put
	_, err = c.PClient3KV.Put(context.TODO(), c.EtcdServerPath+c.ServerType+"/"+AnytoA(c.ServerInfo.GetSid()), string(serverInfo), clientv3.WithLease(c.LeaseId))
	if err != nil {
		logs.Error("put fail：", err.Error())
		PanicError(err)
	}
	logs.Info("Client3KV put ServerInfo:", c.ServerInfo)

	//初始化已有的服务器并监听服务器的变化
	isWatch = append(isWatch, true)
	if isWatch[0] {
		c.InitExistServer()
	}
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
func (c *Client3KVManager) InitExistServer() {
	kv := c.GetClientKV()
	kvs, err := kv.Get(context.TODO(), c.EtcdServerPath, clientv3.WithPrefix())
	PanicError(err)
	var serversInfo []ServerInfo
	for _, srv := range kvs.Kvs {
		s := &ServerInfo{}
		err1 := json.Unmarshal(srv.Value, s)
		PanicError(err1)

		if s.Sid == c.ServerInfo.GetSid() {
			continue
		}
		serversInfo = append(serversInfo, *s)
		c.ServerInfoManager.AddServerInfo(s)
	}
	logs.Info("Already started server:%v", serversInfo)
	//监视login服务器变化
	c.watchToServer()
}

//监听服务器的变化
func (c *Client3KVManager) watchToServer() {
	Spawn(func() {
		wc := c.GetClient().Watch(context.TODO(), c.EtcdServerPath, clientv3.WithPrefix())
		for v := range wc {
			for _, e := range v.Events {
				// logs.Info("type:%v kv:%v  prevKey:%v  value:%v\n ", e.Type, string(e.Kv.Key), e.PrevKv, e.Kv.Value)
				switch e.Type {
				case 1: //删除 mvccpb.DELETE
					//关闭无效连接
					params := strings.Split(string(e.Kv.Key), "/")
					sid := Atoi(params[3])
					c.ServerInfoManager.DelServerInfo(sid)
					logs.Info("remove ServerInfo:id=", sid)
				case 0: //增加mvccpb.PUT
					//如果已经连接
					ss := &ServerInfo{}
					if err := json.Unmarshal(e.Kv.Value, ss); err != nil {
						logs.Info("WatchToLogin err", err)
						continue
					}
					c.ServerInfoManager.AddServerInfo(ss)
					logs.Info("add ServerInfo:", ss)
				}
			}
		}
	})
}
