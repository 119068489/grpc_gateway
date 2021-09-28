package main

import (
	"flag"
	"grpc_gateway/easygo"
	"grpc_gateway/gateway"
	"os"
)

func main() {
	defer easygo.PanicWriter.Flush()
	flagSet := flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	gateway.Entry(flagSet, os.Args[1:])
}
