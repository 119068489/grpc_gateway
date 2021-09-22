package main

import (
	"context"
	"grpc_gateway/easygo"
	"log"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/astaxie/beego/logs"
)

var PulsarM *easygo.PulsarManager

func init() {
	PulsarM = easygo.NewPulsarManager()
}

func main() {

	PulsarM.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://172.16.3.158:6650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	defer PulsarM.GetClient().Close()

	PulsarM.Subscribe(pulsar.ConsumerOptions{
		Topic:            "my-topic",
		SubscriptionName: "my-sub",
		Type:             pulsar.Shared,
	})
	defer PulsarM.GetConsumer().Close()

	for i := 0; i < 10; i++ {
		msg := PulsarM.Receive(context.Background())

		logs.Info("Received message msgId: %#v -- content: '%s'\n",
			msg.ID().EntryID(), string(msg.Payload()))

		PulsarM.Ack(msg)
	}

	if err := PulsarM.GetConsumer().Unsubscribe(); err != nil {
		log.Fatal(err)
	}

}
