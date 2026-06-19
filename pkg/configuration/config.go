package configuration

import (
	"strings"

	"go.deanishe.net/env"
)

type ENVIRONMENT string

const (
	production  ENVIRONMENT = "PRODUCTION"
	development ENVIRONMENT = "DEVELOPMENT"
	sandbox     ENVIRONMENT = "SANDBOX"
)

var baseURL = map[string]string{
	string(production):  "https://api.elo.com",
	string(development): "https://dev.api.elo.com",
}

var Instance *baseConfig

type baseConfig struct {
	AppName       string      `env:"APP_NAME"`
	Namespace     string      `env:"NAMESPACE"`
	BASEURL       string      `env:"BASE_URL"`
	Port          string      `env:"PORT"`
	Environment   ENVIRONMENT `env:"ENVIRONMENT"`
	JWTSecret     string      `env:"JWT_SECRET"`
	TOKENLIFESPAN uint        `env:"TOKEN_LIFE_SPAN"`
	BasicUsername string      `env:"BASIC_USERNAME"`
	BasicPassword string      `env:"BASIC_PASSWORD"`
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

func GetEnv() ENVIRONMENT {
	return ENVIRONMENT(strings.ToLower(string(Instance.Environment)))
}

func IsProduction() bool {
	return Instance.Environment == production
}

func IsSandbox() bool {
	return Instance.Environment == sandbox
}

func IsDevelopment() bool {
	return Instance.Environment == development
}

func GetBaseURL() string {
	if Instance.BASEURL != "" {
		return Instance.BASEURL
	}

	url, ok := baseURL[strings.ToUpper(string(GetEnv()))]
	if !ok {
		return ""
	}
	return url
}
