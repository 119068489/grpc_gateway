# jaeger
OpenTracing 是开放式分布式追踪规范，OpenTracing API 是一致，可表达，与供应商无关的API，用于分布式跟踪和上下文传播。

OpenTracing 的客户端库以及规范，可以到[Github](https://github.com/opentracing/)中查看

Jaeger 是 Uber 开源的分布式跟踪系统

## 环境搭建
环境：
elasticsearch 7
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
chown -R elsearch:elsearch elasticsearch-7
切换用户
su elsearch
启动elasticsearch
切换到elasticsearch/bin目录，执行命令：
./elasticsearch -d

## 安装jaeger
官网下载安装包解压即可

## 启动jaeger
如果是本地调试可以启动一个all-in-one
`jaeger-all-in-one --collector.zipkin.host-port=:9411 --span-storage.type=elasticsearch`

1. collector 用来提交数据到ES数据库
   `jaeger-collector --span-storage.type=elasticsearch`

2. query 用来查询数据UI服务 [UI](http://127.0.0.1:16686/)
   `jaeger-query  --span-storage.type=elasticsearch`
   
3. agent 用于客户端提交数据的接口
   `jaeger-agent --reporter.grpc.host-port=localhost:14250`

## jaeger client go

manage_jaeger.go 封装的工具包
```
package easygo

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	_ "github.com/uber/jaeger-client-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

var Oldspan opentracing.Span

// 配置jaeger addr="127.0.0.1:6831" UDP端口6831
func NewJaegerTracer(serviceName string, jagentHost string) (opentracing.Tracer, io.Closer) {
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1, //字段 Param 是设置采样率或速率，要根据 Type 而定。
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,            //是否写日志在推送中
			BufferFlushInterval: 1 * time.Second, //在提交队列不满的情况下推送的频率
			LocalAgentHostPort:  jagentHost,      //要推送到的 Jaeger agent，默认端口 6831，是 Jaeger 接收压缩格式的 thrift 协议的数据端口。
			// CollectorEndpoint:   "http://127.0.0.1:14268/api/traces",//要推送到的 Jaeger Collector，用 Collector 就不用 agent 了。
		},
	}

	reporter, _ := cfg.Reporter.NewReporter(serviceName, jaeger.NewNullMetrics(), jaeger.NullLogger)
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Reporter(reporter),
	)
	PanicError(err)
	opentracing.InitGlobalTracer(tracer)
	return tracer, closer
}

// SkyTracer agent探针上报对象获取 tracerIp, serviceName, localName = "0.0.0.0:11800", "grpc_gateway", "localhost"
func NewSkyTracer(tracerIp, serviceName, localName string) *go2sky.Tracer {
	r, err := reporter.NewGRPCReporter(tracerIp)
	PanicError(err, "New GRPC Reporter Error")

	tracer, err := go2sky.NewTracer(serviceName, go2sky.WithReporter(r), go2sky.WithInstance(localName))
	PanicError(err, "New Tracer Error")

	return tracer
}

// SkySpan 放置agent探针在请求发出之前
func SetSkySpan(ctx context.Context, method string, tracer *go2sky.Tracer) go2sky.Span {
	span, _, err := tracer.CreateEntrySpan(ctx, method, func(header string) (string, error) {
		return "", nil
	})
	PanicError(err)

	span.SetComponent(5200)
	span.Tag(go2sky.TagHTTPMethod, method)
	span.SetSpanLayer(1)

	return span
}

// metadataTextMap extends a metadata.MD to be an opentracing textmap
type MetadataTextMap struct {
	metadata.MD
}

// Set is a opentracing.TextMapReader interface that extracts values.
func (m MetadataTextMap) Set(key, val string) {
	key = strings.ToLower(key)
	m.MD[key] = append(m.MD[key], val)
}

// ForeachKey is a opentracing.TextMapReader interface that extracts values.
func (c MetadataTextMap) ForeachKey(handler func(key, val string) error) error {
	for k, vs := range c.MD {
		for _, v := range vs {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// DialOption grpc client option
//grpc拦截器注入jaeger
func DialOption(tracer opentracing.Tracer) grpc.DialOption {
	return grpc.WithUnaryInterceptor(ClientInterceptor(tracer))
}

// ServerOption grpc server option
func ServerOption(tracer opentracing.Tracer) grpc.ServerOption {
	return grpc.UnaryInterceptor(ServerInterceptor(tracer))
}

//http服务端Handler拦截注入jaeger
func ServerHttpHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return ServerHookHttpHandler(opentracing.GlobalTracer(), next)
	}
}

// ServerInterceptor http server wrapper
// http服务端Handler拦截注入jaeger
func ServerHookHttpHandler(tracer opentracing.Tracer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var opts []opentracing.StartSpanOption
		opts = append(opts,
			opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
			ext.SpanKindProducer)

		var ParentSpan opentracing.Span
		spCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err == nil {
			opts = append(opts, opentracing.ChildOf(spCtx))
		}

		ParentSpan = opentracing.StartSpan(
			r.URL.Path,
			opts...,
		)
		defer ParentSpan.Finish()

		tracer.Inject(ParentSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header)) //拦截记录span后,刷新r

		Oldspan = ParentSpan

		//请求前
		next.ServeHTTP(w, r) //去请求
		//请求后
	})
}

//ClientInterceptor grpc Client wrapper
//grpc拦截器客户端注入jaeger
func ClientInterceptor(tracer opentracing.Tracer) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string,
		req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		// var parentCtx opentracing.SpanContext
		// parentSpan := opentracing.SpanFromContext(ctx)
		// if parentSpan != nil {
		// 	parentCtx = parentSpan.Context()
		// }

		// span := tracer.StartSpan(
		// 	method,
		// 	opentracing.ChildOf(parentCtx),
		// 	opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
		// 	ext.SpanKindRPCClient,
		// )
		// defer span.Finish()
		if Oldspan != nil {
			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				md = metadata.New(nil)
			} else {
				md = md.Copy()
			}

			mdWriter := MetadataTextMap{md}
			err := tracer.Inject(Oldspan.Context(), opentracing.TextMap, mdWriter)
			if err != nil {
				Oldspan.LogFields(log.String("inject-error", err.Error()))
			}

			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			Oldspan.LogFields(log.String("call-error", err.Error()))
		}
		return err
	}
}

// ServerInterceptor grpc server wrapper
// grpc拦截器服务端注入jaeger
func ServerInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		spanContext, err := tracer.Extract(opentracing.TextMap, MetadataTextMap{md})
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			grpclog.Errorf("extract from metadata err: %v", err)
		} else {
			span := tracer.StartSpan(
				info.FullMethod,
				ext.RPCServerOption(spanContext),
				opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
				ext.SpanKindRPCServer,
			)
			defer span.Finish()
			Oldspan = span
			ctx = opentracing.ContextWithSpan(ctx, span)
		}

		return handler(ctx, req)
	}
}

type SpanOption func(span opentracing.Span)

func SpanWithError(err error) SpanOption {
	return func(span opentracing.Span) {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogFields(log.String("event", "error"), log.String("msg", err.Error()))
		}
	}
}

// example:
// SpanWithLog(
//    "event", "soft error",
//    "type", "cache timeout",
//    "waited.millis", 1500)
func SpanWithLog(arg ...interface{}) SpanOption {
	return func(span opentracing.Span) {
		span.LogKV(arg...)
	}
}

//jaeger子程序或代码块监控采集
//
// newCtx, end := easygo.Start("RpcReq", ctx)
// err := httpReq(newCtx)
// end(easygo.SpanWithError(err), easygo.SpanWithLog("httprequet", "ok1"))
func Start(spanName string, ctx context.Context) (newCtx context.Context, finish func(...SpanOption)) {
	if ctx == nil {
		ctx = context.TODO()
	}
	span, newCtx := opentracing.StartSpanFromContextWithTracer(ctx, opentracing.GlobalTracer(), spanName,
		opentracing.Tag{Key: string(ext.Component), Value: "func"},
	)

	finish = func(ops ...SpanOption) {
		for _, o := range ops {
			o(span)
		}
		span.Finish()
	}

	return
}

//http -> http请求通过request注入request
func RequestWithRequest(r *http.Request, req *http.Request) {
	tracer := opentracing.GlobalTracer()
	spCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		logs.Error(err, "Couldn't Extract headers")
	}

	// http请求创建span
	// span := tracer.StartSpan(
	// 	r.URL.Path,
	// 	opentracing.ChildOf(spCtx),
	// 	opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
	// 	ext.SpanKindConsumer,
	// )
	// defer span.Finish()

	err = tracer.Inject(spCtx, opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		logs.Error(err, "Couldn't inject headers")
	}
}

//http -> grpc请求通过request注入Context
func ContextWithRequst(r *http.Request) context.Context {
	// var opts []opentracing.StartSpanOption
	// opts = append(opts,
	// 	opentracing.Tag{Key: string(ext.Component), Value: "grpc"},
	// 	ext.SpanKindRPCClient)

	// var ParentSpan opentracing.Span
	// spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	// if err == nil {
	// 	opts = append(opts, ext.RPCServerOption(spCtx))
	// }

	// ParentSpan = opentracing.StartSpan(
	// 	r.URL.Path,
	// 	opts...,
	// )
	// defer ParentSpan.Finish()

	ctx := r.Context()
	ctx = opentracing.ContextWithSpan(ctx, Oldspan)

	return ctx
}

//grpc -> http请求通过Context注入request
func RequestWithContext(ctx context.Context, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}

	// http请求创建span
	// span = tracer.StartSpan(
	// 	r.URL.Path,
	// 	opentracing.ChildOf(span.Context()),
	// 	opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
	// 	ext.SpanKindConsumer,
	// )
	// defer span.Finish()

	err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		logs.Error(err, "Couldn't inject headers")
	}
}

//grpc -> grpc请求通过ctx注入ctx
func ContextWithContext(ctx context.Context, method ...string) context.Context {

	parentSpan := opentracing.SpanFromContext(ctx)

	// var parentCtx opentracing.SpanContext
	// if parentSpan != nil {
	// 	parentCtx = parentSpan.Context()
	// }

	// span := opentracing.GlobalTracer().StartSpan(
	// 	method,
	// 	opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
	// 	ext.RPCServerOption(parentCtx),
	// 	ext.SpanKindRPCClient,
	// )
	// defer span.Finish()

	return opentracing.ContextWithSpan(ctx, parentSpan)
}

//grpc -> grpc请求
func GrpcFromGrpc(host string, ctx context.Context, method ...string) *grpc.ClientConn {
	ctx = ContextWithContext(ctx)
	conn, err := grpc.DialContext(ctx, host, grpc.WithInsecure(), grpc.WithBlock(), DialOption(opentracing.GlobalTracer()))
	PanicError(err, "did not connect")
	return conn
}

//http -> grpc请求
func GrpcFromHttp(host string, r *http.Request) (context.Context, *grpc.ClientConn) {
	ctx := ContextWithRequst(r)
	conn, err := grpc.DialContext(ctx, host, grpc.WithInsecure(), grpc.WithBlock(), DialOption(opentracing.GlobalTracer()))
	if err != nil {
		logs.Error("did not connect: %v", err)
	}
	return ctx, conn
}

// grpc -> http请求
func HttpFromGrpc(ctx context.Context, url, method string, body io.Reader) *http.Response {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	PanicError(err)
	RequestWithContext(ctx, req)
	resp, err := http.DefaultClient.Do(req)
	PanicError(err)
	return resp
}

//http -> http请求
func HttpFromHttp(r *http.Request, url, method string, body io.Reader) *http.Response {
	req, _ := http.NewRequest(method, url, body)
	RequestWithRequest(r, req)
	resp, err := http.DefaultClient.Do(req)
	PanicError(err)
	return resp
}
```

以下方法是服务端接收用法，已经封装了创建监控节点span和提交span的代码，因此会产生追踪记录。

1. 创建追踪器 tracer 对象
   ```
   Tracer, _ = easygo.NewJaegerTracer(easygo.Server_Name, "localhost:6831")
   ```

2. rpc服务器中使用方法
   ```
    //创建追踪器对象
    tracer, _ := easygo.NewJaegerTracer(easygo.Server_Name, "127.0.0.1:6831")
    //拦截器注册
    s := grpc.NewServer(easygo.ServerOption(tracer))
	gateway.RegisterGatewayServer(s, &Server{})
   ```

3. http服务器中使用方法
   ```
    //创建追踪器
    easygo.NewJaegerTracer("serverName", "127.0.0.1:6831")
	//拦截器
	next := easygo.ServerHttpHandler()
	err := http.ListenAndServe(listenAddr, next(route))
   ```

4. func代码块监控使用方法
   ```
   //创建span节点
   newCtx, end := easygo.Start("funcName", ctx)
   //要监控的func或代码块
   //TODO 代码块
   returnErr := httpReq(newCtx)
   //写错误日志和操作日志
   end(easygo.SpanWithError(returnErr), easygo.SpanWithLog("logKey", "logValue")) 
   ```

以下方法是客户端发起请求注入追踪参数用法，没有封装创建监控节点span和提交span的代码，不会产生追踪记录。

如果要在请求的发起客户端也产生追踪记录，需要开启注入函数的注释创建span代码，例如：
```
// http请求创建span
// span = tracer.StartSpan(
// 	r.URL.Path,
// 	opentracing.ChildOf(span.Context()),
// 	opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
// 	ext.SpanKindConsumer,
// )
// defer span.Finish()
```

1. grpc服务端 -> grpc请求
   ``` 
   func GrpcFromGrpc(host string, ctx context.Context, method ...string) *grpc.ClientConn {
   	ctx = ContextWithContext(ctx) //通过ctx注入ctx
   	conn, err := grpc.DialContext(ctx, host, grpc.WithInsecure(), grpc.WithBlock(), DialOption(opentracing.GlobalTracer()))
   	PanicError(err, "did not connect")
   	return conn
   }
   ```

2. http服务端 -> grpc请求
   ```
   func GrpcFromHttp(host string, r *http.Request) (context.Context, *grpc.ClientConn) {
      	ctx := ContextWithRequst(r) //通过request注入Context
      	conn, err := grpc.DialContext(ctx, host, grpc.WithInsecure(), grpc.WithBlock(), DialOption(opentracing.GlobalTracer()))
      	if err != nil {
      		logs.Error("did not connect: %v", err)
      	}
      	return ctx, conn
   }
   ```

3. grpc服务端 -> http请求
   ```
   func HttpFromGrpc(ctx context.Context, url, method string, body io.Reader) *http.Response {
      	req, err := http.NewRequestWithContext(ctx, method, url, nil)
      	PanicError(err)
      	RequestWithContext(ctx, req) //通过Context注入request
      	resp, err := http.DefaultClient.Do(req)
      	PanicError(err)
      	return resp
   }
   ```

4. http服务端 -> http请求
   ```
    func HttpFromHttp(r *http.Request, url, method string, body io.Reader) *http.Response {
      	req, _ := http.NewRequest(method, url, body)
      	RequestWithRequest(r, req) //通过request注入request
      	resp, err := http.DefaultClient.Do(req)
      	PanicError(err)
      	return resp
   }
   ```
