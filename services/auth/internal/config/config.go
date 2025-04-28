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

	ServicePort int `env:"GRPC_AUTH_SERVICE_PORT"`
	HttpPort    int `env:"HTTP_AUTH_SERVICE_PORT"`

	ServiceTimeout  time.Duration `env:"AUTH_SERVICE_TIMEOUT"`
	JwtSignedString string        `env:"JWT_SIGNED_STRING"`

	EmailSenderAddress  string `env:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string `env:"EMAIL_SENDER_PASSWORD"`

	Github GithubAPI
}

type GithubAPI struct {
	ClientID     string `env:"GITHUB_CLIENT_ID"`
	ClientSecret string `env:"GITHUB_CLIENT_SECRET"`
	RedirectURL  string `env:"GITHUB_REDIRECT_URL"`
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
