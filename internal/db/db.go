package db

import (
	"fmt"
	"github.com/kamkali/go-timeline/internal/config"
	"github.com/kamkali/go-timeline/internal/db/schema/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(c *config.Config) (*gorm.DB, error) {
	var dsn string
	if c.DB.URI != "" {
		dsn = c.DB.URI
	} else {
		dsn = fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
			c.DB.Host,
			c.DB.Port,
			c.DB.User,
			c.DB.Password,
			c.DB.Name,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("cannot create Postgres DB: %w", err)
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Type{},
		&models.Event{},
		&models.Process{},
		&models.User{},
	)
}
