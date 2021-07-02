package server

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var Client *redis.Client

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr: 	      ViperConfig.Redis.Address,
		Password:     ViperConfig.Redis.Password,
		PoolSize:     1000,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Millisecond * time.Duration(60),
	})
	_, err := Client.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("reids init data base ok")
}

// 右边进
func Rpush(key string, val interface{}) {
	err := Client.RPush(key, val).Err()
	if err != nil {
		panic(err)
	}
}