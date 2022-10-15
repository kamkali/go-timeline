package service

import (
    "github.com/kamkali/go-timeline/internal/config"
    "github.com/kamkali/go-timeline/internal/db"
    "github.com/kamkali/go-timeline/internal/logger"
    "gorm.io/gorm"
    "log"
)

type app struct {
    config   *config.Config
    log      logger.Logger
    database *gorm.DB
}

func (a *app) initConfig() {
    c, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("cannot initialize config for app: %v\n", err)
    }
    a.config = c
}

func (a *app) initLogger() {
    l, err := logger.GetLogger(a.config)
    if err != nil {
        log.Fatalf("cannot init logger: %v\n", err)
    }
    a.log = l
}

func (a *app) initDB() {
    orm, err := db.NewDB(a.config)
    if err != nil {
        log.Fatalf("cannot initialize db: %v\n", err)
    }
    a.database = orm
}

func (a *app) initApp() {
    a.initConfig()
    a.initLogger()
    a.initDB()
}

func (a *app) start() {
    if err := db.Migrate(a.database); err != nil {
        log.Fatalf("couldn't migrate db: %v\n", err)
    }
    a.log.Debug("successfully migrated database")
}

func Run() {
    a := app{}
    a.initApp()
    a.start()
}
