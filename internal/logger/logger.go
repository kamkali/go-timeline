package logger

import (
    "github.com/kamkali/go-timeline/internal/config"
    "go.uber.org/zap"
    "sync"
)

var (
    logger Logger
    once   sync.Once
)

type Logger interface {
    Debug(msg string)
    Warn(msg string)
    Error(msg string)
    Info(msg string)
}

type ZapLogger struct {
    logger *zap.Logger
}

func GetLogger(config *config.Config) (l Logger, err error) {
    if logger == nil {
        z := ZapLogger{}
        once.Do(func() {
            err = z.init(config)
            if err != nil {
                return
            }
            logger = z
        })
    }
    return logger, nil
}

func (z *ZapLogger) init(c *config.Config) (err error) {
    var l *zap.Logger
    switch c.Stage {
    case config.StageProduction:
        l, err = zap.NewProduction()
        if err != nil {
            return err
        }
    case config.StageDevelopment:
        l, err = zap.NewDevelopment()
        if err != nil {
            return err
        }
    case config.StageTest:
        fallthrough
    default:
        l = zap.NewExample()
    }
    z.logger = l

    return nil
}

func (z ZapLogger) Debug(msg string) {
    z.logger.Debug(msg)
}

func (z ZapLogger) Warn(msg string) {
    z.logger.Warn(msg)
}

func (z ZapLogger) Error(msg string) {
    z.logger.Error(msg)
}

func (z ZapLogger) Info(msg string) {
    z.logger.Info(msg)
}
