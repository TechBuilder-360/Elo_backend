package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"
)

type Client struct {
	rdb       *redis.Client
	Namespace string
}

var rdb *redis.Client

type config struct {
	RedisURL          string `env:"REDIS_URL"`
	RedisDB           int    `env:"REDIS_DB"`
	RedisPassword     string `env:"REDIS_PASSWORD"`
	RedisCacheRefresh string `env:"REDIS_CACHE_REFRESH"`
}

func connectRedis() *redis.Client {
	conf := &config{}
	co := configuration.Load(conf)
	conf = co.(*config)

	rdb = redis.NewClient(&redis.Options{
		Addr:     conf.RedisURL,
		Password: conf.RedisPassword,
		DB:       conf.RedisDB,
	})

	// Test redis connection
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Panic("unable to connect to redis: %s", err)
	}

	log.Info("connected to redis DB")
	return rdb
}

func NewClient(i do.Injector) (*Client, error) {
	return &Client{rdb: connectRedis(), Namespace: configuration.Instance.Namespace}, nil
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	key = fmt.Sprintf("%s-%s", c.Namespace, key)
	return rdb.Get(ctx, key).Result()
}

func (c *Client) Close() {
	err := c.rdb.Close()
	if err != nil {
		log.Errorf("Failed to close redis connection: %v", err)
	}
}
