package worker

import "forum/model/kafka/reader"

func Run() {
	//------------------- kafka -----------------------
	WriterInit()

	feedReader := reader.NewKafkaReader(&reader.TopicsAndGroup{
		Topics:  []string{kafkaTopic},
		GroupID: "forum_feed_group",
	})

	go feedReader.BeginConsume(
		nil,
		handleFeedMessage,
		nil,
		reader.DefaultErrorHandler,
		3,
	)
}
