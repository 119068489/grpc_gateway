package rpcserver

import (
	"context"
	"fmt"
	"grpc_gateway/proto/pb/gateway"
	"os"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
)

var c gateway.GatewayClient

func GetClient() *grpc.ClientConn {
	conn, err := grpc.Dial("localhost:9192", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logs.Error("did not connect: %v", err)
	}
	c = gateway.NewGatewayClient(conn)
	return conn
}

func TestMain(m *testing.M) {
	fmt.Println("测试之前的做一些设置,连接rpc服务器")
	conn := GetClient()
	defer conn.Close()
	// 如果 TestMain 使用了 flags，这里应该加上flag.Parse()
	retCode := m.Run() // 执行测试
	fmt.Println("测试之后做一些拆卸工作,关闭rpc服务器连接")
	os.Exit(retCode) // 退出测试
}

//go test -run=Rpcserver
func TestRpcserver(t *testing.T) {
	convey.Convey("RpcServer test", t, func() {
		convey.Convey("Echo", func() {
			res, err := c.Echo(context.Background(), &gateway.StringMessage{Value: "Echo"})
			t.Log(res)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("Gcho", func() {
			res, err := c.Gcho(context.Background(), &gateway.StringMessage{Value: "Gcho"})
			t.Log(res)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

//go test -bench Rpcserver -benchmem
func BenchmarkRpcserver(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c.Echo(context.Background(), &gateway.StringMessage{Value: "Echo"})
	}
}

func BenchmarkRpcserverParallel(b *testing.B) {
	// b.SetParallelism(4) // 设置使用的CPU数
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Echo(context.Background(), &gateway.StringMessage{Value: "Echo"})
		}
	})
}
