package rpcserver

import (
	"context"
	"encoding/base64"
	"errors"
	"io/ioutil"

	"grpc_gateway/api"
	"grpc_gateway/easygo"
	"grpc_gateway/proto/pb/gateway"
	"grpc_gateway/proto/pb/hello_world"

	"github.com/astaxie/beego/logs"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	gateway.GatewayServer
}

func (s *Server) Echo(ctx context.Context, in *gateway.StringMessage) (*gateway.StringMessage, error) {
	logs.Info("Echo: %v\n", in)
	if md, ok := metadata.FromIncomingContext(ctx); !ok {
		return nil, errors.New("无token信息")
	} else {
		l := md["authorization"]
		if len(l) > 0 {
			userName := api.CheckAuth(ctx)
			logs.Debug(userName)
		} else {
			return nil, errors.New("token 无效")
		}
	}

	msg := &gateway.StringMessage{
		Value: "Hello" + in.GetValue(),
	}
	return msg, nil
}

func (s *Server) Gcho(ctx context.Context, in *gateway.StringMessage) (*gateway.StringMessage, error) {
	logs.Info("request: ", in.Value, in.Code)

	//屏蔽链路追踪采集代码
	// newCtx, end := easygo.Start("RpcReq", ctx)
	// err := httpReq(newCtx)
	// end(easygo.SpanWithError(err), easygo.SpanWithLog("httprequet", "ok1"))
	// if err != nil {
	// 	logs.Error(err)
	// }

	return &gateway.StringMessage{Value: "Hello " + in.Value, Code: in.Code}, nil
}

func (s *Server) Upload(ctx context.Context, in *gateway.FSReq) (*gateway.FSResp, error) {
	logs.Info("request: ", in.GetName())
	msg := &gateway.FSResp{
		Status:  true,
		Message: "ok",
	}

	file, errf := base64.StdEncoding.DecodeString(in.GetFile())
	if errf != nil {
		logs.Info(errf)
	}

	// obj := &easygo.RedisAdmin{
	// 	UserId:    10001,
	// 	Role:      0,
	// 	ServerId:  101,
	// 	Timestamp: time.Now().Unix(),
	// 	Token:     "token",
	// }
	// easygo.SetRedisAdmin(obj)

	err := ioutil.WriteFile(in.GetName(), file, 0666)
	if err != nil {
		msg.Message = "fail"
		msg.Status = false
	}

	// admin := easygo.GetRedisAdmin2(10001, "UserId")
	// logs.Debug(admin)
	return msg, nil
}

func RpcReq(ctx context.Context) error {
	conn := easygo.GrpcFromGrpc("localhost:50051", ctx)
	defer conn.Close()
	c := hello_world.NewGreeterClient(conn)

	name := "ok"
	r, err := c.SauHello(ctx, &hello_world.HelloRequest{Name: name})
	if err != nil {
		logs.Error("could not greet: %v", err)
		return err
	}
	logs.Info("Greeting: %s", r.GetMessage())
	return nil
}

func HttpReq(ctx context.Context) error {
	resp := easygo.HttpFromGrpc(ctx, "http://localhost:7072/rpc", "GET", nil)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	logs.Info("http get" + string(body))
	return nil
}

func (s *Server) Login(ctx context.Context, in *gateway.LoginRequest) (*gateway.LoginReply, error) {
	logs.Info("Login: %v\n", in)
	if in.Username == "admin" && in.Password == "123456" {
		tokenString := api.CreateToken(in.Username)
		return &gateway.LoginReply{Status: "200", Token: tokenString}, nil
	} else {
		return &gateway.LoginReply{Status: "403", Token: ""}, nil
	}
}
