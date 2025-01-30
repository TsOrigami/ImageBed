package user

import (
	rds "ImageV2/internal/db/redis"
	"time"
)

func SetUserToken(userID string, token string) error {
	redisClient, err := rds.GetRedis()
	if err != nil {
		return err
	}
	redisClient.Connect.Set(redisClient.Ctx, token, userID, time.Duration(redisClient.Remains)*time.Minute)
	return nil
}
