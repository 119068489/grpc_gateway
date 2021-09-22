package main

import (
	"context"
	"grpc_gateway/easygo"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
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

	PulsarM.NewProducer(pulsar.ProducerOptions{Topic: "my-topic"})
	defer PulsarM.GetProducer().Close()

	PulsarM.Send(context.Background(), &pulsar.ProducerMessage{
		// DeliverAfter: 10 * time.Second, //延时投递
		Payload: []byte("hello"),
	})

}
