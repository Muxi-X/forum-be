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
		handlePostInteraction,
		writer.DefaultEnterDeadLetter,
		reader.DefaultErrorHandler,
		3,
	)
}
