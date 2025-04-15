package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"net"
	"strconv"
	"time"
)

//查看消费者偏移量 kafka-consumer-groups.sh --bootstrap-server 43.138.33.205:9092 --describe --group test
//重置消费者偏移量 kafka-consumer-groups.sh --bootstrap-server 43.138.33.205:9092 --group test --reset-offsets --to-earliest --topic op-records --execute

func InitKafka() (*kafka.Writer, *kafka.Reader, *kafka.Reader) {
	conn, err := kafka.Dial("tcp", global.GVA_CONFIG.Kafka.Brokers[0])
	if err != nil {
		global.GVA_LOG.Info("connect to kafka broker error:", zap.Error(err))
	}
	defer conn.Close()
	initializeTopic(conn)

	return initWriter(), initMysqlReader(), initEsReader()
}

// 初始化创建主题
func initializeTopic(conn *kafka.Conn) {
	//获取controller
	controller, err := conn.Controller()
	if err != nil {
		global.GVA_LOG.Info("get kafka controller error:", zap.Error(err))
	}
	//创建controller连接
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		global.GVA_LOG.Info("connect to kafka controller error:", zap.Error(err))
	}
	defer controllerConn.Close()
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             global.GVA_CONFIG.Kafka.Topic,
			NumPartitions:     3, //创建三个分区
			ReplicationFactor: 2, //每个分区有两个副本
		},
	}
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		global.GVA_LOG.Info("create topic error:", zap.Error(err))
	}
	global.GVA_LOG.Info("create topic success")
}

// 初始化一个生产者实例
func initWriter() *kafka.Writer {
	kw := &kafka.Writer{
		Addr:                   kafka.TCP(global.GVA_CONFIG.Kafka.Brokers...),
		Topic:                  global.GVA_CONFIG.Kafka.Topic,
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireAll,
		Async:                  true,
		AllowAutoTopicCreation: false, //使用自动创建会在第一次创建时进行领导层选举，
	}
	return kw
}

// InitMysqlReader 初始化一个MySQL消费者实例
func initMysqlReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        global.GVA_CONFIG.Kafka.Brokers,
		GroupID:        global.GVA_CONFIG.Kafka.MysqlGroupId,
		Topic:          global.GVA_CONFIG.Kafka.Topic,
		SessionTimeout: 5 * time.Minute,
		CommitInterval: 0,    //关闭自动提交
		MinBytes:       10e3, //10K
		MaxBytes:       10e6,
	})
}

func initEsReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        global.GVA_CONFIG.Kafka.Brokers,
		GroupID:        global.GVA_CONFIG.Kafka.EsGroupId,
		Topic:          global.GVA_CONFIG.Kafka.Topic,
		SessionTimeout: 5 * time.Minute,
		CommitInterval: 0, //关闭自动提交
		MinBytes:       10e3,
		MaxBytes:       10e6,
	})
}
