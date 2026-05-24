package server

import (
	"time"

	"github.com/Toflex/directory_v2/internal/email"
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/hibiken/asynq"
	"github.com/samber/do/v2"
)

type Server struct {
	srv *asynq.Server
	mux *asynq.ServeMux
}

type asynqConfig struct {
	RedisURL      string `env:"REDIS_URL"`
	RedisDB       int    `env:"REDIS_DB"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	Concurrency   int    `env:"ASYNCQ_CONCURRENCY"`
}

func NewServer(i do.Injector) (*Server, error) {
	conf := &asynqConfig{}
	configuration.Load(conf)

	if conf.Concurrency <= 0 {
		conf.Concurrency = 10
	}

	queues := make(map[string]int)
	queues["email"] = 1

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:         conf.RedisURL,
			Password:     conf.RedisPassword,
			DB:           conf.RedisDB,
			PoolSize:     20,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			DialTimeout:  15 * time.Second,
		},
		asynq.Config{
			Concurrency: conf.Concurrency,
			Queues:      queues,
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(constant.TaskTypeWelcomeEmail, email.HandleWelcomeEmailTask)
	mux.HandleFunc(constant.TaskTypeOTPEmail, email.HandleOTPEmailTask)

	return &Server{srv: srv, mux: mux}, nil
}

func (s *Server) Run() error {
	return s.srv.Run(s.mux)
}

func (s *Server) Shutdown() {
	s.srv.Shutdown()
}
