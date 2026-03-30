package limiter

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"
	"golang.org/x/time/rate"
)

const (
	// 限流只作用在创建操作(post,comment,report)的时候
	// 如果用户一周不进行创建操作，删除其限流桶，节约内存空间
	defaultTTL = 7 * 24 * time.Hour
)

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen int64
}

type LimiterManager struct {
	mu  sync.RWMutex
	m   map[uint32]*limiterEntry
	ttl time.Duration
}

func NewLimiterManager() *LimiterManager {
	l := &LimiterManager{
		m:   make(map[uint32]*limiterEntry),
		ttl: defaultTTL,
	}

	c := cron.New()
	_, _ = c.AddFunc("0 2 * * *", func() {
		l.cleanup(time.Now())
	})
	c.Start()

	return l
}

func (l *LimiterManager) AllowN(userId uint32, n int) bool {
	now := time.Now()
	entry := l.getOrCreateEntry(userId, now)
	atomic.StoreInt64(&entry.lastSeen, now.UnixNano())
	return entry.limiter.AllowN(now, n)
}

func (l *LimiterManager) getOrCreateEntry(userId uint32, now time.Time) *limiterEntry {
	l.mu.RLock()
	entry := l.m[userId]
	l.mu.RUnlock()
	if entry != nil {
		return entry
	}

	l.mu.Lock()
	entry = l.m[userId]
	if entry == nil {
		entry = &limiterEntry{
			limiter:  rate.NewLimiter(1, 100), // 容量为100， 每秒产生1个 token
			lastSeen: now.UnixNano(),
		}
		l.m[userId] = entry
	}
	l.mu.Unlock()
	return entry
}

func (l *LimiterManager) cleanup(now time.Time) {
	expireBefore := now.Add(-l.ttl).UnixNano()
	l.mu.Lock()
	for userId, entry := range l.m {
		if atomic.LoadInt64(&entry.lastSeen) < expireBefore {
			delete(l.m, userId)
		}
	}
	l.mu.Unlock()
}
