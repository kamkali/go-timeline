package postgresql

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type event struct {
	gorm.Model

	Name                string
	EventTime           time.Time
	ShortDescription    string
	DetailedDescription string
	Graphic             string
	TypeID              uint
	Type                eventType `gorm:"foreignKey:TypeID"`
}

type eventType struct {
	gorm.Model

	Name  string `gorm:"uniqueIndex;not null"`
	Color string

	Events []event `gorm:"foreignKey:TypeID"`
}

type user struct {
	gorm.Model

	Email    string `gorm:"uniqueIndex:idx_email"`
	Password string
}

func (u *user) BeforeSave(tx *gorm.DB) error {
	if u.Password != "" {
		hash, err := hashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hash
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
