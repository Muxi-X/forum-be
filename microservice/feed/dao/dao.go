package dao

import (
	"forum/model"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	dao *Dao
)

// Dao .
type Dao struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// Interface dao
type Interface interface {
	Create(*FeedModel) (uint32, error)
	Delete(uint32) error
	List(*FeedModel, uint32, uint32, uint32, bool, uint32) ([]*FeedModel, error)

	PublishMsg([]byte) error
}

// Init init dao
func Init() {
	if dao != nil {
		return
	}

	// init db
	model.DB.Init()

	// init redis
	model.RedisDB.Init()

	model.InitKafka(kafkaTopic)

	dao = &Dao{
		DB:    model.DB.Self,
		Redis: model.RedisDB.Self,
	}
}

func GetDao() *Dao {
	return dao
}

const kafkaTopic = "feed"

func (d Dao) PublishMsg(msg []byte) error {
	return model.KafkaWriter.PublishMessage(msg)
}
