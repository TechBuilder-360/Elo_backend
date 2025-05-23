package configuration

import (
	"go.deanishe.net/env"
)

type ENVIRONMENT string

const (
	production  ENVIRONMENT = "PRODUCTION"
	development ENVIRONMENT = "DEVELOPMENT"
	sandbox     ENVIRONMENT = "SANDBOX"
)

var Instance *baseConfig

type baseConfig struct {
	AppName       string      `env:"APP_NAME"`
	Namespace     string      `env:"NAMESPACE"`
	BASEURL       string      `env:"BASE_URL"`
	Port          string      `env:"PORT"`
	Environment   ENVIRONMENT `env:"ENVIRONMENT"`
	Secret        string      `env:"SECRET"`
	TOKENLIFESPAN uint        `env:"TOKEN_LIFE_SPAN"`
}

func LoadBaseConfiguration() {
	c := &baseConfig{}
	if err := env.Bind(c); err != nil {
		panic(err.Error())
	}

	Instance = c
}

func Load(conf interface{}) interface{} {
	if err := env.Bind(conf); err != nil {
		panic(err.Error())
	}

	return conf
}

func (c *baseConfig) GetEnv() ENVIRONMENT {
	return c.Environment
}

func (c *baseConfig) IsProduction() bool {
	return c.Environment == production
}

func (c *baseConfig) IsSandbox() bool {
	return c.Environment == sandbox
}

func (c *baseConfig) IsDevelopment() bool {
	return c.Environment == development
}
