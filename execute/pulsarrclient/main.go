package main

import (
	"context"
	"fmt"
	"grpc_gateway/easygo"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/astaxie/beego/logs"
)

var PulsarM *easygo.PulsarManager

func init() {
	PulsarM = easygo.NewPulsarManager()
}

func main() {

	cOpt := pulsar.ClientOptions{
		URL:               "pulsar://172.16.3.158:6650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	}

	PulsarM.NewClient(cOpt)
	defer PulsarM.GetClient().Close()

	PulsarM.NewReader(pulsar.ReaderOptions{
		Topic: "my-topic",
		// StartMessageID: msgID,
		StartMessageID: pulsar.EarliestMessageID(),
	})
	defer PulsarM.GetRead().Close()

	logs.Debug("first:", pulsar.EarliestMessageID())
	logs.Debug("last", pulsar.LatestMessageID())

	for PulsarM.GetRead().HasNext() {
		msg := PulsarM.Next(context.Background())

		fmt.Printf("Received message msgId: %#v -- content: '%s'\n",
			msg.ID().EntryID(), string(msg.Payload()))
	}

}
