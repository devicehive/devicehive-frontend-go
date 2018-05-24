package dhkafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/devicehive/devicehive-frontend-go/notifications"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

// todo:
//@SerializedName("t")
//private int type;



type KafkaServer struct {
	requestTopic  string
	responseTopic string
	producer      sarama.AsyncProducer
	consumer      sarama.Consumer
	logger        echo.Logger
}

func (k *KafkaServer) Close() {
	if err := k.producer.Close(); err != nil {
		log.Println("Failed to shut down access log producer cleanly", err)
	}
	if err := k.consumer.Close(); err != nil {
		log.Println("Failed to shut down access log consumer cleanly", err)
	}
}

func (k *KafkaServer) Send(msg notifications.MessagesToDevice) {
	kMsg, _ := json.Marshal(KafkaRequestMessage{
		Body:                msg,
		PartionKey:          msg.GetDeviceId(),
		SingleReplyExpected: true,
		ReplyTo:             k.requestTopic,
		CorrelationId:       uuid.NewV4().String(),
	})
	k.producer.Input() <- &sarama.ProducerMessage{
		Value: sarama.ByteEncoder(kMsg),
		Topic: k.requestTopic,
	}
}

func (k *KafkaServer) GetMessages() sarama.PartitionConsumer {
	partitionConsumer, err := k.consumer.ConsumePartition(k.responseTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Failed to connect to RESPONSE_TOPIC")
	}
	return partitionConsumer
}

func New(brokerList []string, logger echo.Logger) *KafkaServer {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama consumer:", err)
	}

	go func() {
		for err := range producer.Errors() {
			logger.Errorf("Failed to write access log entry: %s", err)
		}
	}()

	return &KafkaServer{
		producer:      producer,
		logger:        logger,
		consumer:      consumer,
		requestTopic:  "request_topic",
		responseTopic: "response_topic_" + uuid.NewV1().String(),
	}
}
