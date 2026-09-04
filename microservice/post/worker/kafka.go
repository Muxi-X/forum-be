package worker

import (
	"encoding/json"
	"strconv"
	"time"

	"forum-post/dao"
	"forum/model/kafka/writer"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

const (
	TypeComment    = "comment"
	TypeLike       = "like"
	TypeCollection = "collection"
)

const (
	interactionTopic = "forum_post_interaction"
	interactionGroup = "post-interaction-score"
)

type InteractionMessage struct {
	Type       string `json:"type"`
	PostId     uint32 `json:"post_id"`
	Score      int    `json:"score"`
	EventId    string `json:"event_id"` // 事件id，防止markConsume失败导致重复消费
	NeedRecord bool   `json:"need_record"`
}

var PostInteractionWriter *writer.KafkaWriter

func WriterInit() {
	PostInteractionWriter = writer.NewKafkaWriter(interactionTopic)
}

func PublishPostInteraction(w *writer.KafkaWriter, msg *InteractionMessage) error {
	if msg.EventId == "" {
		msg.EventId = uuid.New().String()
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return w.PublishMessage(strconv.Itoa(int(msg.PostId)), string(data))
}

func handlePostInteraction(msg kafka.Message) error {
	var im InteractionMessage
	if err := json.Unmarshal(msg.Value, &im); err != nil {
		return err
	}

	d := dao.GetDao()
	key := "forum:kafka_event_id:" + im.EventId
	if ok, err := d.SetNX(key, 1, 24*time.Hour); err != nil {
		return err
	} else if !ok {
		return nil
	}

	if im.NeedRecord {
		if err := d.AddChangeRecord(im.PostId); err != nil {
			return err
		}
	}
	if im.Score != 0 {
		if err := d.ChangePostScore(im.PostId, im.Score); err != nil {
			return err
		}
	}

	return nil
}
