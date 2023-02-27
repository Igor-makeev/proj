package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DBURL                string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	AuthConfig
}

type AuthConfig struct {
	Salt       string
	SigningKey string
	TokenTTL   time.Duration
}

func NewAuthConfig() *AuthConfig {
	return &AuthConfig{
		Salt:       os.Getenv("SALT"),
		SigningKey: os.Getenv("SIGNING_KEY"),
		TokenTTL:   6 * time.Hour,
	}
}

func NewConfig() *Config {
	var cfg Config
	flag.StringVar(&cfg.RunAddress, "a", "127.0.0.1:8080", "server addres to listen on")
	flag.StringVar(&cfg.DBURL, "d", "postgres://gopher:qwerty@localhost:5438/gopher?sslmode=disable", "database connection address")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "127.0.0.1:8000", "address of the accrual system")

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("failed to parse config environment variables")
	}

	logrus.Printf("env variable SERVER_ADDRESS=%v", cfg.RunAddress)
	logrus.Printf("env variable DATABASE_URI=%v", cfg.DBURL)
	logrus.Printf("env variable ACCRUAL_SYSTEM_ADDRESS=%v", cfg.AccrualSystemAddress)

	return &cfg
}
