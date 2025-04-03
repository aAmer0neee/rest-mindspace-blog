package repository

import (
	"time"

	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

type Article struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Title     string    `gorm:"type:varchar(100)"`
	Preview   string    `gorm:"type:varchar(200)"`
	Author    string    `gorm:"default:unknown"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type User struct {
	ID       uuid.UUID `gorm:"primaryKey;unique;not null"`
	Username string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	IsAdmin  bool      `gorm:"default:false"`
}
