package main

import (
	"grpc_gateway/easygo"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		for {
			easygo.StringSliceEqual([]string{"hello", "goconvey"}, []string{"hello", "goconvey"})
		}
	}()

	http.ListenAndServe("0.0.0.0:6060", nil)
}
