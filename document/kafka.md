# kafka

## kafka介绍

Kafka是一个分布式的、分区化、可复制提交的日志服务

还可用于不同应用程序之间的松耦和，作为一个可扩展、高可靠的消息系统，还支持高吞吐量的应用

scale out：无需停机即可扩展机器

持久化：通过将数据持久化到硬盘以及replication防止数据丢失

支持online和offline的场景

Kafka通过Zookeeper管理集群配置，选举leader，以及在Consumer Group发生变化时进行rebalance。因此需要依赖Zookeeper

* 常用的场景
  1. 监控
  2. 消息队列
  3. 站点用户活动追踪
  4. 流处理
  5. 日志聚合
  6. 持久性日志

* Kafka中包含以下基础概念
   1. Topic(话题)：Kafka中用于区分不同类别信息的类别名称。由producer指定
   2. Producer(生产者)：将消息发布到Kafka特定的Topic的对象(过程)
   3. Consumers(消费者)：订阅并处理特定的Topic中的消息的对象(过程)
   4. Broker(Kafka服务集群)：已发布的消息保存在一组服务器中，称之为Kafka集群。集群中的每一个服务器都是一个代理(Broker). 消费者可以订阅一个或多个话题，并从Broker拉数据，从而消费这些已发布的消息。
   5. Partition(分区)：Topic物理上的分组，一个topic可以分为多个partition，每个partition是一个有序的队列。partition中的每条消息都会被分配一个有序的id（offset）
   6. Message：消息，是通信的基本单位，每个producer可以向一个topic（主题）发布一些消息。

* 工作流
   1. ⽣产者从Kafka集群获取分区leader信息
   2. ⽣产者将消息发送给leader
   3. leader将消息写入本地磁盘
   4. follower从leader拉取消息数据
   5. follower将消息写入本地磁盘后向leader发送ACK
   6. leader收到所有的follower的ACK之后向生产者发送ACK

* kafka选择partition的原则
  1. partition在写入的时候可以指定需要写入的partition，如果有指定，则写入对应的partition。
  2. 如果没有指定partition，但是设置了数据的key，则会根据key的值hash出一个partition。
  3. 如果既没指定partition，又没有设置key，则会采用轮询⽅式，即每次取一小段时间的数据写入某个partition，下一小段的时间写入下一个partition

* ACK应答机制
  producer在向kafka写入消息的时候，可以设置参数来确定是否确认kafka接收到数据，这个参数可设置 的值为 0,1,all

  0 代表producer往集群发送数据不需要等到集群的返回，不确保消息发送成功。安全性最低但是效率最高。
  1 代表producer往集群发送数据只要leader应答就可以发送下一条，只确保leader发送成功。
  all 代表producer往集群发送数据需要所有的follower都完成从leader的同步才会发送下一条，确保leader发送成功和所有的副本都完成备份。安全性最⾼，但是效率最低
  
以下实例以windows本地调试为例

