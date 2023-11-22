package rate_limiter

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed rate_limiter.lua
var rateLimiterScript string

// SlideWindowOpt is the option of slide window
type SlideWindowOpt struct {
	WindowTime int64  // ms
	Threshold  int64  // times
	Key        string // key
	UniqueId   string // unique id
}

// NewSlideWindowOpt 创建一个滑动窗口限流配置
func NewSlideWindowOpt(windowTime, threshold int64, opts ...Option) *SlideWindowOpt {
	swOpt := &SlideWindowOpt{
		WindowTime: windowTime,
		Threshold:  threshold,
	}
	for _, opt := range opts {
		opt(swOpt)
	}
	return swOpt
}

type Option func(*SlideWindowOpt)

// WithUniqueId 设置唯一标识
func WithUniqueId(uniqueId string) Option {
	return func(swOpt *SlideWindowOpt) {
		swOpt.UniqueId = uniqueId
	}
}

// WithKey 设置key
func WithKey(key string) Option {
	return func(swOpt *SlideWindowOpt) {
		swOpt.Key = key
	}
}

// Limiter is the rate limiter
type Limiter struct {
	SWOpt       *SlideWindowOpt
	RedisClient *redis.Client
}

// NewLimiter 创建一个限流器
func NewLimiter(swOpt *SlideWindowOpt, redisClient *redis.Client) *Limiter {
	return &Limiter{
		SWOpt:       swOpt,
		RedisClient: redisClient,
	}
}

// CheckLimited 检查是否限流, false: 限流, true: 不限流
func (l *Limiter) CheckLimited() bool {
	res := l.RedisClient.Eval(
		context.Background(),
		rateLimiterScript,
		[]string{l.SWOpt.Key},
		l.SWOpt.WindowTime,
		l.SWOpt.Threshold,
		time.Now().UnixMilli(),
		l.SWOpt.UniqueId,
	).Val()
	return res.(int64) == 1
}
