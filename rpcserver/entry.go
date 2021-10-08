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
	defer func() { // 若是异常了,确保异步日志有成功写盘
		logger := initializer.GetBeeLogger()
		if logger != nil {
			logger.Flush()
		}
	}()

	dict := easygo.KWAT{
		"logName":  "rpc_server",
		"yamlPath": "config_rpc.yaml",
	}
	initializer.Execute(dict) //执行公共配置初始化

	Initialize() //初始化本服特有配置

	var serveFunctions = []func(){}
	serveFunctions = append(serveFunctions, SignHandle, RpcServerRun)

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
			logs.Info("rpc服务器关闭:", easygo.PServer.GetInfo())
			//TODO:服务器关闭逻辑处理
			easygo.EtcdMgr.CancleLease()
			easygo.EtcdMgr.Close()
			time.Sleep(time.Second * 10)
			os.Exit(1)
		default:
			logs.Debug("error sign", s)
		}
	}
}

func RpcServerRun() {
	lis, err := net.Listen("tcp", easygo.SERVER_ADDR)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tracer, _ := easygo.NewJaegerTracer(easygo.SERVER_NAME, "127.0.0.1:6831")

	s := grpc.NewServer(easygo.ServerOption(tracer))
	gateway.RegisterGatewayServer(s, &Server{})
	logs.Info("Rpc server start to listen %s", easygo.SERVER_ADDR)

	s.Serve(lis)
}
