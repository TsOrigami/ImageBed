package redis

import (
	conf "ImageV2/configs"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
)

import "github.com/go-redis/redis/v8"

type Redis struct {
	Connect  *redis.Client
	Ctx      context.Context
	addr     string
	password string
	db       int
}

var (
	instance *Redis
	mu       sync.Mutex
	ctx      = context.Background()
)

func GetRedis() (*Redis, error) {
	if instance == nil {
		mu.Lock()
		defer mu.Unlock()
		if instance == nil {
			redisClient, addrInfo, passwdInfo, dbInfo, err := ConnectionRedis()
			if err != nil {
				return nil, err
			}
			instance = &Redis{
				Connect:  redisClient,
				Ctx:      ctx,
				addr:     addrInfo,
				password: passwdInfo,
				db:       dbInfo,
			}
		}
	}
	return instance, nil
}

func ConnectionRedis() (*redis.Client, string, string, int, error) {
	jsonData, err := conf.GetConfigGroupAsJSON("redis")
	if err != nil {
		return nil, "", "", 0, err
	}
	var redisConfig = make(map[string]string)
	err = json.Unmarshal(jsonData, &redisConfig)
	if err != nil {
		return nil, "", "", 0, err
	}
	AddrInfo := redisConfig["addr"] + ":" + redisConfig["port"]
	PasswordInfo := redisConfig["password"]
	DbInfo, err := strconv.Atoi(redisConfig["db"])
	if err != nil {
		return nil, "", "", 0, err
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     AddrInfo,
		Password: PasswordInfo,
		DB:       DbInfo,
	})
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, "", "", 0, err
	}
	fmt.Printf("Redis 连接成功，返回信息: %s\n", pong)
	return redisClient, AddrInfo, PasswordInfo, DbInfo, nil
}
