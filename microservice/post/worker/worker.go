package worker

import (
	"forum/model/kafka/reader"
	"forum/model/kafka/writer"
)

func Run() {
	//------------------- kafka -----------------------
	WriterInit()

	postInteractionReader := reader.NewKafkaReader(&reader.TopicsAndGroup{
		Topics:  []string{interactionTopic},
		GroupID: interactionGroup,
	})

	go postInteractionReader.BeginConsume(
		nil, // 追求强一致的情景再防止重复消费吧
		handlePostInteraction,
		writer.DefaultEnterDeadLetter,
		reader.DefaultErrorHandler,
		3,
	)
}
