package writer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"forum/log"
	"forum/util"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/spf13/viper"
)

const (
	writeTimeout = 5 * time.Second
	maxAttempts  = 3
)

type KafkaWriter struct {
	Self *kafka.Writer
}

func NewKafkaWriter(topic string) *KafkaWriter {
	mechanism := plain.Mechanism{
		Username: viper.GetString("kafka.username"),
		Password: viper.GetString("kafka.password"),
	}

	transport := &kafka.Transport{
		SASL: mechanism,
		TLS:  nil,
	}

	return &KafkaWriter{
		Self: &kafka.Writer{
			Addr:         kafka.TCP(viper.GetString("kafka.addr")),
			Topic:        topic,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireNone,
			MaxAttempts:  maxAttempts,
			BatchTimeout: 10 * time.Millisecond, // 默认值是 1s，调整致 10ms
			Transport:    transport,
		},
	}
}

func (w *KafkaWriter) Close() error {
	return w.Self.Close()
}

func (w *KafkaWriter) PublishMessage(key string, value string, header ...kafka.Header) error {
	msg := kafka.Message{
		Key:     []byte(key),
		Value:   []byte(value),
		Headers: header,
	}

	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	return w.Self.WriteMessages(ctx, msg)
}

// -------------- 死信队列初始化 -----------------

var (
	deadLetterOnce   sync.Once
	DeadLetterWriter *KafkaWriter
)

const deadLetterTopic = "forum_dead_letter"

type deadLetterMessage struct {
	OriginalValue string `json:"original_value"`
	OriginalKey   string `json:"original_key"`
}

// DefaultEnterDeadLetter 将消息写入死信队列
func DefaultEnterDeadLetter(msg kafka.Message, err error) {
	deadLetterOnce.Do(func() {
		DeadLetterWriter = NewKafkaWriter(deadLetterTopic)
	})

	dlm := deadLetterMessage{
		OriginalKey:   string(msg.Key),
		OriginalValue: string(msg.Value),
	}

	header := []kafka.Header{
		{Key: "OriginalTopic", Value: []byte(msg.Topic)},
		{Key: "Error", Value: []byte(err.Error())},
		{Key: "Time", Value: []byte(util.GetCurrentTime())},
		{Key: "OriginalOffset", Value: []byte(fmt.Sprintf("%d", msg.Offset))},
	}

	msgJson, err := json.Marshal(dlm)
	if err != nil {
		return
	}

	err = DeadLetterWriter.PublishMessage("", string(msgJson), header...)
	if err != nil {
		log.Error("Failed to publish dead letter message", log.String(err.Error()))
	}
}
