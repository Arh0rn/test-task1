package config

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http-server"`
	Database   `yaml:"db"`
}

type HTTPServer struct {
	Address         string        `yaml:"host" env-default:"localhost:8081"`
	ReadTimeout     time.Duration `yaml:"read-timeout" env-default:"5s"`
	WriteTimeout    time.Duration `yaml:"write-timeout" env-default:"10s"`
	IdleTimeout     time.Duration `yaml:"idle-timeout" env-default:"60s"`
	ShutdownTimeout time.Duration `yaml:"shutdown-timeout" env-default:"5s"`
	AccessTokenTTL  time.Duration `yaml:"access-token-ttl" env-default:"1h"`
	RefreshTokenTTL time.Duration `yaml:"refresh-token-ttl" env-default:"24h"`
	HashCost        int           `env:"HASH_COST" env-required:"true"`
	JWTSecret       string        `env:"JWT_SECRET" env-required:"true"`
}

type Database struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	DBName   string `yaml:"name" env-required:"true"`
	User     string `yaml:"users" env-default:"postgres"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, errors.New("CONFIG_PATH env-var not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("CONFIG_PATH does not exist: %s", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg) // Read config and env to :D
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
