package gateway

import (
	"context"
	"flag"
	"grpc_gateway/easygo"
	pbgw "grpc_gateway/proto/pb/gateway"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func Entry(flagSet *flag.FlagSet, args []string) {
	initializer := easygo.NewInitializer()
	initializer.Execute()

	Initialize()

	//启动etcd
	PClient3KVMgr.StartClintTV3()
	defer PClient3KVMgr.Close() //关闭etcd
	//etcd注册和发现
	easygo.InitExistServer(PClient3KVMgr, PServerInfoMgr, PServerInfo)

	var serveFunctions = []func(){}
	serveFunctions = append(serveFunctions, SignHandle, GatewayServer, ProxyServer, PprofMonitor)

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
			//TODO:服务器关闭逻辑处理
			logs.Info("gateway服务器关闭:", PServerInfo)
			PClient3KVMgr.CancleLease()
			time.Sleep(time.Second * 10)
			os.Exit(1)
		default:
			logs.Debug("error sign", s)
		}
	}
}

// Interceptor 自定义拦截器
func Interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// logs.Info("method[%v];req[%v];reply[%v];cc[%+v];invoker[%T];\n", method, req, reply, cc, invoker)
	// tracer := easygo.NewSkyTracer("0.0.0.0:11800", "grpc_gateway", "localhost")
	// span := easygo.SetSkySpan(ctx, method, tracer)
	// defer span.End() //提交探针内容

	// TODO
	// 上面都是前置逻辑操作
	err := invoker(ctx, method, req, reply, cc, opts...) // 向服务端发送请求
	easygo.PanicError(err, "invoker err")
	// 下面都是后置逻辑操作
	// TODO

	return err
}

func GatewayServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	tracer, _ := easygo.NewJaegerTracer(easygo.Server_Name, "localhost:6831") //创建jaeger tracer

	//拦截器注
	opts = append(opts, easygo.DialOption(tracer), grpc.WithChainUnaryInterceptor(Interceptor)) //自定义拦截器

	ps := PServerInfoMgr.GetIdelServer(easygo.SERVER_TYPE_RPC) //获取rpc服务器配置
	adds := "localhost:9192"
	if ps != nil {
		adds = ps.InternalIP
	} else {
		logs.Error("No found rpc server,Default listening %s", adds)
	}

	echoEndpoint := flag.String("echo_endpoint", adds, "endpoint of Gateway")
	err := pbgw.RegisterGatewayHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
	easygo.PanicError(err)

	easygo.ServerRun(easygo.Server_IP, mux, "Gateway server")
}

func ProxyServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", getEntry).Methods("GET")
	r.HandleFunc("/", proxyUpEntry).Methods("POST")      //上传请求
	r.PathPrefix("/v1/example/").HandlerFunc(proxyEntry) //前缀匹配rpc请求

	addr := ":9190"
	easygo.ServerRun(addr, r, "Proxy server") //对外的端口
}

func PprofMonitor() {
	addr := ":" + strconv.Itoa(6060)
	easygo.ServerRun(addr, nil, "Pprof Server")
}
