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
//correlationId rpc模式下，每个请求的唯一值,其它模式不用传,rpc客户端传correlationId第2个参数(队列名)
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

//rpc模式发送者
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

//rpc模式接收者
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
