package server

import (
	"github.com/Toflex/directory_v2/internal/email"
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/verification"
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
	queues["verification"] = 1

	opt, err := asynq.ParseRedisURI(conf.RedisURL)
	if err != nil {
		log.WithError(err).Error("Failed to parse redus url")
		return nil, err
	}

	srv := asynq.NewServer(opt,
		asynq.Config{
			Concurrency: conf.Concurrency,
			Queues:      queues,
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(constant.TaskTypeWelcomeEmail, email.HandleWelcomeEmailTask)
	mux.HandleFunc(constant.TaskTypeOTPEmail, email.HandleOTPEmailTask)
	mux.HandleFunc(constant.TaskTypeIdentityVerification, verification.ProcessVerificationTask)
	mux.HandleFunc(constant.TaskUserVerification, email.HandleVericationEmailTask)

	return &Server{srv: srv, mux: mux}, nil
}

func (s *Server) Run() error {
	return s.srv.Run(s.mux)
}

func (s *Server) Shutdown() {
	s.srv.Shutdown()
}
