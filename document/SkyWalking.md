# SkyWalking
SkyWalking 是一个基于 OpenTracing 规范的、开源的 APM 系统，它是专门为微服务架构以及云原生架构而设计的，支持多种语言的客户端，部署简单，快速，目前在业界使用较为广泛。具体的skywalking安装部署参照上一篇博文：SkyWalking搭建。
由于我司部分底层服务用golang实现，为了做apm分析，需要集成skywalking go agent。

## 环境搭建
环境：
elasticsearch 6.3.2
skywalking 6.3.0
jdk 1.8

## 安装elasticsearch
官网下载安装包解压即可

## 配置elasticsearch
切换到elasticsearch配置文件目录，目录为elasticsearch-7/config目录下elasticsearch.yml文件，需要更改的配置见下：
cluster.name: CollectorDBCluster #此名称需要和collector配置文件一致。
node.name: CollectorDBCluster1，
network.host: 127.0.0.1 #本机ip地址
创建用户
elasticsearch无法以root用户身份启动，需要创建用户，创建命令：
useradd elsearch
chown -R elsearch:elsearch elasticsearch-6.3.2
切换用户
su elsearch
启动elasticsearch
切换到elasticsearch/bin目录，执行命令：
./elasticsearch -d

## 安装SkyWalking
官网下载安装包解压即可

## 配置
配置config/application.yml文件
修改storage为elasticsearch-7

## 启动
bin目录下执行startup.bat即可

## agent
1. 在grpc-gateway入口代码中加入代码
    ```注释部分即可
        func HttpRun() {
            flag.StringVar(&listenAddr, "listen-addr", ":9090", "listen address")
            ctx := context.Background()
            ctx, cancel := context.WithCancel(ctx)
            defer cancel()

            mux := runtime.NewServeMux()
            opts := []grpc.DialOption{grpc.WithInsecure()}

            // r, err := reporter.NewGRPCReporter("0.0.0.0:11800")
            // easygo.PanicError(err, "New GRPC Reporter Error")

            // tracer, err := go2sky.NewTracer("grpc_gateway", go2sky.WithReporter(r), go2sky.WithInstance("getway"))
            // easygo.PanicError(err, "New Tracer Error")

            //拦截器注册
            Interceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
                // log.Printf("method[%v];req[%v];reply[%v];cc[%+v];invoker[%T];\n", method, req, reply, cc, invoker)
                // span, ctx, err := tracer.CreateEntrySpan(ctx, method, func(header string) (string, error) {
                // 	return "", nil
                // })
                // easygo.PanicError(err)

                // span.SetComponent(5200)
                // span.Tag(go2sky.TagHTTPMethod, method)
                // span.SetSpanLayer(1)

                // 上面都是前置逻辑操作
                err := invoker(ctx, method, req, reply, cc, opts...) // 向服务端发送请求
                easygo.PanicError(err, "invoker err")
                // 下面都是后置逻辑操作
                // span.End()
                logs.Debug("请求已经发送了")
                return err
            }
            opts = append(opts, grpc.WithChainUnaryInterceptor(Interceptor))

            ps := PServerInfoMgr.GetIdelServer(easygo.SERVER_TYPE_RPC) //获取rpc服务器配置
            adds := "localhost:9192"
            if ps != nil {
                adds = ps.InternalIP
            } else {
                logs.Error("没有发现rpc服务器")
            }

            echoEndpoint := flag.String("echo_endpoint", adds, "endpoint of Gateway")
            err := pbgw.RegisterGatewayHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
            easygo.PanicError(err)

            logs.Info("http服务开启")
            http.ListenAndServe(listenAddr, mux)
        }
    ```

2. http服务器中加入代码
   ```
   func ServerHTTP(writer http.ResponseWriter, request *http.Request) {
   	time.Sleep(time.Duration(500) * time.Millisecond)
   	go2sky.PutCorrelation(request.Context(), "MIDDLE_KEY", "go2sky")

   	str := `{"value":"world"}`
   	clientReq, err := http.NewRequest(http.MethodPost, upstreamURL, strings.NewReader(str))
   	if err != nil {
   		writer.WriteHeader(http.StatusInternalServerError)
   		log.Printf("unable to create http request error: %v \n", err)
   		return
   	}
   	clientReq = clientReq.WithContext(request.Context())
   	res, err := client.Do(clientReq)
   	if err != nil {
   		writer.WriteHeader(http.StatusInternalServerError)
   		log.Printf("unable to do http request error: %v \n", err)
   		return
   	}
   	defer res.Body.Close()
   	body, err := ioutil.ReadAll(res.Body)
   	if err != nil {
   		writer.WriteHeader(http.StatusInternalServerError)
   		log.Printf("read http response error: %v \n", err)
   		return
   	}
   	writer.WriteHeader(res.StatusCode)
   	_, _ = writer.Write(body)
   }

   func main() {
   	flag.Parse()

   	var report go2sky.Reporter
   	var err error
   	report, err = reporter.NewGRPCReporter("0.0.0.0:11800")
   	easygo.PanicError(err, "New GRPC Reporter Error")

   	tracer, err := go2sky.NewTracer(serviceName, go2sky.WithReporter(report), go2sky.WithInstance("getway"))
   	easygo.PanicError(err, "New Tracer Error")

   	client, err = httpPlugin.NewClient(tracer)
   	if err != nil {
   		log.Fatalf("create client error %v \n", err)
   	}
   	sm, err := httpPlugin.NewServerMiddleware(tracer)
   	if err != nil {
   		log.Fatalf("create server middleware error %v \n", err)
   	}

   	route := http.NewServeMux()
   	route.HandleFunc("/", ServerHTTP)

   	err = http.ListenAndServe(listenAddr, sm(route))
   	if err != nil {
   		log.Fatal(err)
   	}
   }
   ```