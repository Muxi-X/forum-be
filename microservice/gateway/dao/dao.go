package dao

import (
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
