package user

import (
	rds "ImageV2/internal/db/redis"
	"time"
)

func GetUsername(token string) (string, error) {
	redisClient, err := rds.GetRedis()
	if err != nil {
		return "", err
	}
	userID, err := redisClient.Connect.Get(redisClient.Ctx, token).Result()
	if err != nil {
		return "", err
	}
	redisClient.Connect.Set(redisClient.Ctx, token, userID, time.Duration(redisClient.Remains)*time.Minute)
	return userID, nil
}
