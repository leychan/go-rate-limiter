package rate_limiter

import (
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestLimiter_Check(t *testing.T) {
	swOpt := NewSlideWindowOpt(1000, 5, WithKey("test3333"))

	redisOpt = &redis.Options{
		Addr: "127.0.0.1:6379",
	}

	redisClient := NewRedisClient()
	type fields struct {
		RedisClient *redis.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "1", fields: fields{
			RedisClient: redisClient,
		}, want: true},
		{name: "2", fields: fields{
			RedisClient: redisClient,
		}, want: true},
		{name: "3", fields: fields{
			RedisClient: redisClient,
		}, want: true},
		{name: "4", fields: fields{
			RedisClient: redisClient,
		}, want: true},
		{name: "5", fields: fields{
			RedisClient: redisClient,
		}, want: true},
		{name: "6", fields: fields{
			RedisClient: redisClient,
		}, want: false},
		{name: "7", fields: fields{
			RedisClient: redisClient,
		}, want: false},
		{name: "8", fields: fields{
			RedisClient: redisClient,
		}, want: false},
		{name: "9", fields: fields{
			RedisClient: redisClient,
		}, want: false},
		{name: "10", fields: fields{
			RedisClient: redisClient,
		}, want: false},
		{name: "11", fields: fields{
			RedisClient: redisClient,
		}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			swOpt.UniqueId = tt.name
			l := &Limiter{
				RedisClient: tt.fields.RedisClient,
				SWOpt:       swOpt,
			}
			if got := l.CheckLimited(); got != tt.want {
				t.Errorf("CheckLimited() = %v, want %v", got, tt.want)
			}
		})
	}
}
