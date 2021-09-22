package main

import (
	"grpc_gateway/easygo"

	"github.com/astaxie/beego/logs"
	"github.com/streadway/amqp"
)

var Rabbitmq *easygo.RabbitMQ

const MQURL = "amqp://guest:guest@localhost:5672/" //"amqp://testuser:123456@127.0.0.1:5672/testhost"

func init() {
	Rabbitmq = easygo.NewRabbitMQ()
	Rabbitmq.Mqurl = MQURL
}

func main() {
	// consume()
	// consumeSub()
	// recieveRouting()
	// recieveRouting2()
	// recieveTopic()
	// recieveTopic2()
	RecieveRpc()
}

//消费者消费消息的函数
func ReadMsg(d amqp.Delivery) {
	logs.Info("Received a message: %s", d.Body)
}

//简单模式 工作模式 接收消息
func consumeSimple() {
	Rabbitmq.QueueName = "testhost" //随便写，生产消费一致就行
	Rabbitmq.NewClient()
	defer Rabbitmq.Destory()
	Rabbitmq.ConsumeSimple(ReadMsg)
}

//订阅模式 接收消息
func consumeSub() {
	Rabbitmq.Exchange = "newProduct" //随便写，生产消费一致就行
	Rabbitmq.NewClient()
	Rabbitmq.RecieveSub(ReadMsg)
}

//路由模式 接收消息
func recieveRouting() {
	Rabbitmq.Exchange = "kuteng"
	Rabbitmq.Key = "kuteng_one"
	Rabbitmq.NewClient()
	Rabbitmq.RecieveRouting(ReadMsg)
}

//路由模式 接收消息
func recieveRouting2() {
	Rabbitmq.Exchange = "kuteng"
	Rabbitmq.Key = "kuteng_two"
	Rabbitmq.NewClient()
	Rabbitmq.RecieveRouting(ReadMsg)
}

//话题模式 接收消息
func recieveTopic() {
	Rabbitmq.Exchange = "exKutengTopic"
	Rabbitmq.Key = "#"
	Rabbitmq.NewClient()
	Rabbitmq.RecieveTopic(ReadMsg)
}

//话题模式 接收消息
func recieveTopic2() {
	Rabbitmq.Exchange = "exKutengTopic"
	Rabbitmq.Key = "kuteng.*.two"
	Rabbitmq.NewClient()
	Rabbitmq.RecieveTopic(ReadMsg)
}

//rpc模式 请求消息
func RecieveRpc() {
	Rabbitmq.Key = "rpc_queue"
	Rabbitmq.NewClient()

	Rabbitmq.RecieveRpc("ok")
}
