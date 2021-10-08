package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"grpc_gateway/easygo"
	pb "grpc_gateway/proto/pb/gateway"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/golang/glog"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:9191"
	defaultName = "world"
	aginName    = "agin"
	downFile    = "./testFile.txt"
)

func main() {
	for {
		<-time.NewTicker(time.Second).C
		easygo.ProtectRun(httpPost)
		easygo.ProtectRun(httpGet)
		easygo.ProtectRun(rpcReq)
	}

	// easygo.ProtectRun(UploadPost)
	// easygo.ProtectRun(rpcUpload)
	// easygo.ProtectRun(UploadFilePost)
	// test()
}

func httpPost() {
	str := `{"value":" ` + defaultName + `"}`
	resp, err := http.Post("http://localhost:9190/v1/example/echo",
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

func httpGet() {
	resp, err := http.Get("http://localhost:9190/v1/example/gcho/world/101")
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

func rpcReq() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logs.Error("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGatewayClient(conn)

	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()
	r, err := c.Echo(ctx, &pb.StringMessage{Value: name})
	if err != nil {
		logs.Error("could not greet: %v", err)
	}
	logs.Info("rpc return:", r)
	logs.Info("rpc Request: %s", r.GetValue())

	for i := 1; i < 6; i++ {
		time.Sleep(time.Second * 1)
		r, err = c.Echo(ctx, &pb.StringMessage{Value: aginName})
		if err != nil {
			logs.Error("could not greet: %v", err)
		}
		logs.Info("rpc Request: %s ：%d", r.GetValue(), i)
	}
}

func rpcUpload() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGatewayClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	bs, _ := ioutil.ReadFile("./123.txt")
	ff := base64.StdEncoding.EncodeToString(bs)
	r, err := c.Upload(ctx, &pb.FSReq{Name: "./testFile.txt", File: ff})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println("rpc return:", r)
	log.Printf("rpc Request: %s", r.GetMessage())
}

func UploadPost() {
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

func UploadFilePost() {
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

func test() {
	mm := make(chan string)
	// f := make(chan bool)

	go func() {
		i := 0
		for {
			time.Sleep(time.Duration(i) * time.Second)
			ms := easygo.AnytoA(i)
			mm <- ms
			i++
		}

	}()

	for m := range mm {
		logs.Debug("m:", m)
	}

	// <-f
}
