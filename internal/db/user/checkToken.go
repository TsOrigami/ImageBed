package user

import (
	rds "ImageV2/internal/db/redis"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

func CheckToken(token string) bool {
	redisClient, err := rds.GetRedis()
	if err != nil {
		return false
	}
	userID, err := redisClient.Connect.Get(redisClient.Ctx, token).Result()
	if errors.Is(err, redis.Nil) {
		return false
	} else if err != nil {
		return false
	}
	redisClient.Connect.Set(redisClient.Ctx, token, userID, 5*time.Minute)
	return true
}
