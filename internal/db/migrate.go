package db

import (
	"github.com/kamkali/go-timeline/internal/db/schema/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Type{},
		&models.Event{},
		&models.Process{},
		&models.User{},
	)
}
