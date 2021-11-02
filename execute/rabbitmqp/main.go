package main

import (
	"grpc_gateway/easygo"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
)

var Rabbitmq *easygo.RabbitMQ
var Rabbitmq2 *easygo.RabbitMQ

const MQURL = "amqp://guest:guest@localhost:5672/" //"amqp://testuser:123456@127.0.0.1:5672/testhost"

func init() {
	Rabbitmq = easygo.NewRabbitMQ()
	Rabbitmq.Mqurl = MQURL

	Rabbitmq2 = easygo.NewRabbitMQ()
	Rabbitmq2.Mqurl = MQURL
}

func main() {
	// publishSimple()
	// publishwork()
	// publishPub()
	// publishRouting()
	// publishTopic()
	publishRpc()
}

//simple 简单模式生产者
func publishSimple() {
	Rabbitmq.QueueName = "testhost"
	Rabbitmq.NewClient()

	Rabbitmq.PublishSimple("Hello testuser!")
	logs.Info(" 发送成功！")
}

//work 工作模式生产者
func publishwork() {
	Rabbitmq.QueueName = "testhost"
	Rabbitmq.NewClient()

	for i := 1; i <= 1000; i++ {
		Rabbitmq.PublishSimple("Hello testuser!" + strconv.Itoa(i))
		time.Sleep(time.Second / 10)
		logs.Info("[%d] 发送成功！", i)
	}
}

// Publish 订阅模式生产者
func publishPub() {
	Rabbitmq.Exchange = "newProduct"
	Rabbitmq.NewClient()

	for i := 1; i <= 10; i++ {
		Rabbitmq.PublishPub("Hello Publish mode!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		logs.Info("[%d] 发送成功！", i)
	}
}

// Routing模式
func publishRouting() {
	Rabbitmq.Exchange = "kuteng"
	Rabbitmq.Key = "kuteng_one"
	Rabbitmq.NewClient()

	Rabbitmq2.Exchange = "kuteng"
	Rabbitmq2.Key = "kuteng_two"
	Rabbitmq2.NewClient()

	for i := 0; i <= 10; i++ {
		Rabbitmq.PublishRouting("Hello kuteng one!" + strconv.Itoa(i))
		Rabbitmq2.PublishRouting("Hello kuteng Two!" + strconv.Itoa(i))

		time.Sleep(1 * time.Second)
		logs.Info("[%d] 发送成功！", i)
	}
}

// Topic模式
func publishTopic() {
	Rabbitmq.Exchange = "exKutengTopic"
	Rabbitmq.Key = "kuteng.topic.one"
	Rabbitmq.NewClient()

	Rabbitmq2.Exchange = "exKutengTopic"
	Rabbitmq2.Key = "kuteng.topic.two"
	Rabbitmq2.NewClient()

	for i := 0; i <= 10; i++ {
		Rabbitmq.PublishTopic("Hello kuteng one!" + strconv.Itoa(i))
		Rabbitmq2.PublishTopic("Hello kuteng Two!" + strconv.Itoa(i))

		time.Sleep(1 * time.Second)
		logs.Info("[%d] 发送成功！", i)
	}
}

//Rpc模式
func publishRpc() {
	Rabbitmq.QueueName = "rpc_queue"
	Rabbitmq.NewClient()

	Rabbitmq.PublishRpc("Hello rpcuser!")

}
