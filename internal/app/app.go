package app

import (
    "github.com/kamkali/go-timeline/internal/config"
    "github.com/kamkali/go-timeline/internal/db"
    "github.com/kamkali/go-timeline/internal/domain"
    "github.com/kamkali/go-timeline/internal/domain/eventservice"
    "github.com/kamkali/go-timeline/internal/logger"
    "github.com/kamkali/go-timeline/internal/server"
    "gorm.io/gorm"
    "log"
)

type app struct {
    log      logger.Logger
    config   *config.Config
    database *gorm.DB

    timelineRepo    domain.EventRepository
    timelineService domain.EventService
    server          *server.Server
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
    a.initTimelineRepo()
    a.initTimelineService()
    a.initHTTPServer()
}

func (a *app) initTimelineRepo() {
    a.timelineRepo = db.NewEventRepository(a.log, a.database)
}

func (a *app) initTimelineService() {
    a.timelineService = eventservice.New(a.log, a.timelineRepo)
}

func (a *app) initHTTPServer() {
    a.server = server.New(a.config, a.timelineService, a.log)
}

func (a *app) start() {
    if err := db.Migrate(a.database); err != nil {
        log.Fatalf("couldn't migrate db: %v\n", err)
    }
    a.log.Debug("successfully migrated database")

    a.server.Start()
}

func Run() {
    a := app{}
    a.initApp()
    a.start()
}
