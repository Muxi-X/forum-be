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
	ToggleFollow(uint32, uint32) (bool, error)
	CountFollowing(uint32) (uint32, error)
	CountFollowers(uint32) (uint32, error)
	IsFollowing(uint32, uint32) (bool, error)
	ListFollowing(uint32, uint32, uint32) ([]*FollowListUser, error)
	ListFollowers(uint32, uint32, uint32) ([]*FollowListUser, error)
	BatchCountFollowing([]uint32) (map[uint32]uint32, error)
	BatchCountFollowers([]uint32) (map[uint32]uint32, error)
	BatchIsFollowing(uint32, []uint32) (map[uint32]bool, error)

	ListMessage() ([]string, error)
	ListPrivateMessage(uint32) ([]string, error)
	ListDeduplicatedPrivateMessage(uint32) ([]string, error)
	CreateMessage(uint32, string) error
	CreateOrUpdateInteractionMessage(uint32, string) error
	DeleteMessage(uint32) error
	DeleteOneMessage(uint32, string) error
	MarkAllMessageRead(uint32) error
	MarkOneMessageRead(uint32, string) error
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
	model.DB.Self.AutoMigrate(&UserModel{}, &UserFollowModel{})

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
