package queue

import (
	"encoding/json"
	"time"

	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/hibiken/asynq"
	"github.com/samber/do/v2"
)

type Client struct {
	client *asynq.Client
}

var ClientQueue *Client

type TaskPayload struct {
	TaskID    string
	QueueName string
	Data      interface{}
	Retention time.Duration
	Retry     int
	WaitTime  time.Duration
	Timeout   time.Duration
}

type asynqConfig struct {
	RedisURL      string `env:"REDIS_URL"`
	RedisDB       int    `env:"REDIS_DB"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	Concurrency   int    `env:"ASYNCQ_CONCURRENCY"`
}

func NewClient(i do.Injector) (*Client, error) {
	conf := &asynqConfig{}
	configuration.Load(conf)

	opt, err := asynq.ParseRedisURI(conf.RedisURL)
	if err != nil {
		log.WithError(err).Error("Failed to parse redus url")
		return nil, err
	}

	asynqclient := asynq.NewClient(opt)

	ClientQueue = &Client{client: asynqclient}

	return ClientQueue, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return c.client.Enqueue(task, opts...)
}

func Enqueue(taskType string, payload TaskPayload) error {
	data, err := json.Marshal(payload.Data)
	if err != nil {
		return err
	}

	task := asynq.NewTask(taskType, data)

	_, err = ClientQueue.Enqueue(task,
		asynq.Queue(payload.QueueName),
		asynq.ProcessIn(payload.WaitTime),
		asynq.MaxRetry(payload.Retry),
		asynq.Retention(payload.Retention),
		asynq.TaskID(payload.TaskID),
		asynq.Timeout(payload.Timeout))
	if err != nil {
		return err
	}

	return nil
}
