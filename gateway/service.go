package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"time"

	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"grpc_gateway/easygo"

	"github.com/astaxie/beego/logs"
	"github.com/minio/minio-go"
)

func proxyEntry(w http.ResponseWriter, r *http.Request) {
	logs.Debug("一个RPC转发请求")
	trueServer := "http://localhost:9191"
	url, err := url.Parse(trueServer)
	if err != nil {
		log.Println(err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)
}

func proxyUpEntry(w http.ResponseWriter, r *http.Request) {
	logs.Debug("一个上传转发请求")
	trueServer := "http://127.0.0.1:10086"
	url, err := url.Parse(trueServer)
	if err != nil {
		log.Println(err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)
}

func upEntry(wi http.ResponseWriter, ri *http.Request) {
	ri.ParseMultipartForm(32 << 20)
	file, _, err := ri.FormFile("uploadfile")
	if err != nil {
		log.Println(err)
		wi.Write([]byte("Error:Upload Error."))
		return
	}
	defer file.Close()

	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()
		part, err := m.CreateFormFile("uploadfile", "filename")
		if err != nil {
			log.Println(err)
		}
		if _, err = io.Copy(part, file); err != nil {
			log.Println(err)
		}
	}()

	res, errq := http.Post("http://127.0.0.1:10086", m.FormDataContentType(), r)
	if errq != nil {
		log.Println(errq)
	}
	body, errr := ioutil.ReadAll(res.Body)
	if errr != nil {
		log.Println(errr)
	}

	logs.Debug(string(body))

	wi.Write(body)
}

func getEntry(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<!doctype html><html><head><meta charset=\"utf-8\"><title>Upload</title></head><body><form method=\"POST\" action=\"\" enctype=\"multipart/form-data\"><input name=\"uploadfile\" type=\"file\" /><input type=\"submit\" value=\"Upload\" /></form></body></html>"))
}

func upMinio(w http.ResponseWriter, r *http.Request) {
	endpoint := "127.0.0.1:9000"
	accessKeyID := "VP9R4AYZDA2WO0T9RCZ3"
	secretAccessKey := "YsmeLkmszaLpbEFG0EsZbgcHZsAKGxg+Cb8CYeNe"
	useSSL := false

	// 初使化minio client对象。
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	// 创建一个叫mymusic的存储桶。
	bucketName := "mymusic"
	location := "us-east-1"

	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// 检查存储桶是否已经存在。
		exists, err := minioClient.BucketExists(bucketName)
		if err != nil && !exists {
			log.Fatalln(err)

		}
	}

	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("uploadfile")
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Upload Error."))
		return
	}
	defer file.Close()

	objectName := easygo.AnytoA(time.Now().Unix())

	n, err := minioClient.PutObject(bucketName, objectName, file, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		fmt.Println(err)
		return
	}

	type Result struct {
		BucketName string
		ObjectName string
		Size       int64
	}

	result := &Result{bucketName, objectName, n}
	b, err := json.Marshal(result)
	if err != nil {
		return
	}
	w.Write(b)
}
