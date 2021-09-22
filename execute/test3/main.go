package main

import (
	"flag"
	"grpc_gateway/easygo"
	"grpc_gateway/proto/pb/gateway"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego/logs"
)

var (
	listenAddr string
)

func init() {
	flag.StringVar(&listenAddr, "listen-addr", ":7073", "listen address")
}

func ServerRpc(w http.ResponseWriter, r *http.Request) {
	logs.Info(r.URL.Path)

	logs.Info(r.URL.Path)
	ctx, conn := easygo.GrpcFromHttp("localhost:9192", r)
	defer conn.Close()
	c := gateway.NewGatewayClient(conn)

	name := "test3"

	res, err := c.Echo(ctx, &gateway.StringMessage{Value: name})
	if err != nil {
		logs.Error("could not greet: %v", err)
	}
	logs.Info("rpc Request: %s", res.GetValue())
	logs.Info("ok")
}

func ServerHttp(w http.ResponseWriter, r *http.Request) {
	logs.Info(r.URL.Path)
	// str := `{"value":" test"}`
	// url := "http://localhost:9191/v1/example/echo"
	// res := easygo.HttpFromHttp(r, url, "POST", strings.NewReader(str))
	// defer res.Body.Close()
	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	logs.Error(err)
	// }

	// logs.Info("http post" + string(body))
	// logs.Info("http get ok")

	resp := easygo.HttpFromHttp(r, "http://localhost:7072/rpc", "GET", nil)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	easygo.PanicError(err)

	logs.Info("http get" + string(body))
}

func main() {
	easygo.NewJaegerTracer("test3", "127.0.0.1:6831")

	next := easygo.ServerHttpHandler()

	route := http.NewServeMux()
	route.HandleFunc("/http", ServerHttp)

	route.HandleFunc("/rpc", ServerRpc)

	err := http.ListenAndServe(listenAddr, next(route))
	easygo.PanicError(err)
}
