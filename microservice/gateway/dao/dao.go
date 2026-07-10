package dao

import (
	"encoding/json"
	"fmt"
	"time"

	"forum/model"
	"forum/pkg/limiter"
	"forum/pkg/obfuscate"

	"github.com/spf13/viper"
)

var (
	dao *Dao
)

// Dao .
type Dao struct {
	LimiterManager *limiter.LimiterManager
	Obfuscator     *obfuscate.Obfuscator
}

// Interface dao
type Interface interface {
	AllowN(userId uint32, n int) bool
	Obfuscate(id uint32) string
	Deobfuscate(hid string) (uint32, error)

	SavePending(prefix string, id uint, data *PendingData) error
	GetPending(prefix string, id uint) (*PendingData, error)
	DeletePending(prefix string, id uint) error
	FindPending(id uint) (*PendingData, string, error)
	NextPendingID() (uint, error)
}

// Init init dao
func Init() {
	if dao != nil {
		return
	}

	// 黑名单过期数据定时清理
	// go service.TidyBlacklist()
	// 同步黑名单数据
	// service.SynchronizeBlacklistToRedis()

	// init redis
	model.RedisDB.Init()

	// init casbin
	model.CB.Init()

	limiterManager := limiter.NewLimiterManager()
	obfuscator := obfuscate.NewObfuscator(viper.GetString("hashids.salt"), viper.GetInt("hashids.minlength"))

	dao = &Dao{
		LimiterManager: limiterManager,
		Obfuscator:     obfuscator,
	}
}

func GetDao() *Dao {
	return dao
}

func (d Dao) AllowN(userId uint32, n int) bool {
	return d.LimiterManager.AllowN(userId, n)
}

func (d Dao) Obfuscate(id uint32) string {
	return d.Obfuscator.Obfuscate(uint(id))
}

func (d Dao) Deobfuscate(hid string) (uint32, error) {
	id, err := d.Obfuscator.Deobfuscate(hid)
	if err != nil {
		return 0, err
	}

	return uint32(id), nil
}

// PendingData 全量存储原始请求，不做字段删减
// ResourceType 用于 webhook 分发，UserId 用于创建资源，RawRequest 为完整原始请求 JSON
type PendingData struct {
	ResourceType string          `json:"resource_type"`
	UserId       uint32          `json:"user_id"`
	RawRequest   json.RawMessage `json:"raw_request"`
}

const (
	PendingPrefixPost     = "audit:post:"
	PendingPrefixSipScore = "audit:sipscore:"
	pendingIDKey          = "audit:next_id"
	pendingExpiration     = 7 * 24 * time.Hour
)

var AllPendingPrefixes = []string{PendingPrefixPost, PendingPrefixSipScore}

func (d *Dao) SavePending(prefix string, id uint, data *PendingData) error {
	key := fmt.Sprintf("%s%d", prefix, id)
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化待审核数据失败: %w", err)
	}
	return model.SetStringInRedis(key, string(jsonBytes), pendingExpiration)
}

func (d *Dao) GetPending(prefix string, id uint) (*PendingData, error) {
	key := fmt.Sprintf("%s%d", prefix, id)
	val, found, err := model.GetStringFromRedis(key)
	if err != nil {
		return nil, fmt.Errorf("读取待审核数据失败: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("待审核数据不存在: prefix=%s id=%d", prefix, id)
	}

	var data PendingData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, fmt.Errorf("反序列化待审核数据失败: %w", err)
	}
	return &data, nil
}

func (d *Dao) DeletePending(prefix string, id uint) error {
	key := fmt.Sprintf("%s%d", prefix, id)
	return model.RedisDB.Self.Del(key).Err()
}

func (d *Dao) FindPending(id uint) (*PendingData, string, error) {
	for _, prefix := range AllPendingPrefixes {
		data, err := d.GetPending(prefix, id)
		if err == nil {
			return data, prefix, nil
		}
	}
	return nil, "", fmt.Errorf("待审核数据不存在: id=%d", id)
}

// NextPendingID 生成全局唯一 pendingID（所有资源类型共享）
func (d *Dao) NextPendingID() (uint, error) {
	id, err := model.RedisDB.Self.Incr(pendingIDKey).Result()
	if err != nil {
		return 0, fmt.Errorf("生成待审核ID失败: %w", err)
	}
	return uint(id), nil
}
