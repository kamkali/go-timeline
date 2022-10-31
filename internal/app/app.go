package app

import (
	"github.com/kamkali/go-timeline/internal/config"
	"github.com/kamkali/go-timeline/internal/db"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/domain/eventservice"
	"github.com/kamkali/go-timeline/internal/domain/processservice"
	"github.com/kamkali/go-timeline/internal/domain/typeservice"
	"github.com/kamkali/go-timeline/internal/logger"
	"github.com/kamkali/go-timeline/internal/server"
	"gorm.io/gorm"
	"log"
)

type app struct {
	log      logger.Logger
	config   *config.Config
	database *gorm.DB

	eventRepo      domain.EventRepository
	eventService   domain.EventService
	typeRepo       domain.TypeRepository
	typeService    domain.TypeService
	processRepo    domain.ProcessRepository
	processService domain.ProcessService
	server         *server.Server
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
	a.initTimelineRepositories()
	a.initTimelineServices()
	a.initHTTPServer()
}

func (a *app) initTimelineRepositories() {
	a.eventRepo = db.NewEventRepository(a.log, a.database)
	a.typeRepo = db.NewTypeRepository(a.log, a.database)
	a.processRepo = db.NewProcessRepository(a.log, a.database)
}

func (a *app) initTimelineServices() {
	a.eventService = eventservice.New(a.log, a.eventRepo)
	a.typeService = typeservice.New(a.log, a.typeRepo)
	a.processService = processservice.New(a.log, a.processRepo)
}

func (a *app) initHTTPServer() {
	a.server = server.New(a.config, a.log, a.eventService, a.typeService, a.processService)
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
