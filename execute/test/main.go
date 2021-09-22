package main

import (
	"flag"
	"grpc_gateway/easygo"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/SkyAPM/go2sky"
	httpPlugin "github.com/SkyAPM/go2sky/plugins/http"
	"github.com/SkyAPM/go2sky/reporter"
)

var (
	grpc        bool
	oapServer   string
	upstreamURL string
	listenAddr  string
	serviceName string

	client *http.Client
)

func init() {
	flag.BoolVar(&grpc, "grpc", false, "use grpc reporter")
	flag.StringVar(&oapServer, "oap-server", "0.0.0.0:11800", "oap server address")
	flag.StringVar(&upstreamURL, "upstream-url", "http://localhost:9090/v1/example/echo", "upstream service url")
	flag.StringVar(&listenAddr, "listen-addr", ":7070", "listen address")
	flag.StringVar(&serviceName, "service-name", "go2sky", "service name")
}

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
