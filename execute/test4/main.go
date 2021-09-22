package main

import (
	"context"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"google.golang.org/grpc/metadata"
)

func main() {
	mds := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	logs.Debug(mds)
	ctx := metadata.NewOutgoingContext(context.Background(), mds)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Printf("get metadata error")
	}
	if t, ok := md["timestamp"]; ok {
		fmt.Printf("timestamp from metadata:\n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	}

	// logs.Debug(mds["timestamp"])
}
