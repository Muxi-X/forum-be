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
	GetUser(id uint32) (*UserModel, error)
	BatchGetUser(ids []uint32) (map[uint32]*UserModel, error)
	GetUserByIds(ids []uint32) ([]*UserModel, error)
	GetUserByEmail(email string) (*UserModel, error)
	UpdatePassword(userID uint32, newPassword string) error
	ListUser(offset, limit, lastId uint32, filter *UserModel) ([]*UserModel, error)
	GetUserByStudentId(studentId string) (*UserModel, error)
	RegisterUser(info *RegisterInfo) error
	AddPublicPolicy(string, uint32) error

	ListMessage() ([]string, error)
	ListPrivateMessage(uint32) ([]string, error)
	CreateMessage(uint32, string) error
	DeleteMessage(uint32) error
	DeleteOneMessage(uint32, string) error
}

// Init init dao
func Init() {
	// lock.Lock()
	// defer lock.Unlock()
	if dao != nil {
		return
	}

	// init db
	model.DB.Init()

	// init redis
	model.RedisDB.Init()

	// init casbin
	model.CB.Init()

	dao = &Dao{
		DB:    model.DB.Self,
		Redis: model.RedisDB.Self,
	}
}

func GetDao() *Dao {
	return dao
}
