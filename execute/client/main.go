package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"grpc_gateway/api"
	"grpc_gateway/proto/pb/gateway"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	address     = "localhost:9191"
	defaultName = "world"
	aginName    = "agin"
	downFile    = "./testFile.txt"
)

func main() {
	// for {
	// 	<-time.Tick(time.Second)
	// 	easygo.ProtectRun(httpPost)
	// 	easygo.ProtectRun(httpGet)
	// 	easygo.ProtectRun(rpcReq)
	// }

	// easygo.ProtectRun(UploadPost)
	// easygo.ProtectRun(rpcUpload)
	// easygo.ProtectRun(UploadFilePost)
	// test()
	// result := login("admin", "123456")
	// ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", result.Token)
	// rpcReq(ctx)
	// httpPost(result.Token)

	c := NewClient()
	c.Run()
}

type Client struct {
	Token string
}

func NewClient() *Client {
	p := &Client{}
	p.Init()
	return p
}

//初始化
func (c *Client) Init() {
}

func (c *Client) doCmd(params []string) {
	fname := []rune(params[0])
	fname[0] = []rune(strings.ToUpper(string(fname[0])))[0] //把首字母转换成大写
	methodName := string(fname)
	method := reflect.ValueOf(c).MethodByName(methodName)
	if !method.IsValid() || method.Kind() != reflect.Func {
		logs.Info("%v 不能识别的命令,方法没有实现", methodName)
		return
	}
	args := make([]reflect.Value, 0, len(params)-1)
	for _, para := range params[1:] {
		v := reflect.ValueOf(para)
		args = append(args, v)
	}

	method.Call(args) // 反射调用方法
}

func (c *Client) Run() {
	for {
		logs.Info("please input")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		s := input.Text()
		if s != "" {
			params := strings.Split(s, " ")
			c.doCmd(params)
		}
	}
}

func (c *Client) Post() {
	str := `{"value":" ` + defaultName + `"}`
	req, err := http.NewRequest(http.MethodPost, "http://localhost:9190/v1/example/echo", strings.NewReader(str))
	if err != nil {
		fmt.Println(err)
	}
	if c.Token != "" {
		req.Header.Set("authorization", c.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Fatal(err)
	}

	fmt.Println("http post" + string(body))
}

//get gcho world 101
func (c *Client) Get(params ...string) {
	host := "http://localhost:9190/v1/example/"
	var par string
	for i := range params {
		par = path.Join(par, params[i])
	}
	host += par

	resp, err := http.Get(host)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Fatal(err)
	}

	fmt.Println("http get " + string(body))
}

func (c *Client) Login(u, p string) *gateway.LoginReply {
	str := `{"username":"` + u + `","password":"` + p + `"}`
	resp, err := http.Post("http://localhost:9190/v1/example/login",
		"application/x-www-form-urlencoded",
		strings.NewReader(str))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Fatal(err)
	}
	result := &gateway.LoginReply{}
	err = json.Unmarshal(body, result)
	if err != nil {
		fmt.Println(err)
	}

	if result.Status == "200" {
		c.Token = result.GetToken()
	}

	fmt.Println("login", result.Status)
	return result
}

func (c *Client) RpcReq(ctx context.Context) {

	if md, ok := metadata.FromIncomingContext(ctx); !ok {
		logs.Error("无Token认证信息")
	} else {
		l := md["authorization"]
		if len(l) > 0 {
			userName := api.CheckAuth(ctx)
			logs.Debug(userName, l[0])
		}
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logs.Error("did not connect: %v", err)
	}
	defer conn.Close()
	clt := gateway.NewGatewayClient(conn)

	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*6)
	defer cancel()

	r, err := clt.Echo(ctx, &gateway.StringMessage{Value: name})
	if err != nil {
		logs.Error("could not greet: %v", err)
	}
	logs.Info("rpc return:", r)
	logs.Info("rpc Request: %s", r.GetValue())

	// for i := 1; i < 6; i++ {
	// 	time.Sleep(time.Second * 1)
	// 	r, err = c.Echo(ctx, &pb.StringMessage{Value: aginName})
	// 	if err != nil {
	// 		logs.Error("could not greet: %v", err)
	// 	}
	// 	logs.Info("rpc Request: %s ：%d", r.GetValue(), i)
	// }
}

func (c *Client) RpcUpload() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	clt := gateway.NewGatewayClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	bs, _ := ioutil.ReadFile("./123.txt")
	ff := base64.StdEncoding.EncodeToString(bs)
	r, err := clt.Upload(ctx, &gateway.FSReq{Name: "./testFile.txt", File: ff})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println("rpc return:", r)
	log.Printf("rpc Request: %s", r.GetMessage())
}

func (c *Client) UploadPost() {
	bs, _ := ioutil.ReadFile("./123.txt")
	ff := base64.StdEncoding.EncodeToString(bs)
	str := `{"name":"` + downFile + `","file":"` + ff + `"}`
	resp, err := http.Post("http://localhost:8080/v1/example/upload",
		"application/x-www-form-urlencoded",
		strings.NewReader(str))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Fatal(err)
	}

	fmt.Println("http post" + string(body))
}

func (c *Client) UploadFilePost() {
	client := http.Client{}
	bodyBuf := &bytes.Buffer{}
	bodyWrite := multipart.NewWriter(bodyBuf)
	file, err := os.Open("./test.png")
	if err != nil {
		log.Println("err")
	}
	defer file.Close()
	if err != nil {
		log.Println("err")
	}
	// file 为key
	fileWrite, err := bodyWrite.CreateFormFile("uploadfile", "file")
	if err != nil {
		log.Println("err")
	}
	_, err = io.Copy(fileWrite, file)
	if err != nil {
		log.Println("err")
	}
	bodyWrite.Close() //要关闭，会将w.w.boundary刷写到w.writer中
	// 创建请求
	contentType := bodyWrite.FormDataContentType()
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:9190", bodyBuf)
	if err != nil {
		log.Println("err")
	}
	// 设置头
	req.Header.Set("Content-Type", contentType)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("err")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err")
	}
	fmt.Println(string(b))
}
