# PULSAR

Pulsar是由雅虎创建的开源的、分布式pub-sub系统，现在是Apache基金会的一个孵化项目。

## 安装单机模式Pulsar

- 系统要求
  
  目前，Pulsar 可用于 64 位 macOS、Linux 和 Windows。 要使用 Pulsar，您需要安装 64 位 JRE/JDK 8 或更高版本。

  默认情况下，Pulsar 会分配 2G JVM 堆内存来启动。 它可以在 PULSAR_MEM 下的 conf/pulsar_env.sh 文件中更改。 这是传递给 JVM 的额外选项。

- 安装二进制版本 Pulsar
  
  从 Pulsar [官网下载](https://pulsar.apache.org/download)

- 启动单机模式 Pulsar
  `bin/pulsar standalone`

  成功启动 Pulsar 后，可以看到如下所示的 INFO 级日志消息：
  ```
  2017-06-01 14:46:29,192 - INFO  - [main:WebSocketService@95] - Configuration Store cache started
  2017-06-01 14:46:29,192 - INFO  - [main:AuthenticationService@61] - Authentication is disabled
  2017-06-01 14:46:29,192 - INFO  - [main:WebSocketService@108] - Pulsar WebSocket Service started
  ```

- 使用单机模式 Pulsar
  Pulsar 中有一个名为 pulsar-client 的 CLI 工具。 Pulsar-client 工具允许使用者在运行的集群中 consume 并 produce 消息到 Pulsar topic。

  - Consume 一条消息
    在 first-subscription 订阅中 consume 一条消息到 my-topic 的命令如下所示：
    `bin/pulsar-client consume my-topic -s "first-subscription"`

  - Produce 一条消息
    向名称为 my-topic 的 topic 发送一条简单的消息 hello-pulsar，命令如下所示：
    `bin/pulsar-client produce my-topic --messages "hello-pulsar"`

- 终止单机模式 Pulsar
  使用 Ctrl+C 终止单机模式 Pulsar 的运行。

  如果服务使用 pulsar-daemon start standalone 命令作为后台进程运行，则使用 pulsar-daemon stop standalone 命令停止服务。

## Pulsar Go client

### 安装

- 安装 go 工具包
  `go get -u github.com/apache/pulsar-client-go/pulsar`

### 连接 URL
要使用客户端库连接到 Pulsar，您需要指定 Pulsar 协议 URL。
Pulsar 协议 URL 分配给特定集群，使用 pulsar 方案并具有默认端口 6650。以下是 localhost 的示例：

`pulsar://localhost:6650`

如果你有多个 broker，你可以使用下面的方法设置 URl：

`pulsar://localhost:6550,localhost:6651,localhost:6652`

生产 Pulsar 集群的 URL 可能如下所示：

`pulsar://pulsar.us-west.example.com:6650`

如果您使用 TLS 身份验证，则 URL 将如下所示：

`pulsar+ssl://pulsar.us-west.example.com:6651`

### 创建客户端

为了与 Pulsar 交互，您首先需要一个 Client 对象。 您可以使用 NewClient 函数创建一个客户端对象，传入一个 ClientOptions 对象。下面是一个示例：
```
import (
    "log"
    "time"

    "github.com/apache/pulsar-client-go/pulsar"
)

func main() {
    client, err := pulsar.NewClient(pulsar.ClientOptions{
        URL:               "pulsar://localhost:6650",   //单个broker的初始化
        //URL: "pulsar://localhost:6650,localhost:6651,localhost:6652",  //多个broker初始化
        OperationTimeout:  30 * time.Second,
        ConnectionTimeout: 30 * time.Second,
    })
    if err != nil {
        log.Fatalf("Could not instantiate Pulsar client: %v", err)
    }

    defer client.Close()
}
```

### 生产者
Pulsar 生产者向 Pulsar 主题发布消息。 您可以使用 ProducerOptions 对象配置 Go 生产者。下面是一个示例：
```
producer, err := client.CreateProducer(pulsar.ProducerOptions{
    Topic: "my-topic",
})

if err != nil {
    log.Fatal(err)
}

_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
    Payload: []byte("hello"),
})

defer producer.Close()

if err != nil {
    fmt.Println("Failed to publish message", err)
}
fmt.Println("Published message")
```

- Producer operations
  Pulsar Go 生产者有以下可用的方法：
  
  * Topic() string 
    获取生产者的话题

  * Name() string 
    获取生产者的名字

  * Send(context.Context, *ProducerMessage) (MessageID, error) 
    向生产者的主题发布消息。 此调用将阻塞，直到消息被 Pulsar 代理成功确认，或者如果超过了在生产者配置中使用 SendTimeout 设置的超时，则会抛出错误。

  * SendAsync(context.Context, *ProducerMessage, func(MessageID, *ProducerMessage, error))
    发送一条消息，这个调用将被阻塞，直到被 Pulsar broker 成功确认。

  * LastSequenceID() int64
    获取此生产者发布的最后一个序列 ID。 表示由代理发布和确认的自动分配或自定义序列 ID（在 ProducerMessage 上设置）。

  * Flush() error
    刷新客户端中缓存的所有消息，并等待所有消息成功持久化。

  * Close()
    关闭生产者并释放分配给它的所有资源。 如果调用 Close() 则不会再接受来自发布者的消息。 这个方法会一直阻塞，直到所有挂起的发布请求都被 Pulsar 持久化。 如果抛出错误，则不会重试任何挂起的写入。

- Producer 示例
  
  如何在生产者中使用消息路由器
  ```
  client, err := NewClient(pulsar.ClientOptions{
      URL: serviceURL,
  })

  if err != nil {
      log.Fatal(err)
  }
  defer client.Close()

  // 仅订阅特定分区
  consumer, err := client.Subscribe(pulsar.ConsumerOptions{
      Topic:            "my-partitioned-topic-partition-2",
      SubscriptionName: "my-sub",
  })

  if err != nil {
      log.Fatal(err)
  }
  defer consumer.Close()

  producer, err := client.CreateProducer(pulsar.ProducerOptions{
      Topic: "my-partitioned-topic",
      MessageRouter: func(msg *ProducerMessage, tm TopicMetadata) int {
          fmt.Println("Routing message ", msg, " -- Partitions: ", tm.NumPartitions())
          return 2
      },
  })

  if err != nil {
      log.Fatal(err)
  }
  defer producer.Close()
  ```

  生产者如何使用 schema 接口
  ```
  type testJSON struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    }
  ```
  ```
    var (
      exampleSchemaDef = "{\"type\":\"record\",\"name\":\"Example\",\"namespace\":\"test\"," +
          "\"fields\":[{\"name\":\"ID\",\"type\":\"int\"},{\"name\":\"Name\",\"type\":\"string\"}]}"
  )
  ```
  ```
  client, err := NewClient(pulsar.ClientOptions{
      URL: "pulsar://localhost:6650",
  })
  if err != nil {
      log.Fatal(err)
  }
  defer client.Close()

  properties := make(map[string]string)
  properties["pulsar"] = "hello"
  jsonSchemaWithProperties := NewJSONSchema(exampleSchemaDef, properties)
  producer, err := client.CreateProducer(ProducerOptions{
      Topic:  "jsonTopic",
      Schema: jsonSchemaWithProperties,
  })
  assert.Nil(t, err)

  _, err = producer.Send(context.Background(), &ProducerMessage{
      Value: &testJSON{
          ID:   100,
          Name: "pulsar",
      },
  })
  if err != nil {
      log.Fatal(err)
  }
  producer.Close()
  ```

  如何在生产者中使用相应的延迟
  ```
  client, err := NewClient(pulsar.ClientOptions{
      URL: "pulsar://localhost:6650",
  })
  if err != nil {
      log.Fatal(err)
  }
  defer client.Close()

  topicName := newTopicName()
  producer, err := client.CreateProducer(pulsar.ProducerOptions{
      Topic: topicName,
  })
  if err != nil {
      log.Fatal(err)
  }
  defer producer.Close()

  consumer, err := client.Subscribe(pulsar.ConsumerOptions{
      Topic:            topicName,
      SubscriptionName: "subName",
      Type:             Shared,
  })
  if err != nil {
      log.Fatal(err)
  }
  defer consumer.Close()

  ID, err := producer.Send(context.Background(), &pulsar.ProducerMessage{
      Payload:      []byte(fmt.Sprintf("test")),
      DeliverAfter: 3 * time.Second,
  })
  if err != nil {
      log.Fatal(err)
  }
  fmt.Println(ID)

  ctx, canc := context.WithTimeout(context.Background(), 1*time.Second)
  msg, err := consumer.Receive(ctx)
  if err != nil {
      log.Fatal(err)
  }
  fmt.Println(msg.Payload())
  canc()

  ctx, canc = context.WithTimeout(context.Background(), 5*time.Second)
  msg, err = consumer.Receive(ctx)
  if err != nil {
      log.Fatal(err)
  }
  fmt.Println(msg.Payload())
  canc()
  ```

- Producer 配置
  type ProducerOptions struct {
  	// 主题指定此生产者将发布的主题。构造生产者时需要此参数
  	Topic string

  	// Name 为生产者指定一个名称 
    // 如果未分配，系统将生成一个全局唯一的名称，该名称可以通过 // Producer.ProducerName() 访问。 
    // 在指定名称时，由用户来确保对于给定的主题，生产者名称在所有 Pulsar 集群中是唯一的 。 Brokers将强制要求只有一个给定名称的生产者才能在一个主题上发布。
  	Name string

  	// Properties 将附加一组应用程序定义的属性。这个属性将会在主题统计中可见 
  	Properties map[string]string

  	// SendTimeout 设置自发送后未被服务器确认的消息的超时时间。 
    // Send 和 SendAsync 在超时后返回错误。 
    // 默认为 30 秒，负数如 -1 为禁用。
  	SendTimeout time.Duration

    // 控制如果生产者的消息队列已满，Send 和 SendAsync 是否阻塞。 
    // 默认为 false，如果设置为 true，则当队列已满时 Send 和 SendAsync 返回错误
  	DisableBlockIfQueueFull bool

  	// 设置队列的最大大小，该队列包含待处理的消息以接收来自代理的确认。 
  	MaxPendingMessages int

  	// 更改用于选择发布特定消息的分区的 `HashingScheme`。 可用的标准散列函数有： 
    // - `JavaStringHash`：Java String.hashCode() 等效 
    // - `Murmur3_32Hash`：使用 Murmur3 散列函数。 
    // https://en.wikipedia.org/wiki/MurmurHash">https://en.wikipedia.org/wiki/MurmurHash
  	//默认值 `JavaStringHash`.
  	HashingScheme

  	// CompressionType 设置生产者的压缩类型。 
    // 默认情况下，不压缩消息有效负载。 支持的压缩类型有：- LZ4 - ZLIB - ZSTD
    // 注意：从 Pulsar 2.3 开始支持 ZSTD。 消费者至少需要在该版本中才能接收使用 ZSTD 压缩的消息
  	CompressionType

  	// 定义所需的压缩级别。 选项：
  	// - Default
  	// - Faster
  	// - Better
  	CompressionLevel

  	// 通过传递 MessageRouter 的实现来设置自定义消息路由策略 
    // 路由器是一个函数，它给出特定消息和主题元数据，返回消息应该路由到的分区索引
  	MessageRouter func(*ProducerMessage, TopicMetadata) int

    // 控制是否为生产者启用消息的自动批处理。 默认情况下启用批处理 。 
    // 启用批处理时，多次调用 Producer.sendAsync 可以将单个批处理发送到broler，从而提高吞吐量，尤其是在发布小消息时。 
    // 如果启用压缩，消息将在批处理级别进行压缩，从而为类似的标题或内容带来更好的压缩率。 
    // 当启用默认批处理延迟设置为 1 ms 并且默认批处理大小为 1000 条消息时 
    // 设置 `DisableBatching: true` 将使生产者单独发送消息
  	DisableBatching bool

  	// 设置发送的消息将被批处理的时间段（默认值：10ms） 
    // 如果启用了批处理消息。 如果设置为非零值，消息将排队直到这个时间
  	// interval or until
  	BatchingMaxPublishDelay time.Duration

    // 设置批处理中允许的最大消息数。 （默认值：1000） 
    // 如果设置为大于 1 的值，消息将排队直到达到此阈值或BatchingMaxSize（见下文）已达到或批处理间隔已过。
  	BatchingMaxMessages uint

    // 设置批处理中允许的最大字节数。 (默认 128 KB) 
    // 如果设置为大于 1 的值，消息将排队直到达到此阈值或BatchingMaxMessages（见上文）已达到或批处理间隔已过。
  	BatchingMaxSize uint

  	// 拦截器链 这些拦截器在 ProducerInterceptor 接口中定义的某些点被调用。
  	Interceptors ProducerInterceptors

    // Schema 通过传递 Schema 的实现来设置自定义模式类型 | bytes[]
  	Schema Schema

  	// 设置 reconnectToBroker 的最大重试次数。 （默认：ultimate）
  	MaxReconnectToBroker *uint

    // 设置批处理构建器类型（默认 DefaultBatchBuilder） 
    // 这将用于在启用批处理时创建批处理容器。 
    // 选项： - DefaultBatchBuilder - KeyBasedBatchBuilder
  	BatcherBuilderType

  	// 是后台进程发现新分区的时间间隔
  	// 默认 1 minute
  	PartitionsAutoDiscoveryInterval time.Duration
  }

### 消费者
Pulsar 消费者订阅一个或多个 Pulsar 主题并监听在该主题/这些主题上产生的传入消息。 您可以使用 ConsumerOptions 对象配置 Go 消费者。 这是一个使用通道的基本示例：
```
consumer, err := client.Subscribe(pulsar.ConsumerOptions{
    Topic:            "topic-1",
    SubscriptionName: "my-sub",
    Type:             pulsar.Shared,
})
if err != nil {
    log.Fatal(err)
}
defer consumer.Close()

for i := 0; i < 10; i++ {
    msg, err := consumer.Receive(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Received message msgId: %#v -- content: '%s'\n",
        msg.ID(), string(msg.Payload()))

    consumer.Ack(msg)
}

if err := consumer.Unsubscribe(); err != nil {
    log.Fatal(err)
}
```

- Consumers operations
   Pulsar Go 消费者可以使用以下方法：

   * Subscription()	string
     返回消费者的订阅名称

   * Unsubcribe()	error
     从分配的主题中取消订阅消费者。 如果取消订阅操作以某种方式不成功，则会引发错误。

   * Receive(context.Context)	(Message, error)
     从主题接收一条消息。 此方法会阻塞，直到消息可用。

   * Chan()	<-chan ConsumerMessage
     Chan 返回一个传递消息的 channel

   * Ack(Message)	
     向 Pulsar broker 确认消息

   * AckID(MessageID)
     通过消息 ID 向 Pulsar broker 确认消息

   * ReconsumeLater(msg Message, delay time.Duration)
     ReconsumeLater 标记消息以在自定义延迟后重新发送

   * Nack(Message)	
     确认未能处理单个消息。
     
   * NackID(MessageID)	
     确认未能处理单个消息。

   * Seek(msgID MessageID)	error
     将与此消费者关联的订阅重置为特定的消息 ID。 消息 ID 可以是特定消息，也可以表示主题中的第一条或最后一条消息。

   * SeekByTime(time time.Time)	error
     将与此消费者关联的订阅重置为特定的消息发布时间。

   * Close()	
     关闭消费者，禁用其从代理接收消息的能力

   * Name()	string
     Name 返回消费者名称

- Receive 示例
  
  1. 如何使用正则表达式
     ```
     client, err := pulsar.NewClient(pulsar.ClientOptions{
         URL: "pulsar://localhost:6650",
     })

     defer client.Close()

     p, err := client.CreateProducer(pulsar.ProducerOptions{
         Topic:           topicInRegex,
         DisableBatching: true,
     })
     if err != nil {
         log.Fatal(err)
     }
     defer p.Close()

     topicsPattern := fmt.Sprintf("persistent://%s/foo.*", namespace)
     opts := pulsar.ConsumerOptions{
         TopicsPattern:    topicsPattern,
         SubscriptionName: "regex-sub",
     }
     consumer, err := client.Subscribe(opts)
     if err != nil {
         log.Fatal(err)
     }
     defer consumer.Close()
     ```

  2. 如何使用多topic 的Consumer
     ```
     func newTopicName() string {
         return fmt.Sprintf("my-topic-%v", time.Now().Nanosecond())
     }


     topic1 := "topic-1"
     topic2 := "topic-2"

     client, err := NewClient(pulsar.ClientOptions{
         URL: "pulsar://localhost:6650",
     })
     if err != nil {
         log.Fatal(err)
     }
     topics := []string{topic1, topic2}
     consumer, err := client.Subscribe(pulsar.ConsumerOptions{
         Topics:           topics,
         SubscriptionName: "multi-topic-sub",
     })
     if err != nil {
         log.Fatal(err)
     }
     defer consumer.Close()
     ```

  3. 如何使用消费监听器
     ```
      client, err := pulsar.NewClient(pulsar.ClientOptions{URL: "pulsar://localhost:6650"})
       if err != nil {
           log.Fatal(err)
       }

       defer client.Close()

       channel := make(chan pulsar.ConsumerMessage, 100)

       options := pulsar.ConsumerOptions{
           Topic:            "topic-1",
           SubscriptionName: "my-subscription",
           Type:             pulsar.Shared,
       }

       options.MessageChannel = channel

       consumer, err := client.Subscribe(options)
       if err != nil {
           log.Fatal(err)
       }

       defer consumer.Close()

       // Receive messages from channel. The channel returns a struct which contains message and the consumer from where
       // the message was received. It's not necessary here since we have 1 single consumer, but the channel could be
       // shared across multiple consumers as well
       for cm := range channel {
           msg := cm.Message
           fmt.Printf("Received message  msgId: %v -- content: '%s'\n",
               msg.ID(), string(msg.Payload()))

           consumer.Ack(msg)
       }
     ```

  4. 如何使用消费者接收超时器
     ```
     client, err := NewClient(pulsar.ClientOptions{
         URL: "pulsar://localhost:6650",
     })
     if err != nil {
         log.Fatal(err)
     }
     defer client.Close()

     topic := "test-topic-with-no-messages"
     ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
     defer cancel()

     // create consumer
     consumer, err := client.Subscribe(pulsar.ConsumerOptions{
         Topic:            topic,
         SubscriptionName: "my-sub1",
         Type:             Shared,
     })
     if err != nil {
         log.Fatal(err)
     }
     defer consumer.Close()

     msg, err := consumer.Receive(ctx)
     fmt.Println(msg.Payload())
     if err != nil {
         log.Fatal(err)
     }
     ```

  5. 如何在消费者中使用schema
     ```
     type testJSON struct {
         ID   int    `json:"id"`
         Name string `json:"name"`
     }
     ```
     ```
     var (
         exampleSchemaDef = "{\"type\":\"record\",\"name\":\"Example\",\"namespace\":\"test\"," +
             "\"fields\":[{\"name\":\"ID\",\"type\":\"int\"},{\"name\":\"Name\",\"type\":\"string\"}]}"
     )
     ```
     ```
     client, err := NewClient(pulsar.ClientOptions{
         URL: "pulsar://localhost:6650",
     })
     if err != nil {
         log.Fatal(err)
     }
     defer client.Close()

     var s testJSON

     consumerJS := NewJSONSchema(exampleSchemaDef, nil)
     consumer, err := client.Subscribe(ConsumerOptions{
         Topic:                       "jsonTopic",
         SubscriptionName:            "sub-1",
         Schema:                      consumerJS,
         SubscriptionInitialPosition: SubscriptionPositionEarliest,
     })
     assert.Nil(t, err)
     msg, err := consumer.Receive(context.Background())
     assert.Nil(t, err)
     err = msg.GetSchemaValue(&s)
     if err != nil {
         log.Fatal(err)
     }

     defer consumer.Close()
     ```

- Consumer 配置
  type ConsumerOptions struct {
  	// 主题指定此消费者将订阅的主题。 构造读取器时需要此参数
  	Topic string

  	// 指定此消费者将订阅的主题列表。 
    // 订阅时需要主题、主题列表或主题模式。
  	Topics []string

  	// 指定正则表达式订阅同一命名空间下的多个主题。 
    // 订阅时需要主题、主题列表或主题模式
  	TopicsPattern string

  	// 如果使用TopicsPattern，则指定轮询新分区或新主题的时间间隔。
  	AutoDiscoveryPeriod time.Duration

  	// 指定此使用者的订阅名称。 订阅时需要此参数
  	SubscriptionName string

  	// 将一组应用程序定义的属性附加到使用者 
    // 这些属性将在主题统计信息中可见
  	Properties map[string]string

  	// 选择订阅主题时使用的订阅类型。 
    // 默认为`Exclusive`
  	Type SubscriptionType

  	// 订阅时将设置光标的初始位置
  	// 默认为`Latest`
  	SubscriptionInitialPosition

  	// 死信队列消费者策略的配置。 
    // 例如。 在 N 次尝试处理失败后将消息路由到主题 X 
    // 默认情况下为nil,并且没有 DLQ
  	DLQ *DLQPolicy

  	// 密钥共享消费者策略的配置。
  	KeySharedPolicy *KeySharedPolicy

  	// 自动重试将消息发送到默认填充的 DLQPolicy 主题 
    // 默认为 false
  	RetryEnable bool

    // 为消费者设置一个`MessageChannel` 
    // 当接收到消息时，将其推送到通道进行消费
  	MessageChannel chan ConsumerMessage

    // 设置消费者接收队列的大小。 
    // 消费者接收队列控制在应用程序调用 `Consumer.receive()` 之前，`Consumer` 可以累积多少消息。 使用更高的值可能会增加使用者吞吐量，但会牺牲更大的内存利用率。 
    // 默认值为 `1000` 消息，应该适用于大多数用例。
  	ReceiverQueueSize int

  	// 重新传递未能处理的消息的延迟时间。 默认为 1 分钟。 （参见`Consumer.Nack()`）
  	NackRedeliveryDelay time.Duration

  	// 设置消费者名称。
  	Name string

  	// 如果启用，消费者将从压缩的主题中读取消息，而不是读取主题的完整消息积压。 
    // 这意味着，如果主题已被压缩，则消费者只会看到主题中每个键的最新值，直到主题消息积压中已被压缩的点。 超过该点消息将照常发送。 
    
    // ReadCompacted 只能启用对具有单个活动使用者的持久主题的订阅（即failure或exclusive订阅）。 尝试在订阅非持久主题或共享订阅时启用它，将导致订阅调用抛出 PulsarClientException。
  	ReadCompacted bool

  	// 将订阅标记为已复制以使其跨集群保持同步
  	ReplicateSubscriptionState bool

  	// 一个拦截器链，这些拦截器将在 ConsumerInterceptor 接口中定义的某些点被调用。
  	Interceptors ConsumerInterceptors

    // 通过传递 Schema 的实现来设置自定义模式类型| bytes[]
  	Schema Schema

  	// 设置 reconnectToBroker 的最大重试次数。 (默认: ultimate)
  	MaxReconnectToBroker *uint
  }

### 阅读器
Pulsar Reader处理来自 Pulsar 主题的消息。 读者与消费者不同，因为对于读者，您需要明确指定要从流中的哪条消息开始（另一方面，消费者会自动从最新的未确认消息开始）。 您可以使用 ReaderOptions 对象配置 Go 阅读器。下面是示例：
```
reader, err := client.CreateReader(pulsar.ReaderOptions{
    Topic:          "topic-1",
    StartMessageID: pulsar.EarliestMessageID(),
})
if err != nil {
    log.Fatal(err)
}
defer reader.Close()
```

- Reader operations
  * Topic()	string
    返回读者的话题

  * Next(context.Context)	(Message, error)
    接收关于主题的下一条消息（类似于消费者的 Receive 方法）。 此方法会阻塞，直到消息可用。
  
  * HasNext()	(bool, error)
    检查当前位置是否有可供阅读的消息
 
  * Close()	error
    关闭读取器，禁用其从代理接收消息的能力
  
  * Seek(MessageID)	error
    将与此 reader 关联的订阅重置为特定的消息ID
   
  * SeekByTime(time time.Time)	error
    将与此 reader 关联的订阅重置为特定的消息投递时间

- Reader 示例
  1. 如何使用阅读器读取“下一个”消息
     下面是使用 Next() 方法处理传入消息的 Go 阅读器的示例用法：
     ```
     client, err := pulsar.NewClient(pulsar.ClientOptions{URL: "pulsar://localhost:6650"})
       if err != nil {
           log.Fatal(err)
       }

       defer client.Close()

       reader, err := client.CreateReader(pulsar.ReaderOptions{
           Topic:          "topic-1",
           StartMessageID: pulsar.EarliestMessageID(),
       })
       if err != nil {
           log.Fatal(err)
       }
       defer reader.Close()

       for reader.HasNext() {
           msg, err := reader.Next(context.Background())
           if err != nil {
               log.Fatal(err)
           }

           fmt.Printf("Received message msgId: %#v -- content: '%s'\n",
               msg.ID(), string(msg.Payload()))
       }
     ```

     在上面的示例中，阅读器从最早的可用消息（由 pulsar.EarliestMessage 指定）开始阅读。 阅读器还可以使用 DeserializeMessageID 函数从最新的消息 (pulsar.LatestMessage) 或其他一些由字节指定的消息 ID 开始读取，该函数接受一个字节数组并返回一个 MessageID 对象。
     ```
     lastSavedId := // 从外部存储中读取上次保存的消息 ID 是一个 byte[]

     reader, err := client.CreateReader(pulsar.ReaderOptions{
         Topic:          "my-golang-topic",
         StartMessageID: pulsar.DeserializeMessageID(lastSavedId),
     })
     ```

  2. 如何使用阅读器读取特定消息
     ```
     client, err := NewClient(pulsar.ClientOptions{
          URL: lookupURL,
      })

      if err != nil {
          log.Fatal(err)
      }
      defer client.Close()

      topic := "topic-1"
      ctx := context.Background()

      // create producer
      producer, err := client.CreateProducer(pulsar.ProducerOptions{
          Topic:           topic,
          DisableBatching: true,
      })
      if err != nil {
          log.Fatal(err)
      }
      defer producer.Close()

      // send 10 messages
      msgIDs := [10]MessageID{}
      for i := 0; i < 10; i++ {
          msgID, err := producer.Send(ctx, &pulsar.ProducerMessage{
              Payload: []byte(fmt.Sprintf("hello-%d", i)),
          })
          assert.NoError(t, err)
          assert.NotNil(t, msgID)
          msgIDs[i] = msgID
      }

      // create reader on 5th message (not included)
      reader, err := client.CreateReader(pulsar.ReaderOptions{
          Topic:          topic,
          StartMessageID: msgIDs[4],
      })

      if err != nil {
          log.Fatal(err)
      }
      defer reader.Close()

      // receive the remaining 5 messages
      for i := 5; i < 10; i++ {
          msg, err := reader.Next(context.Background())
          if err != nil {
          log.Fatal(err)
      }

      // create reader on 5th message (included)
      readerInclusive, err := client.CreateReader(pulsar.ReaderOptions{
          Topic:                   topic,
          StartMessageID:          msgIDs[4],
          StartMessageIDInclusive: true,
      })

      if err != nil {
          log.Fatal(err)
      }
      defer readerInclusive.Close()
     ```

- Reader 配置
  type ReaderOptions struct {
    // 指定此消费者将订阅的主题。 构造读取器时需要此参数。
    Topic string

    // 设置阅读器名称。
    Name string

    // 将一组应用程序定义的属性附加到阅读器 ,此属性将在主题统计信息中可见
    Properties map[string]string

    // StartMessageID 初始阅读器定位是通过指定消息 ID 来完成的。 选项是： 
    // * `pulsar.EarliestMessage` ：从主题中可用的最早消息开始读取 
    // * `pulsar.LatestMessage` ：从结束话题开始阅读，只获取阅读器创建后发布的消息
    // * `MessageID` : 从特定的消息 ID 开始读取，阅读器会将自己定位在该特定位置。 要读取的第一条消息将是指定messageID 的next消息
    StartMessageID MessageID

    // 如果为 true，则读取器将从包含的 `StartMessageID` 开始。 默认为`false`，阅读器将从“下一个”消息开始
    StartMessageIDInclusive bool

    // MessageChannel 为消费者设置了一个`MessageChannel` 
    // 当收到消息时，会推送到该通道进行消费
    MessageChannel chan ReaderMessage

    // ReceiverQueueSize 设置消费者接收队列的大小。 
    // 消费者接收队列控制在应用程序调用 Reader.readNext() 之前，Reader 可以累积多少消息。 
    // 使用更高的值可能会增加使用者吞吐量，但会牺牲更大的内存利用率。 
    // 默认值为 {@code 1000} 消息，应该适用于大多数用例。
    ReceiverQueueSize int

    // 设置订阅角色前缀。 默认前缀是“reader”。
    SubscriptionRolePrefix string

    // 如果启用，阅读器将从压缩的主题中读取消息，而不是读取该主题的完整消息积压。 
    // 这意味着，如果主题已被压缩，则读者只会看到主题中每个键的最新值，直到主题消息积压中已被压缩的点。 超过该点消息将照常发送。
    
    // ReadCompacted 只能在从持久主题读取时启用。 尝试在非持久性主题上启用它会导致读者创建调用抛出 PulsarClientException。
    // 默认值 false
    ReadCompacted bool
  }

### 消息
Pulsar Go 客户端提供了一个 ProducerMessage 接口，您可以使用该接口构造 Pulsar 主题上的生产者消息。 这是一个示例消息：
```
msg := pulsar.ProducerMessage{
    Payload: []byte("Here is some message data"),
    Key: "message-key",
    Properties: map[string]string{
        "foo": "bar",
    },
    EventTime: time.Now(),
    ReplicationClusters: []string{"cluster1", "cluster3"},
}

if _, err := producer.send(msg); err != nil {
    log.Fatalf("Could not publish message due to: %v", err)
}
```

- ProducerMessage 配置
  
  type ProducerMessage struct {
  	// 消息的实际数据
  	`Payload` []byte

  	//值和负载是互斥的，Value interface{} for schema message.
  	`Value` interface{}

  	// 与消息关联的可选键（对于主题压缩等特别有用）
  	`Key` string

  	// OrderingKey 设置消息的排序键。
  	`OrderingKey` string

  	// 附加到消息的任何特定于应用程序的元数据的键值映射（键和值都必须是字符串）
  	`Properties` map[string]string

  	// 设置给定消息的事件时间
    // 默认情况下，消息没有关联的事件时间，而是发布时间 
    // 始终存在。 
    // 将事件时间设置为非零时间戳以明确声明事件“发生”的时间，而不是发布消息的时间。
  	`EventTime` time.Time

  	// 覆盖此消息的复制集群。
  	`ReplicationClusters` []string

  	// 禁用此消息的复制
  	`DisableReplication` bool

  	// 设置要分配给当前消息的序列 ID
  	`SequenceID` *int64

    // 请求仅在指定的相对延迟后投递消息。 
    // 注意：只有当消费者通过 `SubscriptionType=Shared` 订阅消费时，消息才会延迟传递。 使用其他订阅类型，消息仍将立即传递。
  	`DeliverAfter` time.Duration

    // 仅在指定的绝对时间戳或之后传递消息。 
    // 注意：只有当消费者通过 `SubscriptionType=Shared` 订阅消费时，消息才会延迟传递。 使用其他订阅类型，消息仍将立即传递。
  	`DeliverAt` time.Time
  }
  
- TLS 加密和身份验证
  
  为了使用 TLS 加密，您需要配置您的客户端来这样做：
  * 使用 pulsar+ssl URL 类型
  * 设置 TLSTrustCertCertsFilePath 到你的客户端和 Pulsar broker 使用的 TLS 证书路径
  * 配置 认证 选项
  
   下面是一个示例：
   ```
   opts := pulsar.ClientOptions{
        URL: "pulsar+ssl://my-cluster.com:6651",
        TLSTrustCertsFilePath: "/path/to/certs/my-cert.csr",
        Authentication: NewAuthenticationTLS("my-cert.pem", "my-key.pem"),
    }
   ```

- OAuth2 身份验证
  要使用 OAuth2 身份验证，您需要配置您的客户端以执行以下操作。 此示例显示如何配置 OAuth2 身份验证。
  ```
  oauth := pulsar.NewAuthenticationOAuth2(map[string]string{
          "type":       "client_credentials",
          "issuerUrl":  "https://dev-kt-aa9ne.us.auth0.com",
          "audience":   "https://dev-kt-aa9ne.us.auth0.com/api/v2/",
          "privateKey": "/path/to/privateKey",
          "clientId":   "0Xx...Yyxeny",
      })
  client, err := pulsar.NewClient(pulsar.ClientOptions{
          URL:              "pulsar://my-cluster:6650",
          Authentication:   oauth,
  })
  ```
  