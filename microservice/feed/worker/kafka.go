package worker

import (
	"encoding/json"
	"forum-feed/dao"

	"forum/model/kafka/writer"

	"github.com/segmentio/kafka-go"
)

const kafkaTopic = "forum_feed"

var FeedWriter *writer.KafkaWriter

func WriterInit() {
	FeedWriter = writer.NewKafkaWriter(kafkaTopic)
}

func PublishFeedMessage(w *writer.KafkaWriter, key string, value string) error {
	return w.PublishMessage(key, value)
}

func handleFeedMessage(msg kafka.Message) error {
	var feed dao.FeedModel

	if err := json.Unmarshal(msg.Value, &feed); err != nil {
		return err
	}

	if err := feed.Create(); err != nil {
		return err
	}

	return nil
}
