package app

import (
	"github.com/kamkali/go-timeline/internal/auth"
	"github.com/kamkali/go-timeline/internal/config"
	"github.com/kamkali/go-timeline/internal/db"
	"github.com/kamkali/go-timeline/internal/domain"
	"github.com/kamkali/go-timeline/internal/domain/eventservice"
	"github.com/kamkali/go-timeline/internal/domain/processservice"
	"github.com/kamkali/go-timeline/internal/domain/typeservice"
	"github.com/kamkali/go-timeline/internal/domain/userservice"
	"github.com/kamkali/go-timeline/internal/logger"
	"github.com/kamkali/go-timeline/internal/server"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"log"
	"strings"
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
	userService    domain.UserService
	userRepository domain.UserRepository

	jwtManager *auth.JWTManager
	server     *server.Server
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
	a.initJWTManager()
	a.initHTTPServer()
}

func (a *app) initTimelineRepositories() {
	a.eventRepo = db.NewEventRepository(a.log, a.database)
	a.typeRepo = db.NewTypeRepository(a.log, a.database)
	a.processRepo = db.NewProcessRepository(a.log, a.database)
	a.userRepository = db.NewUserRepository(a.log, a.database)
}

func (a *app) initTimelineServices() {
	a.eventService = eventservice.New(a.log, a.eventRepo)
	a.typeService = typeservice.New(a.log, a.typeRepo)
	a.processService = processservice.New(a.log, a.processRepo)
	a.userService = userservice.New(a.log, a.userRepository)
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
		a.eventService, a.typeService, a.processService, a.userService,
	)
	if err != nil {
		log.Fatalf("cannot init server: %v\n", err)
	}
	a.server = s
}

func (a *app) start() {
	if err := db.Migrate(a.database); err != nil {
		log.Fatalf("couldn't migrate db: %v\n", err)
	}
	a.log.Debug("successfully migrated database")
	if a.config.SeedDB {
		if err := a.seedDBWithAdmin(a.config); err != nil {
			// little hack for development purposes
			if !strings.Contains(errors.Cause(err).Error(), "duplicate key value violates unique constraint \"idx_email\"") {
				log.Fatalf("cannot seed DB with admin info")
			}
		}
	}

	a.server.Start()
}

func (a *app) seedDBWithAdmin(c *config.Config) error {
	u := domain.User{
		Email:    c.AdminEmail,
		Password: c.AdminPassword,
	}
	if err := a.userService.CreateUser(context.Background(), u); err != nil {
		return err
	}
	a.log.Info("Seeded the DB with admin user")
	return nil
}

func Run() {
	a := app{}
	a.initApp()
	a.start()
}
