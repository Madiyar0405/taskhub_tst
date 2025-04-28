package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"sync"
	"time"
)

type Config struct {
	Env         string `env:"ENV" env-default:"dev"`
	PostgresURI string `env:"POSTGRES_URI" env-required:"true"`

	TokenTTL time.Duration `env:"TOKEN_TTL" env-required:"true"`

	ServicePort int `env:"GRPC_COMPANY_SERVICE_PORT"`
	HttpPort    int `env:"HTTP_COMPANY_SERVICE_PORT"`

	GithubWebhookSecret string `env:"GITHUB_WEBHOOK_SECRET"`

	ServiceTimeout time.Duration `env:"COMPANY_SERVICE_TIMEOUT"`

	Clients ClientsConfig
}

type Client struct {
	Address      string        `env:"CLIENT_ADDRESS"`
	Timeout      time.Duration `env:"CLIENT_TIMEOUT"`
	RetriesCount int           `env:"CLIENT_RETRIES"`
}

type ClientsConfig struct {
	Auth Client `env-prefix:"AUTH_"`
}

var (
	cfg  *Config
	once sync.Once
)

func MustLoad() *Config {
	once.Do(func() {
		err := godotenv.Load("../.env")
		if err != nil {
			panic("No .env file found")
		}

		cfg = &Config{}
		err = cleanenv.ReadEnv(cfg)
		if err != nil {
			panic("failed to read config " + err.Error())
		}
	})

	return cfg
}
