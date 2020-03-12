package db

import (
	"fmt"
	"blachat-server/config"
	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func InitRedis () {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.GetConfig().GetString("redis_host") + ":" + config.GetConfig().GetString("redis_port"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := redisClient.Ping().Result()
	fmt.Println(pong, err)
}

func GetRedisClient() *redis.Client {
	return redisClient
}