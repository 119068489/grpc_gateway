package main

import (
	"flag"
	"os"

	"grpc_gateway/easygo"
	"grpc_gateway/rpcserver"
)

func main() {
	defer easygo.PanicWriter.Flush()
	defer easygo.RecoverAndLog()

	flagSet := flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	rpcserver.Entry(flagSet, os.Args[1:])
}
