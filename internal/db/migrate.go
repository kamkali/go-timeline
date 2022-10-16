package db

import (
	"github.com/kamkali/go-timeline/internal/db/schema"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&schema.Type{},
		&schema.Event{},
		&schema.Process{},
	)
}
