package main

import (
	"grpc_gateway/easygo"
	"time"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var PulsarM *easygo.PulsarManager

func init() {
	PulsarM = easygo.NewPulsarManager()
}

// 基于sarama第三方库开发的kafka client product

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "auto-log"
	msg.Value = sarama.StringEncoder("this is a web log test")
	// 连接kafka
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		logs.Error("producer closed, err:", err)
		return
	}
	defer client.Close()

	mb, _ := msg.Value.Encode()
	for i := 1; i <= 5; i++ {
		ms := string(mb) + " " + easygo.AnytoA(i)
		msg.Value = sarama.StringEncoder(ms)
		// 发送消息
		pid, offset, err := client.SendMessage(msg)
		if err != nil {
			logs.Error("send msg failed, err:", err)
			return
		}
		logs.Info("pid:%v offset:%v\n", pid, offset)
		time.Sleep(1 * time.Second)
	}

}
