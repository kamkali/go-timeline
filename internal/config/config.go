package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type AppStage string

const (
	StageTest        AppStage = "TEST"
	StageDevelopment AppStage = "DEV"
	StageProduction  AppStage = "PROD"
)

type Config struct {
	Stage  AppStage `envconfig:"STAGE" default:"DEV"`
	SeedDB bool     `envconfig:"SEED_DB" default:"false"`

	Admin struct {
		Email    string `envconfig:"ADMIN_EMAIL"`
		Password string `envconfig:"ADMIN_PASS"`
	}

	DB struct {
		Host     string `envconfig:"DB_HOST" default:"localhost"`
		Port     string `envconfig:"DB_PORT" default:"5432"`
		User     string `envconfig:"DB_USER" default:"postgres"`
		Password string `envconfig:"DB_PASSWORD" default:"postgres"`
		Name     string `envconfig:"DB_NAME" default:"timeline"`
		URI      string `envconfig:"DATABASE_URL"`
	}

	Server struct {
		Host           string `envconfig:"SERVER_HOST"`
		Port           string `envconfig:"PORT" default:"8080" required:"true"`
		TimeoutSeconds uint   `envconfig:"SERVER_TIMEOUT" default:"30"`
	}

	Auth struct {
		SecretKey string `envconfig:"SECRET_KEY" required:"true"`
		PublicKey string `envconfig:"PUBLIC_KEY" required:"true"`
	}
}

func LoadConfig() (*Config, error) {
	var c Config
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := envconfig.Process("timeline", &c); err != nil {
		return nil, err
	}
	return &c, nil
}
