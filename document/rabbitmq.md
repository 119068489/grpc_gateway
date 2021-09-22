# RabbitMQ

## 介绍
消息总线(Message Queue)，是一种跨进程、异步的通信机制，用于上下游传递消息。由消息系统来确保消息的可靠传递。

### 适应场景
* 上下游逻辑解耦&&物理解耦
* 保证数据最终一致性
* 广播
* 错峰流控等等

### 特点
RabbitMQ是由Erlang语言开发的AMQP的开源实现。

AMQP：Advanced Message Queue，高级消息队列协议。它是应用层协议的一个开放标准，为面向消息的中间件设计，基于此协议的客户端与消息中间件可传递消息，并不受产品、开发语言灯条件的限制。

* 可靠性(Reliablity)：使用了一些机制来保证可靠性，比如持久化、传输确认、发布确认。
* 灵活的路由(Flexible Routing)：
  在消息进入队列之前，通过Exchange来路由消息。对于典型的路由功能，Rabbit已经提供了一些内置的Exchange来实现。针对更复杂的路由功能，可以将多个Exchange绑定在一起，也通过插件机制实现自己的Exchange。
* 消息集群(Clustering)：多个RabbitMQ服务器可以组成一个集群，形成一个逻辑Broker。
* 高可用(Highly Avaliable Queues)：队列可以在集群中的机器上进行镜像，使得在部分节点出问题的情况下队列仍然可用。
* 多种协议(Multi-protocol)：支持多种消息队列协议，如STOMP、MQTT等。
* 多种语言客户端(Many Clients)：几乎支持所有常用语言，比如Java、.NET、Ruby等。
* 管理界面(Management UI)：提供了易用的用户界面，使得用户可以监控和管理消息Broker的许多方面。
* 跟踪机制(Tracing)：如果消息异常，RabbitMQ提供了消息的跟踪机制，使用者可以找出发生了什么。
* 插件机制(Plugin System)：提供了许多插件，来从多方面进行扩展，也可以编辑自己的插件。

### 基本概念
* Broker：标识消息队列服务器实体.

* Virtual Host：
  虚拟主机。标识一批交换机、消息队列和相关对象。虚拟主机是共享相同的身份认证和加密环境的独立服务器域。每个vhost本质上就是一个mini版的RabbitMQ服务器，拥有自己的队列、交换器、绑定和权限机制。vhost是AMQP概念的基础，必须在链接时指定，RabbitMQ默认的vhost是 /。

* Exchange：交换器，用来接收生产者发送的消息并将这些消息路由给服务器中的队列。

* Queue：
  消息队列，用来保存消息直到发送给消费者。它是消息的容器，也是消息的终点。一个消息可投入一个或多个队列。消息一直在队列里面，等待消费者连接到这个队列将其取走。

* Banding：
  绑定，用于消息队列和交换机之间的关联。一个绑定就是基于路由键将交换机和消息队列连接起来的路由规则，所以可以将交换器理解成一个由绑定构成的路由表。

* Channel：
  信道，多路复用连接中的一条独立的双向数据流通道。信道是建立在真实的TCP连接内的虚拟链接，AMQP命令都是通过信道发出去的，不管是发布消息、订阅队列还是接收消息，这些动作都是通过信道完成。因为对于操作系统来说，建立和销毁TCP都是非常昂贵的开销，所以引入了信道的概念，以复用一条TCP连接。

* Connection：网络连接，比如一个TCP连接。

* Publisher：消息的生产者，也是一个向交换器发布消息的客户端应用程序。

* Consumer：消息的消费者，表示一个从一个消息队列中取得消息的客户端应用程序。

* Message：
  消息，消息是不具名的，它是由消息头和消息体组成。消息体是不透明的，而消息头则是由一系列的可选属性组成，这些属性包括routing-key(路由键)、priority(优先级)、delivery-mode(消息可能需要持久性存储[消息的路由模式])等。

