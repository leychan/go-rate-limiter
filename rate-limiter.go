package rate_limiter

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed rate_limiter.lua
var rateLimiterScript string

type SlideWindowOpt struct {
	WindowTime int64
	Threshold  int64
	Key        string
	UniqueId   string // unique id
}

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

func WithUniqueId(uniqueId string) Option {
	return func(swOpt *SlideWindowOpt) {
		swOpt.UniqueId = uniqueId
	}
}

func WithKey(key string) Option {
	return func(swOpt *SlideWindowOpt) {
		swOpt.Key = key
	}
}

type Limiter struct {
	SWOpt       *SlideWindowOpt
	RedisClient *redis.Client
}

func NewLimiter(swOpt *SlideWindowOpt, redisClient *redis.Client) *Limiter {
	return &Limiter{
		SWOpt:       swOpt,
		RedisClient: redisClient,
	}
}

func (l *Limiter) Check() bool {
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
