package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
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

	opt, err := redis.ParseURL(conf.RedisURL)
	if err != nil {
		log.Panic("Failed to parse redus url %s", err.Error())
	}

	opt.MaxRetries = 10
	opt.DialTimeout = time.Second * 30
	if strings.HasPrefix(strings.ToLower(conf.RedisURL), "rediss://") {
		opt.TLSConfig = &tls.Config{}
	}

	rdb = redis.NewClient(opt)

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
	return c.rdb.Get(ctx, key).Result()
}

func (c *Client) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	key = fmt.Sprintf("%s-%s", c.Namespace, key)
	_, err := c.rdb.Set(ctx, key, value, expiration).Result()
	return err
}

func (c *Client) IncrWithExpire(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	key = fmt.Sprintf("%s-%s", c.Namespace, key)
	count, err := c.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		_, err = c.rdb.Expire(ctx, key, expiration).Result()
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	key = fmt.Sprintf("%s-%s", c.Namespace, key)
	return c.rdb.Del(ctx, key).Err()
}

func (c *Client) Close() {
	err := c.rdb.Close()
	if err != nil {
		log.Errorf("Failed to close redis connection: %v", err)
	}
}