### 六种工作模式
1. simple简单模式
   
   ![图例](https://topgoer.com/static/9.3/0.png)

   * 消息产生着(Publisher)将消息放入队列
   * 消息的消费者(consumer)监听(while)消息队列,如果队列中有消息,就消费掉,消息被拿走后,自动从队列中删除(隐患 消息可能没有被消费者正确处理,已经从队列中消失了,造成消息的丢失)应用场景:聊天(中间有一个过度的服务器;p端,c端)

   [示例代码](#simple模式-简单模式)

2. work工作模式(资源的竞争)

   ![图例](https://topgoer.com/static/9.3/2.png)

   * 消息产生者将消息放入队列消费者可以有多个,消费者1,消费者2,同时监听同一个队列,消息被消费? C1 C2共同争抢当前的消息队列内容,谁先拿到谁负责消费消息(隐患,高并发情况下,默认会产生某一个消息被多个消费者共同使用,可以设置一个开关(syncronize,与同步锁的性能不一样) 保证一条消息只能被一个消费者使用)
   * 应用场景:红包;大项目中的资源调度(任务分配系统不需知道哪一个任务执行系统在空闲,直接将任务扔到消息队列中,空闲的系统自动争抢)

   [示例代码](#work模式-工作模式)
  
3. publish/subscribe发布订阅(共享资源)

   ![图例](https://topgoer.com/static/9.3/3.png)

   * X代表交换机rabbitMQ内部组件,erlang 消息产生者是代码完成,代码的执行效率不高,消息产生者将消息放入交换机,交换机发布订阅把消息发送到所有消息队列中,对应消息队列的消费者拿到消息进行消费
   * 相关场景:邮件群发,群聊天,广播(广告)

   [示例代码](#publish模式-订阅模式)
  
4. routing路由模式

   ![图例](https://topgoer.com/static/9.3/4.png)

   * 消息生产者将消息发送给交换机按照路由判断,路由是字符串(info) 当前产生的消息携带路由字符(对象的方法),交换机根据路由的key,只能匹配上路由key对应的消息队列,对应的消费者才能消费消息;
   * 根据业务功能定义路由字符串
   * 从系统的代码逻辑中获取对应的功能字符串,将消息任务扔到对应的队列中业务场景:error 通知;EXCEPTION;错误通知的功能;传统意义的错误通知;客户通知;利用key路由,可以将程序中的错误封装成消息传入到消息队列中,开发者可以自定义消费者,实时接收错误;

   [示例代码](#routing模式-路由模式)
  
5. topic 主题模式(路由模式的一种)

   ![图例](https://topgoer.com/static/9.3/5.png)

   * 星号井号代表通配符
   * 星号代表多个单词,井号代表一个单词
   * 路由功能添加模糊匹配
   * 消息产生者产生消息,把消息交给交换机
   * 交换机根据key的规则模糊匹配到对应的队列,由队列的监听消费者接收消息消费

    [示例代码](#topic模式-话题模式)

6. rpc模式
   
   ![图例](https://s1.ax1x.com/2020/09/02/dz0WU1.png)

   * 对于RPC请求，客户端发送一条带有两个属性的消息:replyTo,设置为仅为请求创建的匿名独占队列,和correlationId,设置为每个请求的惟一id值。
   * 请求被发送到rpc_queue队列。
   * RPC工作进程(即:服务器)在队列上等待请求。当一个请求出现时，它执行任务,并使用replyTo字段中的队列将结果发回客户机。
   * 客户机在回应消息队列上等待数据。当消息出现时，它检查correlationId属性。如果匹配请求中的值，则向程序返回该响应数据。

    [示例代码](#rpc模式)

## windows下安装
1. 下载并安装erlang
   原因：RabbitMQ服务端代码是使用并发式语言Erlang编写的，安装Rabbit MQ的前提是安装Erlang。

   [下载地址](http://www.erlang.org/downloads)

   安装Erlang后将bin目录加入到path中

2. 下载并安装RabbitMQ
   [下载地址](http://www.rabbitmq.com/download.html)

   RabbitMQ安装好后接下来安装RabbitMQ-Plugins。

   安装目录 E:\RabbitMQ\rabbitmq_server-3.8.1\sbin目录下打开cmd

   输入`rabbitmq-plugins enable rabbitmq_management`命令进行安装

   如果出现下面的提示表示运行成功

   ![图例](https://topgoer.com/static/9.3/23.png)

   启动服务：`rabbitmq-server.bat`

   如果出现下面的提示表示启动成功

   ![图例](https://topgoer.com/static/9.3/24.png)

   如果出现以下错误

   ![图例](https://topgoer.com/static/9.3/27.png)

   这个是因为rabbit已经启动了，不能再次启动,解决办法：
   
   在服务中找到RabbitMQ停止即可再次运行启动命令。其实也可以直接下一步进入管理系统

   rabbitmq启动成功，浏览器中进入管理登录http://localhost:15672

   输入guest,guest进入rabbitMQ管理控制台

3. 重复安装Rabbit Server的坑
   如果不是第一次在Windows上安装Rabbit Server一定要把Rabbit和Erlang卸载干净之后，找到注册表：HKEY_LOCAL_MACHINE\SOFTWARE\Ericsson\Erlang\ErlSrv 删除其下的所有项。

   不然会出现Rabbit安装之后启动不了的情况，理论上卸载的顺序也是先Rabbit再Erlang。

## 代码实现

### 准备工作
首先我们新建一个Virtual Host并且给他分配一个用户名，用来隔离数据，根据自己需要自行创建,也可以不创建使用系统帐户

如果不创建连接MQ时使用：amqp://guest:guest@localhost:5672/
如果创建了连接MQ时使用：amqp://testuser:123456@127.0.0.1:5672/testhost
(testuser:123456 是我创建的用户名和密码。testhost是我创建的Virtual Host。)

- 第一步：创建Virtual Host
  
  ![图例](https://topgoer.com/static/9.3/30.png)

- 第二步：创建user
  
  ![图例](https://topgoer.com/static/9.3/31.png)

- 第三步：点击新创建的user
  
  ![图例](https://topgoer.com/static/9.3/32.png)

- 第四步：设置权限
  
  ![图例](https://topgoer.com/static/9.3/33.png)

- 第五步：查询user权限
  
  ![图例](https://topgoer.com/static/9.3/34.png)





封装一个RabbitMQ
rabitmq.go //这个是RabbitMQ的封装
```
package easygo

import (
	"fmt"
	"log"

	"github.com/astaxie/beego/logs"
	"github.com/streadway/amqp"
)

//连接信息amqp://kuteng:kuteng@127.0.0.1:5672/kuteng这个信息是固定不变的amqp://事固定参数后面两个是用户名密码ip地址端口号Virtual Host
// const MQURL = "amqp://testuser:123456@127.0.0.1:5672/testhost"

type DeclareReq struct {
	Queue      string     //队列名称
	Exchange   string     //交换机名称
	Kind       string     //
	Passive    bool       //
	Durable    bool       //是否持久化
	Exclusive  bool       //是否具有排他性
	AutoDelete bool       //是否自动删除
	NoWait     bool       //是否阻塞处理
	Internal   bool       //true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
	Arguments  amqp.Table //额外的属性
}

type BasicConsume struct {
	Queue       string     //队列名称
	ConsumerTag string     //用来区分多个消费者
	NoLocal     bool       //设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
	NoAck       bool       //是否自动应答
	Exclusive   bool       //是否独有
	NoWait      bool       //是否阻塞处理
	Arguments   amqp.Table //额外的属性
}

//rabbitMQ结构体
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机名称
	Exchange string
	//bind Key 名称
	Key string
	//连接信息
	Mqurl string
}

func NewRabbitMQ() *RabbitMQ { // services map[string]interface{},
	p := &RabbitMQ{}
	p.Init()
	return p
}

//初始化
func (r *RabbitMQ) Init() {
	// r.Mqurl = MQURL
}

//断开channel 和 connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

//创建简单模式下RabbitMQ实例
func (r *RabbitMQ) NewClient() {
	var err error
	//获取connection
	r.conn, err = amqp.Dial(r.Mqurl)
	r.failOnErr(err, "failed to connect rabb"+
		"itmq!")
	//获取channel
	r.channel, err = r.conn.Channel()
	r.failOnErr(err, "failed to open a channel")
}

//申请队列
func (r *RabbitMQ) QueueDeclare(req *DeclareReq) amqp.Queue {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		req.Queue,
		//是否持久化
		req.Durable,
		//是否自动删除
		req.AutoDelete,
		//是否具有排他性
		req.Exclusive,
		//是否阻塞处理
		req.NoWait,
		//额外的属性
		amqp.Table(req.Arguments),
	)

	r.failOnErr(err, "Failed to declare a queue")

	return q
}

//创建交换机
func (r *RabbitMQ) ExchangeDeclare(req *DeclareReq) {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		req.Kind,
		req.Durable,
		req.AutoDelete,
		req.Internal,
		req.NoWait,
		amqp.Table(req.Arguments),
	)

	r.failOnErr(err, "Failed to declare an exchange")
}

//绑定队列到交换机中
func (r *RabbitMQ) QueueBind(queueName string, nowait bool, args amqp.Table) {
	//绑定队列到 exchange 中
	err := r.channel.QueueBind(
		queueName,
		//在订阅模式下，这里的key要为空
		//路由模式下，需要绑定key
		r.Key,
		r.Exchange,
		nowait,
		args)

	r.failOnErr(err, "Failed to Bind the queue to the exchange")
}

//调用channel 发送消息到队列中
//
//mandatory 如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
//
//immediate 如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
//
//correlationId rpc模式下，每个请求的唯一值,其它模式不用传
func (r *RabbitMQ) Publish(message, key string, mandatory, immediate bool, correlationId ...string) {
	pmsg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	}

	if len(correlationId) > 0 {
		pmsg.CorrelationId = correlationId[0]
	}

	if len(correlationId) > 1 {
		pmsg.ReplyTo = correlationId[1]
	}

	err := r.channel.Publish(
		r.Exchange,
		key,
		mandatory,
		immediate,
		pmsg)

	if err != nil {
		r.failOnErr(err, "Failed to send message to the queue")
	}
}

//接收消息返回消息列表
func (r *RabbitMQ) RecieveReMsgs(b *BasicConsume) <-chan amqp.Delivery {
	msgs, err := r.channel.Consume(
		b.Queue, // queue
		//用来区分多个消费者
		b.ConsumerTag, // consumer
		//是否自动应答
		b.NoAck, // auto-ack
		//是否独有
		b.Exclusive, // exclusive
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		b.NoLocal, // no-local
		//列是否阻塞
		b.NoWait,                // no-wait
		amqp.Table(b.Arguments), // args
	)
	r.failOnErr(err, "Failed to register a consumer")

	return msgs
}

//接收消息并消耗
func (r *RabbitMQ) Consume(b *BasicConsume, f ...func(amqp.Delivery)) {
	msgs := r.RecieveReMsgs(b)

	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for d := range msgs {
			//消息逻辑处理，可以自行设计逻辑
			// log.Printf("Received a message: %s", d.Body)
			if len(f) > 0 {
				f[0](d) //传进来的消息处理函数
			}

		}
	}()

	logs.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

//simple 模式生产者
func (r *RabbitMQ) PublishSimple(message string) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	r.QueueDeclare(&DeclareReq{
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
	})
	//调用channel 发送消息到队列中
	r.Publish(message, r.QueueName, false, false)
}

//simple 模式下消费者
func (r *RabbitMQ) ConsumeSimple(f func(amqp.Delivery)) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q := r.QueueDeclare(&DeclareReq{
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
	})
	//接收消息
	r.Consume(&BasicConsume{
		Queue:       q.Name,
		ConsumerTag: "",
		NoAck:       true,
		Exclusive:   false,
		NoLocal:     false,
		NoWait:      false,
	}, f)
}

//订阅模式生产者
func (r *RabbitMQ) PublishPub(message string) {
	//1.尝试创建交换机
	r.ExchangeDeclare(&DeclareReq{
		Exchange:   r.Exchange,
		Kind:       "fanout",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
	})

	//2.发送消息
	r.Publish(message, "", false, false)
}

//订阅模式消费者
func (r *RabbitMQ) RecieveSub(f func(amqp.Delivery)) {
	//1.试探性创建交换机
	r.ExchangeDeclare(&DeclareReq{
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
	})
	//2.试探性创建队列，这里注意队列名称不要写
	q := r.QueueDeclare(&DeclareReq{
		Durable:    false,
		AutoDelete: false,
		Exclusive:  true,
		NoWait:     false,
	})

	//绑定队列到 exchange 中
	r.QueueBind(q.Name, false, nil)

	//消费消息
	r.Consume(&BasicConsume{
		Queue:       q.Name,
		ConsumerTag: "",
		NoAck:       true,
		Exclusive:   false,
		NoLocal:     false,
		NoWait:      false,
	}, f)
}

//路由模式发送消息
func (r *RabbitMQ) PublishRouting(message string) {
	//1.尝试创建交换机
	r.ExchangeDeclare(&DeclareReq{
		Exchange:   r.Exchange,
		Kind:       "direct",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
	})

	//2.发送消息
	r.Publish(message, r.Key, false, false)
}

//路由模式接受消息
func (r *RabbitMQ) RecieveRouting(f func(amqp.Delivery)) {
	//1.试探性创建交换机
	r.ExchangeDeclare(&DeclareReq{
		Exchange:   r.Exchange,
		Kind:       "direct",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
	})

	//2.试探性创建队列，这里注意队列名称不要写
	q := r.QueueDeclare(&DeclareReq{
		Durable:    false,
		AutoDelete: false,
		Exclusive:  true,
		NoWait:     false,
	})

	//绑定队列到 exchange 中
	r.QueueBind(q.Name, false, nil)

	//消费消息
	r.Consume(&BasicConsume{
		Queue:       q.Name,
		ConsumerTag: "",
		NoAck:       true,
		Exclusive:   false,
		NoLocal:     false,
		NoWait:      false,
	}, f)
}

//话题模式发送消息
func (r *RabbitMQ) PublishTopic(message string) {
	//1.尝试创建交换机
	r.ExchangeDeclare(&DeclareReq{
		Exchange:   r.Exchange,
		Kind:       "topic",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
	})
	//2.发送消息
	r.Publish(message, r.Key, false, false)
}

//话题模式接受消息
//要注意key,规则
//其中“*”用于匹配一个单词，“#”用于匹配多个单词（可以是零个）
//匹配 kuteng.* 表示匹配 kuteng.hello, kuteng.hello.one需要用kuteng.#才能匹配到
func (r *RabbitMQ) RecieveTopic(f func(amqp.Delivery)) {
	//1.试探性创建交换机
	r.ExchangeDeclare(&DeclareReq{
		Exchange:   r.Exchange,
		Kind:       "topic",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
	})
	//2.试探性创建队列，这里注意队列名称不要写
	q := r.QueueDeclare(&DeclareReq{
		Durable:    false,
		AutoDelete: false,
		Exclusive:  true,
		NoWait:     false,
	})

	//绑定队列到 exchange 中
	r.QueueBind(q.Name, false, nil)

	//消费消息
	r.Consume(&BasicConsume{
		Queue:       q.Name,
		ConsumerTag: "",
		NoAck:       true,
		Exclusive:   false,
		NoLocal:     false,
		NoWait:      false,
	}, f)
}

func (r *RabbitMQ) PublishRpc(message string) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q := r.QueueDeclare(&DeclareReq{
		Queue:      r.QueueName,
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
	})

	// 我们可能要运行多个服务器进程。为了将负载平均分配给多个服务器，我们需要在通道上设置prefetch设置。
	r.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	msgs := r.RecieveReMsgs(&BasicConsume{
		Queue:     q.Name,
		NoAck:     false,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
	})

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			response := message + string(d.Body)

			r.Publish(response, d.ReplyTo, false, false, d.CorrelationId)
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}

func (r *RabbitMQ) RecieveRpc(message string) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q := r.QueueDeclare(&DeclareReq{
		Queue:      r.QueueName,
		Durable:    false,
		AutoDelete: false,
		Exclusive:  true,
		NoWait:     false,
	})

	msgs := r.RecieveReMsgs(&BasicConsume{
		Queue:       q.Name,
		ConsumerTag: "",
		NoAck:       true,
		Exclusive:   false,
		NoLocal:     false,
		NoWait:      false,
	})

	corrId := RandomString(32)
	r.Publish(message, r.Key, false, false, corrId, q.Name)

	for d := range msgs {
		if corrId == d.CorrelationId {
			logs.Info(" [.] Got %s", d.Body)
			break
		}
	}
}
```


### Simple模式 简单模式

Simple模式下用QueueName创建RabbitMQ实例

- 生产者
  mainSimlpePublish.go //Publish 先启动
  ```
  package main

  import (
  	"grpc_gateway/easygo" //这个是RabbitMQ的封装rabitmq.go的包名

  	"github.com/astaxie/beego/logs"
  )

  var Rabbitmq *easygo.RabbitMQ

  const MQURL = "amqp://testuser:123456@127.0.0.1:5672/testhost"

  func init() {
  	Rabbitmq = easygo.NewRabbitMQ()
  	Rabbitmq.Mqurl = MQURL //可以用常量，也可以根据需要用变量
  }

  func main() {
  	Rabbitmq.QueueName = "testhost" //随便写，生产消费一致就行
  	Rabbitmq.NewClient()
  	Rabbitmq.PublishSimple("Hello testuser111!")
  	logs.Info("发送成功！")
  }

  ```


- 消费者
  mainSimpleRecieve.go
  ```
  package main

  import (
  	"grpc_gateway/easygo" //这个是RabbitMQ的封装rabitmq.go的包名
  )

  var Rabbitmq *easygo.RabbitMQ

  const MQURL = "amqp://testuser:123456@127.0.0.1:5672/testhost"

  func init() {
  	Rabbitmq = easygo.NewRabbitMQ()
  	Rabbitmq.Mqurl = MQURL //可以用常量，也可以根据需要用变量
  }

  func main() {
  	Rabbitmq.QueueName = "testhost" //随便写，生产消费一致就行
  	Rabbitmq.NewClient()
  	Rabbitmq.ConsumeSimple(ReadMsg)
  }
  
  //消费者消费消息的函数
  func ReadMsg(d amqp.Delivery) {
  	logs.Info("Received a message: %s", d.Body)
  }
  ```

### Work模式 工作模式
工作模式，一个消息只能被一个消费者获取

Work模式和Simple模式相比代码并没有发生变化只是多了一个消费者

修改以下生产者代码，多发送几条
```
package main

import (
	"grpc_gateway/easygo"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
)

var Rabbitmq *easygo.RabbitMQ

const MQURL = "amqp://testuser:123456@127.0.0.1:5672/testhost"

func init() {
	Rabbitmq = easygo.NewRabbitMQ()
	Rabbitmq.Mqurl = MQURL
}

func main() {
	Rabbitmq.QueueName = "testhost"
	Rabbitmq.NewClient()

	for i := 0; i <= 100; i++ {
		Rabbitmq.PublishSimple("Hello testuser!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		logs.Info("[%d] 发送成功！", i)
	}

}
```

启动2个消费者即可

结果是轮询接收，每个消费者轮流接收消息

### Publish模式 订阅模式
订阅模式，消息被路由投递给多个队列，一个消息被多个消费者获取

Publish模式下用Exchange创建RabbitMQ实例

订阅模式下所有消费者收到同样的消息


- 生产者
  mainPub.go 
  ```
  package main

  import (
  	"grpc_gateway/easygo" //这个是RabbitMQ的封装rabitmq.go的包名

  	"github.com/astaxie/beego/logs"
  )

  var Rabbitmq *easygo.RabbitMQ

  const MQURL = "amqp://testuser:123456@127.0.0.1:5672/testhost"

  func init() {
  	Rabbitmq = easygo.NewRabbitMQ()
  	Rabbitmq.Mqurl = MQURL //可以用常量，也可以根据需要用变量
  }

  func main() {
  	//随便写，生产消费一致就行
  	Rabbitmq.Exchange = "newProduct"
	Rabbitmq.NewClient()

	for i := 1; i <= 10; i++ {
		Rabbitmq.PublishPub("Hello Publish mode!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		logs.Info("[%d] 发送成功！", i)
	}
  }

  ```


- 消费者
  mainSub.go (两个消费者代码是一样的)
  ```
  package main

  import (
  	"grpc_gateway/easygo" //这个是RabbitMQ的封装rabitmq.go的包名
  )

  var Rabbitmq *easygo.RabbitMQ

  const MQURL = "amqp://testuser:123456@127.0.0.1:5672/testhost"

  func init() {
  	Rabbitmq = easygo.NewRabbitMQ()
  	Rabbitmq.Mqurl = MQURL //可以用常量，也可以根据需要用变量
  }

  func main() {
  	Rabbitmq.Exchange = "newProduct" //随便写，生产消费一致就行
	Rabbitmq.NewClient()
	Rabbitmq.RecieveSub(ReadMsg)
  }

  //消费者消费消息的函数
  func ReadMsg(d amqp.Delivery) {
  	logs.Info("Received a message: %s", d.Body)
  }
  ```

### Routing模式 路由模式
路由模式，一个消息被多个消费者获取，并且消息的目标队列可被生产者指定

Routing模式用Exchange,Key创建RabbitMQ实例

路由模式下不同的key消费者收到不一样的消息

- 生产者
  mainpublish.go
  ```
  package main

  import (
  	"grpc_gateway/easygo"
  	"strconv"
  	"time"

  	"github.com/astaxie/beego/logs"
  )

  var Rabbitmq *easygo.RabbitMQ
  var Rabbitmq2 *easygo.RabbitMQ

  const MQURL = "amqp://testuser:123456@127.0.0.1:5672/testhost"

  func init() {
  	Rabbitmq = easygo.NewRabbitMQ()
  	Rabbitmq.Mqurl = MQURL

  	Rabbitmq2 = easygo.NewRabbitMQ()
  	Rabbitmq2.Mqurl = MQURL
  }

  func main() {
  	publishRouting()
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
  ```

- 消费者
  mainrecieve1.go (两个消费者代码是不一样的)
  ```
  package main

  import (
  	"grpc_gateway/easygo" //这个是RabbitMQ的封装rabitmq.go的包名
  )

  var Rabbitmq *easygo.RabbitMQ

  const MQURL = "amqp://testuser:123456@127.0.0.1:5672/testhost"

  func init() {
  	Rabbitmq = easygo.NewRabbitMQ()
  	Rabbitmq.Mqurl = MQURL //可以用常量，也可以根据需要用变量
  }

  func main() {
    Rabbitmq.Exchange = "kuteng"
	Rabbitmq.Key = "kuteng_one"
	Rabbitmq.NewClient()
	Rabbitmq.RecieveRouting(ReadMsg)
  }

  //消费者消费消息的函数
  func ReadMsg(d amqp.Delivery) {
  	logs.Info("Received a message: %s", d.Body)
  }
  ```
  mainrecieve2.go (两个消费者代码是不一样的)
  ```
  package main

  import (
  	"grpc_gateway/easygo" //这个是RabbitMQ的封装rabitmq.go的包名
  )

  var Rabbitmq *easygo.RabbitMQ

  const MQURL = "amqp://testuser:123456@127.0.0.1:5672/testhost"

  func init() {
  	Rabbitmq = easygo.NewRabbitMQ()
  	Rabbitmq.Mqurl = MQURL //可以用常量，也可以根据需要用变量
  }

  func main() {
    Rabbitmq.Exchange = "kuteng"
	Rabbitmq.Key = "kuteng_two"
	Rabbitmq.NewClient()
	Rabbitmq.RecieveRouting(ReadMsg)
  }

  //消费者消费消息的函数
  func ReadMsg(d amqp.Delivery) {
  	logs.Info("Received a message: %s", d.Body)
  }
  ```

### Topic模式 话题模式
一个消息被多个消费者获取，消息的目标queue可用BindingKey以通配符，（#：一个或多个词，*：一个词）的方式指定

Topic模式用Exchange,Key创建RabbitMQ实例

由于匹配规则不一样，消费者1接收2条消息，消费者2接收1条消息

- 生产者
  mainpublish.go
  ```
  package main

  import (
  	"grpc_gateway/easygo"
  	"strconv"
  	"time"

  	"github.com/astaxie/beego/logs"
  )

  var Rabbitmq *easygo.RabbitMQ
  var Rabbitmq2 *easygo.RabbitMQ

  const MQURL = "amqp://testuser:123456@127.0.0.1:5672/testhost"

  func init() {
  	Rabbitmq = easygo.NewRabbitMQ()
  	Rabbitmq.Mqurl = MQURL

  	Rabbitmq2 = easygo.NewRabbitMQ()
  	Rabbitmq2.Mqurl = MQURL
  }

  func main() {
  	publishTopic()
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
  ```

- 消费者
  mainrecieve1.go (两个消费者代码是不一样的)
  ```
  package main

  import (
  	"grpc_gateway/easygo" //这个是RabbitMQ的封装rabitmq.go的包名
  )

  var Rabbitmq *easygo.RabbitMQ

  const MQURL = "amqp://testuser:123456@127.0.0.1:5672/testhost"

  func init() {
  	Rabbitmq = easygo.NewRabbitMQ()
  	Rabbitmq.Mqurl = MQURL //可以用常量，也可以根据需要用变量
  }

  func main() {
    Rabbitmq.Exchange = "exKutengTopic"
  	Rabbitmq.Key = "#"
  	Rabbitmq.NewClient()
  	Rabbitmq.RecieveTopic(ReadMsg)
  }

  //消费者消费消息的函数
  func ReadMsg(d amqp.Delivery) {
  	logs.Info("Received a message: %s", d.Body)
  }
  ```
  mainrecieve2.go (两个消费者代码是不一样的)
  ```
  package main

  import (
  	"grpc_gateway/easygo" //这个是RabbitMQ的封装rabitmq.go的包名
  )

  var Rabbitmq *easygo.RabbitMQ

  const MQURL = "amqp://testuser:123456@127.0.0.1:5672/testhost"

  func init() {
  	Rabbitmq = easygo.NewRabbitMQ()
  	Rabbitmq.Mqurl = MQURL //可以用常量，也可以根据需要用变量
  }

  func main() {
    Rabbitmq.Exchange = "exKutengTopic"
  	Rabbitmq.Key = "kuteng.*.two"
  	Rabbitmq.NewClient()
  	Rabbitmq.RecieveTopic(ReadMsg)
  }

  //消费者消费消息的函数
  func ReadMsg(d amqp.Delivery) {
  	logs.Info("Received a message: %s", d.Body)
  }
  ```

### Rpc模式

Rpc模式 生产者用QueueName创建RabbitMQ实例，消费者用Key创建RabbitMQ实例(创建一个匿名排他回调队列)

- 生产者
  mainpublish.go
  ```
  package main

  import (
  	"grpc_gateway/easygo"
  	"strconv"
  	"time"

  	"github.com/astaxie/beego/logs"
  )

  var Rabbitmq *easygo.RabbitMQ
  var Rabbitmq2 *easygo.RabbitMQ

  const MQURL = "amqp://testuser:123456@127.0.0.1:5672/testhost"

  func init() {
  	Rabbitmq = easygo.NewRabbitMQ()
  	Rabbitmq.Mqurl = MQURL
  }

  func main() {
  	publishRpc()
  }

  //Rpc模式
  func publishRpc() {
  	Rabbitmq.QueueName = "rpc_queue"
  	Rabbitmq.NewClient()

  	Rabbitmq.PublishRpc("Hello rpcuser!")
  }
  ```

- 消费者
  mainrecieve1.go (两个消费者代码是不一样的)
  ```
  package main

  import (
  	"grpc_gateway/easygo" //这个是RabbitMQ的封装rabitmq.go的包名
  )

  var Rabbitmq *easygo.RabbitMQ

  const MQURL = "amqp://testuser:123456@127.0.0.1:5672/testhost"

  func init() {
  	Rabbitmq = easygo.NewRabbitMQ()
  	Rabbitmq.Mqurl = MQURL //可以用常量，也可以根据需要用变量
  }

  func main() {
    Rabbitmq.Key = "rpc_queue"
  	Rabbitmq.NewClient()

  	Rabbitmq.RecieveRpc("ok")
  }

  ```

