package queue

import (
	"time"

	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/hibiken/asynq"
	"github.com/samber/do/v2"
)

type Client struct {
	client *asynq.Client
}

var ClientQueue *Client

type asynqConfig struct {
	RedisURL      string `env:"REDIS_URL"`
	RedisDB       int    `env:"REDIS_DB"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	Concurrency   int    `env:"ASYNCQ_CONCURRENCY"`
}

func NewClient(i do.Injector) (*Client, error) {
	conf := &asynqConfig{}
	configuration.Load(conf)

	asynqclient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     conf.RedisURL,
		Password: conf.RedisPassword,
		DB:       conf.RedisDB,
	})

	ClientQueue = &Client{client: asynqclient}

	return ClientQueue, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return c.client.Enqueue(task, opts...)
}

func Enqueue(taskType string, payload []byte, waitTime int) error {
	task := asynq.NewTask(taskType, payload)

	_, err := ClientQueue.Enqueue(task, asynq.ProcessIn(time.Duration(waitTime)*time.Second))
	if err != nil {
		return err
	}

	return nil
}
