package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() (*redis.Client, error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	if host == "" {
		host = "redis"
	}
	if port == "" {
		port = "6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Password:     password,
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis bağlantı hatası: %v", err)
	}

	log.Println("✅ Redis bağlantısı başarıyla kuruldu.")
	RedisClient = client
	return client, nil
}
