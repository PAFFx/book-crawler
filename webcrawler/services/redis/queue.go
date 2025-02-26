package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"

	"book-search/webcrawler/config"
)

var redisClient *redis.Client
var env *config.EnvVariables

func InitRedisClient(ctx context.Context) error {
	if env == nil {
		var err error
		env, err = config.GetEnv()
		if err != nil {
			log.Fatalf("Error getting environment variables: %v", err)
		}
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     env.RedisHost,
		Password: env.RedisPassword,
		DB:       env.RedisDB,
	})

	_, err := redisClient.Ping(ctx).Result()
	return err
}

func PushURLQueue(url string) {
	ctx := context.Background()
	err := redisClient.RPush(ctx, "urlQueue", url).Err()
	if err != nil {
		panic(err)
	}
}

func PopURLQueue() string {
	ctx := context.Background()
	url, err := redisClient.LPop(ctx, "urlQueue").Result()
	if err != nil {
		panic(err)
	}
	return url
}
