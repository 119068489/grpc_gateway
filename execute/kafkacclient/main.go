package main

import (
	"grpc_gateway/easygo"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var PulsarM *easygo.PulsarManager

func init() {
	PulsarM = easygo.NewPulsarManager()
}

// 基于sarama第三方库开发的kafka client consumer

func main() {
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		logs.Error("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions("auto-log") // 根据topic取到所有的分区
	if err != nil {
		logs.Error("fail to get list of partition:err%v\n", err)
		return
	}
	forever := make(chan bool)
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("auto-log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			logs.Error("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				logs.Info("Partition:%d Offset:%d Key:%v Value:%s", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
			}
		}(pc)
	}
	<-forever
}
