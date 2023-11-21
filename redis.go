package rate_limiter

import (
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var redisOpt *redis.Options

func NewRedisClient() *redis.Client {
	if redisClient != nil {
		return redisClient
	}
	redisClient = redis.NewClient(redisOpt)
	return redisClient
}
