package db

import (
    "fmt"
    "github.com/kamkali/go-timeline/internal/config"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func NewDB(c *config.Config) (*gorm.DB, error) {
    psql := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
        c.DB.Host,
        c.DB.Port,
        c.DB.User,
        c.DB.Password,
        c.DB.Name,
    )

    db, err := gorm.Open(postgres.Open(psql), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("cannot create Postgres DB: %w", err)
    }

    return db, nil
}
