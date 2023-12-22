package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Http *HttpConfig
}

type HttpConfig struct {
	Port string
}

const GO_ENV = "development"

func Init() (*Config, error) {
	err := godotenv.Load(fmt.Sprintf(".env.%s", GO_ENV))
	if err != nil {
		return nil, err
	}

	httpConfig := &HttpConfig{Port: os.Getenv("PORT")}

	return &Config{
		Http: httpConfig,
	}, nil
}
