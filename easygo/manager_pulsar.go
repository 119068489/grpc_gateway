package easygo

import (
	"context"
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
)

type IPulsarManager interface {
	GetClient() pulsar.Client
	GetProducer() pulsar.Producer
	GetConsumer() pulsar.Consumer
	GetRead() pulsar.Reader
}

//连接管理
type PulsarManager struct {
	Client   *pulsar.Client
	Producer *pulsar.Producer
	Consumer *pulsar.Consumer
	Reader   *pulsar.Reader
}

func NewPulsarManager() *PulsarManager { // services map[string]interface{},
	p := &PulsarManager{}
	p.Init()
	return p
}

//初始化
func (p *PulsarManager) Init() {
}

// 创建一个pulsar Client实例
func (p *PulsarManager) NewClient(copt pulsar.ClientOptions) {
	client, err := pulsar.NewClient(copt)
	if err != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", err)
	}

	p.Client = &client
}

// 创建Producer实例
// 此方法会阻塞直到生产者创建成功
func (p *PulsarManager) NewProducer(popt pulsar.ProducerOptions) {
	producer, err := p.GetClient().CreateProducer(popt)
	if err != nil {
		log.Fatal(err)
	}

	p.Producer = &producer
}

// 发送消息
// 在 Pulsar broker 成功确认之前，此调用将被阻塞。
// 示例： producer.Send(ctx, pulsar.ProducerMessage{ Payload: myPayload })
func (p *PulsarManager) Send(ctx context.Context, message *pulsar.ProducerMessage) {
	_, err := p.GetProducer().Send(ctx, message)

	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}
}

// 订阅 通过订阅一个主题来创建一个 `Consumer`。
// 如果订阅不存在，将创建一个新的订阅，所有在创建之后发布的消息将被保留直到被确认，即使消费者没有连接
func (p *PulsarManager) Subscribe(copt pulsar.ConsumerOptions) {
	consumer, err := p.GetClient().Subscribe(copt)
	if err != nil {
		log.Fatal(err)
	}

	p.Consumer = &consumer
}

// 接收一条消息。
// 这会调用阻塞直到消息可用。(收到消息)
func (p *PulsarManager) Receive(ctx context.Context) pulsar.Message {
	msg, err := p.GetConsumer().Receive(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return msg
}

// 确认单个消息的消费
func (p *PulsarManager) Ack(msg pulsar.Message) {
	p.GetConsumer().Ack(msg)
}

func (p *PulsarManager) NewReader(ropt pulsar.ReaderOptions) {
	reader, err := p.GetClient().CreateReader(ropt)
	if err != nil {
		log.Fatal(err)
	}
	p.Reader = &reader
}

func (p *PulsarManager) Next(ctx context.Context) pulsar.Message {
	msg, err := p.GetRead().Next(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}

func (p *PulsarManager) GetClient() pulsar.Client {
	return *p.Client
}

func (p *PulsarManager) GetProducer() pulsar.Producer {
	return *p.Producer
}

func (p *PulsarManager) GetConsumer() pulsar.Consumer {
	return *p.Consumer
}

func (p *PulsarManager) GetRead() pulsar.Reader {
	return *p.Reader
}
