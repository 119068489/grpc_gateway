package rpcserver

import (
	"flag"
	"grpc_gateway/easygo"
	"grpc_gateway/proto/pb/gateway"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/astaxie/beego/logs"
	"google.golang.org/grpc"
)

func Entry(flagSet *flag.FlagSet, args []string) {
	initializer := easygo.NewInitializer()
	initializer.Execute()

	Initialize()

	//启动etcd
	PClient3KVMgr.StartClintTV3()
	defer PClient3KVMgr.Close()

	//etcd已存在的服务器
	easygo.InitExistServer(PClient3KVMgr, PServerInfoMgr, PServerInfo)

	var serveFunctions = []func(){}
	serveFunctions = append(serveFunctions, SignHandle, RpcRun)

	jobs := []easygo.IGoroutine{}
	for _, f := range serveFunctions {
		job := easygo.Spawn(f)
		jobs = append(jobs, job)
	}
	easygo.JoinAllJobs(jobs...)
}

func SignHandle() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	for {
		s := <-c
		switch s {
		case syscall.SIGTERM:
			logs.Info("rpc服务器关闭:", PServerInfo)
			//TODO:服务器关闭逻辑处理
			PClient3KVMgr.CancleLease()
			time.Sleep(time.Second * 10)
			os.Exit(1)
		default:
			logs.Debug("error sign", s)
		}
	}
}

func RpcRun() {

	lis, err := net.Listen("tcp", easygo.Server_IP)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tracer, _ := easygo.NewJaegerTracer(easygo.Server_Name, "127.0.0.1:6831")

	s := grpc.NewServer(easygo.ServerOption(tracer))
	gateway.RegisterGatewayServer(s, &Server{})
	log.Println("rpc服务已经开启")

	s.Serve(lis)
}
