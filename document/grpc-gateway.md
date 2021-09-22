# grpc-gateway
* 简介
grpc-gateway是一个protoc插件。它读取gRPC服务定义并生成反向代理服务，将 gRPC 服务代理为 RESTful JSON API。该服务是根据 .proto 文件定义中的自定义选项生成的。

## 安装
```
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0 \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.5.0 \
    google.golang.org/protobuf/cmd/protoc-gen-go@v1.26.0 \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
```
安装后将二进制文件加入PATH

## Usage

1. 定义服务
```helloworld.proto
// Copyright 2015 gRPC authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

syntax = "proto3";

option go_package = "study/grpc/gateway/helloworld/pb";
option java_multiple_files = true;
option java_package = "study.grpc.gateway.helloworld.pb";
option java_outer_classname = "HelloWorldProto";

package pb;

// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
  // Sends a greeting
  rpc SayHelloClientStream (stream HelloRequest) returns (HelloReply) {}
  // Sends a greeting
  rpc SayHelloServerStream (HelloRequest) returns (stream HelloReply) {}
  // Sends a greeting
  rpc SayHelloBidStream (stream HelloRequest) returns (stream HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}

```
2. 根据定义生成grpc服务
```
protoc -I . \
    --go_out . --go_opt paths=source_relative \
    --go-grpc_out . --go-grpc_opt paths=source_relative \
    study/grpc/gateway/helloworld/pb/helloworld.proto
```

3. 实现服务
```
package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes" // grpc 响应状态码
	"google.golang.org/grpc/status"

	//"google.golang.org/grpc/credentials" // grpc认证包
	// "google.golang.org/grpc/grpclog"
	pb "xk/test/study/grpc/gateway/helloworld/pb"

	"google.golang.org/grpc/metadata" // grpc metadata包
)

const (
	timestampFormat = time.StampNano
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) SayHelloServerStream(in *pb.HelloRequest, stream pb.Greeter_SayHelloServerStreamServer) error {
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
		stream.SetTrailer(trailer)
	}()
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.DataLoss, "SayHelloServerStream: failed to get metadata")
	}
	if t, ok := md["timestamp"]; ok {
		log.Printf("timestamp from metadata:\n")
		for i, e := range t {
			log.Printf("%d. %s\n", i, e)
		}
	}
	header := metadata.New(map[string]string{"location": "MTV", "timestamp": time.Now().Format(timestampFormat)})
	stream.SendHeader(header)
	log.Printf("request receive: %v\n", in)

	for i := 0; i < 10; i++ {
		log.Printf("echo message %v\n", in.Name)
		err := stream.Send(&pb.HelloReply{Message: "Hello " + in.Name})
		if err != nil {
			return err
		}
	}
	return nil
}

func auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "无Token认证信息")
	}
	var (
		appid  string
		appkey string
	)
	if val, ok := md["appid"]; ok {
		appid = val[0]
	}

	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}

	if appid != "101010" || appkey != "i am key" {
		return grpc.Errorf(codes.Unauthenticated, "Token认证信息无效: appid=%s, appkey=%s", appid, appkey)
	}

	return nil

}

// 流拦截器， 除了一元消息， 剩下的三种流消息都走这个过滤器
func streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Printf("start stream interceptor\n")
	err := auth(ss.Context())
	if err != nil {
		return err
	}

	// 继续处理请求
	return handler(srv, ss)
}

// interceptor 一元拦截器
func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}

	// 继续处理请求
	return handler(ctx, req)
}

func Start(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	// TLS认证
	// creds, err := credentials.NewServerTLSFromFile("./keys/server.crt", "./keys/server_rsa_private.pem")
	// if err != nil {
	//     grpclog.Fatalf("Failed to generate credentials %v", err)
	// }

	// opts = append(opts, grpc.Creds(creds))
	// 注册一元 interceptor
	opts = append(opts, grpc.ChainUnaryInterceptor([]grpc.UnaryServerInterceptor{interceptor}...))

	// 注册流 interceptor
	opts = append(opts, grpc.StreamInterceptor(streamInterceptor))

	s := grpc.NewServer(opts...)
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

```

4. 生成反向代理服务
```
protoc -I . --grpc-gateway_out ./gen/go \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
    study/grpc/gateway/helloworld/pb/helloworld.proto
```

当前目录下 proto文件定义的 go_package位置会生成文件
helloworld.pb.gw.go

5. 启动反向代理服务
```main.go
package main

import (
	"context"
	"flag"
	"net/http"

	gw "xk/test/study/grpc/gateway/helloworld/pb"
	"xk/test/study/grpc/interceptor"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:80", "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
    // 注入了客户端拦截器
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithPerRPCCredentials(new(interceptor.CustomCredential))}
	err := gw.RegisterGreeterHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8081", mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}

```
启动反向代理服务

`go run main.go`


6. curl 访问 反向代理服务
url格式为 "{hostname}:{port}/{package}.{service}/{api}"
- package proto文件中 package的名字
- service proto文件中 service的名字
- api proto文件中的接口名

例如：

`curl -XPOST "http://localhost:8081/pb.Greeter/SayHello" -d '{"name": "world"}'`

7. 自定义访问方式

除了使用默认的.proto文件生成gateway服务之外，还可以修改.proto文件，生成自定义的url访问方式


## 自定义访问

### 修改proto文件，注入url访问方式

1. 下载依赖.proto文件

`git clone https://github.com/googleapis/googleapis.git`
`git clone https://github.com/protocolbuffers/protobuf.git`

2. 修改自定义的proto文件

```helloworld.proto

+import "google/api/annotations.proto"; // 引入依赖

-  rpc SayHello (HelloRequest) returns (HelloReply) {}
+  rpc SayHello (HelloRequest) returns (HelloReply) {
+	// 自定义url, name 参数为 HelloRequest中的字段
+         option (google.api.http) = {
+            get: "/v1/sayhello/{name}"
+          };
+  }

```

3. 生成gateway 服务代码

```
protoc -I . -I "依赖proto文件所在位置，多个可以使用多个-I引入" --grpc-gateway_out ./gen/go \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
    study/grpc/gateway/helloworld/pb/helloworld.proto
```

### 使用服务配置文件

1. 配置yaml服务文件
```helloworld-service.yaml
type: google.api.Service
config_version: 3

http:
  rules:
    - selector: pb.Greeter.SayHello
      get: /v1/sayhello/{name}
```

## 来源

[ github.com ](https://github.com/grpc-ecosystem/grpc-gateway)

[ grpc-gateway doc ](https://grpc-ecosystem.github.io/grpc-gateway/#getting-started)