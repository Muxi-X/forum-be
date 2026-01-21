package dao

import (
	pb "forum-chat/proto"
	"forum/model"
	"time"

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
	Create(*ChatData) error
	GetList(uint32, time.Duration, bool) ([]string, error)
	Rewrite(uint32, []string) error
	ListHistory(uint32, uint32, uint32, uint32, bool) ([]*pb.Message, error)
	CreateHistory(uint32, []string) error
	GetUserList(uint32) ([]uint32, error)
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

	dao = &Dao{
		DB:    model.DB.Self,
		Redis: model.RedisDB.Self,
	}

}

// ChatData 发送到redis里面的数据
type ChatData struct {
	Content  string `json:"content"`
	Time     string `json:"time"`
	Receiver uint32 `json:"-"`
	Sender   uint32 `json:"sender"`
	TypeName string `json:"type_name"`
}

func GetDao() *Dao {
	return dao
}

type DBdata struct {
	ReceiverID uint32
	SenderID   uint32
	Content    string
	Time       string
	TypeName   string
}