## 安装
版本[下载](http://kafka.apache.org/downloads.html)

安装文档：[官方文档](http://kafka.apache.org/documentation/#quickstart)

前置条件要有JDK8+的环境

解压安装包kafka_2.13-2.8.0.tgz，目录如下：

kafka_2.13-2.8.0
  ├── bin
  ├── config
  ├── libs
  ├── licenses
  ├── logs
  ├── site-docs
  ... LICENSE
  ... NOTICE

## 启动服务
1. 启动ZooKeeper
   打开kafka_2.12-2.1.0\bin\windows目录，在此目录下打开cmd，执行命令
   `zookeeper-server-start.bat ..\..\config\zookeeper.properties`

2. 启动Kafka
   打开kafka_2.12-2.1.0\bin\windows目录，在此目录下打开cmd，执行命令
   `kafka-server-start.bat ..\..\config\server.properties`

### 注意：
如果出现‘命令语法不正确’ ，导致不能正常运行，尝试修改配置文件的dataDir（zookeeper.properties），log.dirs（server.properties）。因为默认的是linux的文件目录格式。

如果出现‘命令过长’，导致不能正常运行，请将安装包解压的目录放到盘符根目录，比如：F:\kafka_2.13-2.8.0

## 命令测试

_注意：不管是用命令测试，还是用代码测试，必须先创建topic主题_

* 创建主题
  
  创建一个名称为test-log的主题，2个分区，1个副本，尽量不要用xx_topic这样的名字，容易和命令搞混，导致不必要的错误。
  `kafka-topics.bat --create --zookeeper localhost:2181 --replication-factor 1 --partitions 2 --topic test-log`

  参数说明：
    –zookeeper ：为kafka所配置的zk地址，多个zk以“，”号分隔开，用于存放kafka服务的主题元数据信息及节点信息。
    –partitions ：用于设置主题的分区数，每一个线程负责一个分区
    –replication-factor ：用于设置主题的副本数，每个副本分配在不同的节点，但不能超过总结点数，比如上面的示例若副本设置为2就会报错。


  如创建成功，命令行会有提示输出：
  ```
  F:\kafka_2.13-2.8.0\bin\windows>kafka-topics.bat --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test-log
  Created topic test-log.
  ```

* 修改主题
  - 为主题增加配置
    `kafka-topics.bat --alter --zookeeper localhost:2181 --topic test-log --config flush.messages=1`

  - 删除新增的配置
    `kafka-topics.bat --alter --zookeeper localhost:2181 --topic test-log --delete-config flush.messages`


* 删除主题
  `kafka-topics.bat --delete --zookeeper localhost:2181 --topic test-log`

  1. 若delete.topic.enable=true直接彻底删除该Topic。
  2. 若delete.topic.enable=false若当前被删除的topic没有使用过即没有传输过信息，可以彻底删除； 反之，则不会把这个topic直接删除，而是将其标记为（marked for deletion），重启kafka server后删除。
  

  
* 查看主题列表
  `kafka-topics.bat --list --zookeeper localhost:2181`

  执行结果如下：
  ```
  F:\kafka_2.13-2.8.0\bin\windows>kafka-topics.bat --list --zookeeper localhost:2181
  __consumer_offsets
  test-log
  ```

* 查看主题详情
  `kafka-topics.bat --zookeeper localhost:2181 --describe --topic test-log`

  执行结果如下：
  ```
  F:\kafka_2.13-2.8.0\bin\windows>kafka-topics.bat --zookeeper localhost:2181 --describe --topic test-log
  Topic: test-log TopicId: m5pfj8nXQXyz-2Ixn5WE3Q PartitionCount: 1       ReplicationFactor: 1    Configs:
          Topic: test-log Partition: 0    Leader: 0       Replicas: 0     Isr: 0
  ```

  
* 启动生产者
  `kafka-console-producer.bat --broker-list localhost:9092 --topic test-log`

  如果不报错，进入待输入状态，可以持续发送消息
  ```
  F:\kafka_2.13-2.8.0\bin\windows>kafka-console-producer.bat --broker-list localhost:9092 --topic test-log
  >hello kafka
  >hello
  >ok
  >test
  >
  ```
  
* 启动消费者
  `kafka-console-consumer.bat --bootstrap-server localhost:9092 --from-beginning --topic test-log`

  如果不报错，进入待输入状态，可以持续收到消息
  ```
  F:\kafka_2.13-2.8.0\bin\windows>kafka-console-consumer.bat --bootstrap-server localhost:9092 --from-beginning --topic test-log
  hello kafka
  ok
  hello
  test
  ```

## go代码测试
1. 生产者
   ```
   package main

   import (
   	"grpc_gateway/easygo"

   	"github.com/Shopify/sarama"
   	"github.com/astaxie/beego/logs"
   )

   // 基于sarama第三方库开发的kafka client product

   func main() {
   	config := sarama.NewConfig()
   	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
   	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
   	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

   	// 构造一个消息
   	msg := &sarama.ProducerMessage{}
   	msg.Topic = "test-log"
   	msg.Value = sarama.StringEncoder("this is a web log test")
   	// 连接kafka
   	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
   	if err != nil {
   		logs.Error("producer closed, err:", err)
   		return
   	}
   	defer client.Close()
   	// 发送消息
   	pid, offset, err := client.SendMessage(msg)
   	if err != nil {
   		logs.Error("send msg failed, err:", err)
   		return
   	}
   	logs.Info("pid:%v offset:%v\n", pid, offset)
   }
   ```

2. 消费者
   ```
   package main

   import (
   	"grpc_gateway/easygo"
   	"time"

   	"github.com/Shopify/sarama"
   	"github.com/astaxie/beego/logs"
   )

   // 基于sarama第三方库开发的kafka client consumer

   func main() {
   	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
   	if err != nil {
   		logs.Error("fail to start consumer, err:%v\n", err)
   		return
   	}
   	partitionList, err := consumer.Partitions("test-log") // 根据topic取到所有的分区
   	if err != nil {
   		logs.Error("fail to get list of partition:err%v\n", err)
   		return
   	}
   	logs.Info(len(partitionList))
   	for partition := range partitionList { // 遍历所有的分区
   		// 针对每个分区创建一个对应的分区消费者
   		pc, err := consumer.ConsumePartition("test-log", int32(partition), sarama.OffsetNewest)
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
   		// time.Sleep(1 * time.Hour)
   	}
   }
   ```
