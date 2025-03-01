package redis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"book-search/webcrawler/config"
)

var redisClient *redis.Client

func InitRedisClient(ctx context.Context) error {
	env, err := config.GetEnv()
	if err != nil {
		return err
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     env.RedisHost,
		Password: env.RedisPassword,
		DB:       env.RedisDB,
	})

	_, err = redisClient.Ping(ctx).Result()
	return err
}

func CloseRedisClient() error {
	return redisClient.Close()
}

func PushURLQueue(url string) error {
	ctx := context.Background()
	return redisClient.RPush(ctx, "urlQueue", url).Err()
}

func PushMultipleURLQueue(urls []string) error {
	ctx := context.Background()
	return redisClient.RPush(ctx, "urlQueue", urls).Err()
}

func PopURLQueue() (string, error) {
	ctx := context.Background()
	url, err := redisClient.LPop(ctx, "urlQueue").Result()
	return url, err
}

func PopMultipleURLQueue(count int) ([]string, error) {
	ctx := context.Background()
	urls, err := redisClient.LPopCount(ctx, "urlQueue", count).Result()
	return urls, err
}

func IsURLQueueEmpty() bool {
	ctx := context.Background()
	return redisClient.LLen(ctx, "urlQueue").Val() == 0
}

func GetURLQueueLength() int64 {
	ctx := context.Background()
	return redisClient.LLen(ctx, "urlQueue").Val()
}
