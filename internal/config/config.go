package config

import (
    "github.com/kelseyhightower/envconfig"
)

type AppStage int

const (
    StageTest AppStage = iota + 1
    StageDevelopment
    StageProduction
)

type Config struct {
    Stage AppStage

    DB struct {
        Host     string `envconfig:"DB_HOST" default:"localhost"`
        Port     string `envconfig:"DB_PORT" default:"5432"`
        User     string `envconfig:"DB_USER" default:"postgres"`
        Password string `envconfig:"DB_PASSWORD" default:"postgres"`
        Name     string `envconfig:"DB_NAME" default:"timeline"`
    }
}

func LoadConfig() (*Config, error) {
    var c Config
    if err := envconfig.Process("timeline", &c); err != nil {
        return nil, err
    }
    return &c, nil
}
