package app

import (
	"github.com/kamkali/go-timeline/internal/auth"
	"github.com/kamkali/go-timeline/internal/config"
	postgresql2 "github.com/kamkali/go-timeline/internal/postgresql"
	"github.com/kamkali/go-timeline/internal/server"
	service2 "github.com/kamkali/go-timeline/internal/service"
	timeline2 "github.com/kamkali/go-timeline/internal/timeline"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
)

type app struct {
	log        *zap.Logger
	config     *config.Config
	database   *gorm.DB
	jwtManager *auth.JWTManager
	server     *server.Server

	eventRepo      timeline2.EventRepository
	eventService   timeline2.EventService
	typeRepo       timeline2.TypeRepository
	typeService    timeline2.TypeService
	userService    timeline2.UserService
	userRepository timeline2.UserRepository
}

func (a *app) initConfig() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot initialize config for app: %v\n", err)
	}
	a.config = c
}

func (a *app) initLogger() {
	var zapLogger *zap.Logger
	switch a.config.Stage {
	case config.StageProduction:
		l, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("cannot init zapLogger")
		}
		zapLogger = l
	case config.StageDevelopment:
		l, err := zap.NewDevelopment()
		if err != nil {
			log.Fatalf("cannot init zapLogger")
		}
		zapLogger = l
	case config.StageTest:
		fallthrough
	default:
		l := zap.NewExample()
		zapLogger = l
	}

	a.log = zapLogger
}

func (a *app) initDB() {
	orm, err := postgresql2.NewDB(a.config)
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
	a.initJWTManager()
	a.initHTTPServer()
}

func (a *app) initTimelineRepositories() {
	a.eventRepo = postgresql2.NewEventRepository(a.log, a.database)
	a.typeRepo = postgresql2.NewTypeRepository(a.log, a.database)
	a.userRepository = postgresql2.NewUserRepository(a.log, a.database)
}

func (a *app) initTimelineServices() {
	a.eventService = service2.NewEventService(a.log, a.eventRepo)
	a.typeService = service2.NewTypeService(a.log, a.typeRepo)
	a.userService = service2.NewUserService(a.log, a.userRepository)
}

func (a *app) initJWTManager() {
	manager, err := auth.NewJWTManager(a.log, a.config.Auth.SecretKey, a.config.Auth.PublicKey)
	if err != nil {
		log.Fatalf("cannot instantiate JWT Manager")
	}
	a.jwtManager = manager
}

func (a *app) initHTTPServer() {
	s, err := server.New(
		a.config,
		a.log,
		a.jwtManager,
		a.eventService, a.typeService, a.userService,
	)
	if err != nil {
		log.Fatalf("cannot init server: %v\n", err)
	}
	a.server = s
}

func (a *app) start() {
	if err := postgresql2.Migrate(a.database); err != nil {
		log.Fatalf("couldn't migrate db: %v\n", err)
	}
	a.log.Info("successfully migrated database")
	if a.config.SeedDB {
		if err := a.seedDBWithAdmin(a.config); err != nil {
			log.Fatalf("cannot seed DB with admin info")
		}
		if err := a.seedDBWithExampleValues(); err != nil {
			log.Fatalf("cannot seed DB with example values")
		}
	}

	a.server.Start()
}

func Run() {
	a := app{}
	a.initApp()
	a.start()
}
