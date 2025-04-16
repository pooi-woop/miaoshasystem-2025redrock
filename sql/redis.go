package sql

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
)

var ctx = context.Background()
var R *redis.Client

func Init() {
	viper.SetConfigFile("C:/miaoshaSystem/global/redisConfig.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading Redis config file: %v", err)
	}

	// 确保正确读取 Redis 配置
	addr := viper.GetString("DB.addr")
	password := viper.GetString("DB.password") // 确保密码字段正确读取
	db := viper.GetInt("DB.DB")

	R = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // 确保密码字段正确传递
		DB:       db,
	})

	_, err := R.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis successfully")
}
