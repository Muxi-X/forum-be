package reader

import (
	"context"
	"fmt"
	"time"

	"forum/log"
	"forum/model"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/spf13/viper"
)

type TopicsAndGroup struct {
	Topics  []string
	GroupID string
}

type KafkaReader struct {
	Self *kafka.Reader
}

const (
	backoff       = 5 * time.Second
	commitTimeout = 5 * time.Second
	dedupTTL      = 24 * time.Hour
)

func NewKafkaReader(tg *TopicsAndGroup) *KafkaReader {
	// SASL PLAIN 配置
	mechanism := plain.Mechanism{
		Username: viper.GetString("kafka.username"),
		Password: viper.GetString("kafka.password"),
	}

	dialer := &kafka.Dialer{
		SASLMechanism: mechanism,
		TLS:           nil,
	}

	return &KafkaReader{
		Self: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{viper.GetString("kafka.addr")},
			GroupTopics: tg.Topics,
			GroupID:     tg.GroupID,
			Dialer:      dialer,
		}),
	}
}

func (r *KafkaReader) Close() error {
	return r.Self.Close()
}

type ErrorHandler func(msg kafka.Message, err error) Action

type Action int

const (
	ActionRetry      Action = iota // 可重试错误，退避后继续
	ActionDeadLetter               // 不可重试，进死信并 commit 跳过
	ActionStop                     // kafka错误，退出循环
)

// PreventDuplicate 防止重复消费
type PreventDuplicate func(msg kafka.Message) (duplicated bool, err error)

func (r *KafkaReader) BeginConsume(
	preventDuplicate PreventDuplicate,
	consumeLogic func(msg kafka.Message) error,
	deadLetter func(msg kafka.Message, err error),
	onError ErrorHandler,
	maxAttempts int,
) {
	if maxAttempts < 1 {
		maxAttempts = 1
	}

	for {
		msg, err := r.Self.FetchMessage(context.Background())
		if err != nil {
			log.Error("Kafka fetch message failed", log.String(err.Error()))
			return
		}

		duplicated, err := preventDuplicate(msg)
		if err != nil {
			log.Error("Kafka prevent duplicate failed", log.String(err.Error()))
		} else if duplicated {
			if err := r.MarkConsumeComplete(msg); err != nil {
				log.Error("Kafka mark consume failed", log.String(err.Error()))
				return
			}
			continue
		}

		var lastErr error
		success := false

	retryLoop:
		for retries := 0; retries < maxAttempts; retries++ {
			var consumeErr error
			if consumeErr = consumeLogic(msg); consumeErr == nil {
				success = true
				break
			}

			lastErr = consumeErr
			switch onError(msg, consumeErr) {
			case ActionRetry:
				time.Sleep(backoff)
			case ActionDeadLetter:
				break retryLoop
			case ActionStop:
				return
			}
		}

		if !success {
			if deadLetter != nil {
				deadLetter(msg, lastErr)
			}
		}

		if err := r.MarkConsumeComplete(msg); err != nil {
			log.Error("Kafka mark consume failed", log.String(err.Error()))
			return
		}
	}
}

func (r *KafkaReader) MarkConsumeComplete(msg kafka.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), commitTimeout)
	defer cancel()
	return r.Self.CommitMessages(ctx, msg)
}

// ---------------- 默认失败处理 ----------------

func DefaultErrorHandler(msg kafka.Message, err error) Action {
	log.Error(fmt.Sprintf("Kafka消费失败(Topic: %s)", msg.Topic), log.String(err.Error()))

	return ActionDeadLetter
}

// ---------------- 默认防重处理 ----------------

func DefaultPreventDuplicate(msg kafka.Message) (bool, error) {
	key := fmt.Sprintf("forum:kafka_dedup:%s:%d:%d", msg.Topic, msg.Partition, msg.Offset)
	ok, err := model.RedisDB.Self.SetNX(key, 1, dedupTTL).Result()
	if err != nil {
		return false, err
	}
	return !ok, nil
}
