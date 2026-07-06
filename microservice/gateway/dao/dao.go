package dao

import (
	"encoding/json"
	"fmt"

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

	// 待审核帖子暂存
	SavePendingPost(id uint, data *PendingPost) error
	GetPendingPost(id uint) (*PendingPost, error)
	DeletePendingPost(id uint) error
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

const (
	postKeyPrefix  = "audit:post:"
	postIDKey      = "audit:post:next_id"
	postExpiration = 0 // 不设置过期，webhook 回调后手动删除
)

// PendingPost 待审核帖子临时数据
type PendingPost struct {
	UserId          uint32   `json:"user_id"`
	Content         string   `json:"content"`
	Domain          string   `json:"domain"`
	Title           string   `json:"title"`
	Category        string   `json:"category"`
	ContentType     string   `json:"content_type"`
	Tags            []string `json:"tags"`
	CompiledContent string   `json:"compiled_content"`
	Summary         string   `json:"summary"`
}

// 用于临时保存待审核数据
func (d *Dao) SavePendingPost(id uint, data *PendingPost) error {
	key := fmt.Sprintf("%s%d", postKeyPrefix, id)
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化待审核数据失败: %w", err)
	}
	return model.SetStringInRedis(key, string(jsonBytes), postExpiration)
}

func (d *Dao) GetPendingPost(id uint) (*PendingPost, error) {
	key := fmt.Sprintf("%s%d", postKeyPrefix, id)
	val, found, err := model.GetStringFromRedis(key)
	if err != nil {
		return nil, fmt.Errorf("读取待审核数据失败: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("待审核数据不存在: id=%d", id)
	}

	var data *PendingPost
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, fmt.Errorf("反序列化待审核数据失败: %w", err)
	}
	return data, nil
}

func (d *Dao) DeletePendingPost(id uint) error {
	key := fmt.Sprintf("%s%d", postKeyPrefix, id)
	return model.RedisDB.Self.Del(key).Err()
}

func (d *Dao) NextPendingID() (uint, error) {
	id, err := model.RedisDB.Self.Incr(postIDKey).Result()
	if err != nil {
		return 0, fmt.Errorf("生成待审核ID失败: %w", err)
	}
	return uint(id), nil
}
