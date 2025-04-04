package utils

import (
	"context"
	"delivery-service/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Redis connection failed:", err)
	} else {
		fmt.Println("Connected to Redis")
	}
}
func updateRedisCache(ctx context.Context, campaigns []models.Campaign) error {
	if len(campaigns) == 0 {
		fmt.Println("No campaigns to cache.")
		return nil
	}

	campaignsJSON, err := json.Marshal(campaigns)
	if err != nil {
		return err
	}
	err = RedisClient.Set(ctx, "active_campaigns", campaignsJSON, 10*time.Minute).Err()
	if err != nil {
		return err
	}

	fmt.Println("Redis cache updated")
	return nil
}
