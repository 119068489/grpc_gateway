package main

import (
	"grpc_gateway/easygo"
	"runtime"
	"testing"
	"time"
)

func BenchmarkMain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		easygo.PrintMsg("asdddddddaaaaaaaasefefefefefeesdfafsafaf")
	}
	time.Sleep(time.Second * 1)
	b.Log(runtime.NumGoroutine())
}

// BenchmarkMain-8   	       1	8627570600 ns/op	   85328 B/op	    2339 allocs/op
