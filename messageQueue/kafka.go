package messagequeue

import (
	"time"

	"doovvvDP/config"

	"github.com/segmentio/kafka-go"
)

type KafkaService struct {
	OrderProducer *kafka.Writer
	OrderConsumer *kafka.Reader
	KafkaConn     *kafka.Conn
}

func NewKafkaService() *KafkaService {
	return &KafkaService{}
}

func (k *KafkaService) KafkaInit() {
	// k.CreateTopic()
	k.OrderProducer = &kafka.Writer{
		Addr:                   kafka.TCP(config.MyConfig.KafkaConfig.HostPort),
		Topic:                  config.MyConfig.KafkaConfig.OrderTopic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           time.Duration(config.MyConfig.KafkaConfig.Timeout) * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: false,
	}
	k.OrderConsumer = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.MyConfig.KafkaConfig.HostPort},
		Topic:   config.MyConfig.KafkaConfig.OrderTopic,
		// CommitInterval: time.Duration(config.MyConfig.KafkaConfig.Timeout) * time.Second,
		// GroupID:     "order",
		StartOffset: kafka.FirstOffset,
		MaxWait:     10 * time.Second,
	})
}

func (k *KafkaService) Close() {
	k.OrderProducer.Close()
	k.OrderConsumer.Close()
}

func (k *KafkaService) CreateTopic() {
	topic := config.MyConfig.KafkaConfig.OrderTopic

	conn, err := kafka.Dial("tcp", config.MyConfig.KafkaConfig.HostPort)
	if err != nil {
		panic("kafka connect:" + err.Error())
	}
	k.KafkaConn = conn

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     config.MyConfig.KafkaConfig.Partition,
			ReplicationFactor: 1,
		},
	}
	err = k.KafkaConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic("kafka create topic:" + err.Error())
	}
}
