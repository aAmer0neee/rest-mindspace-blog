package domain

import (
	"time"
)

type ArticleEntity struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title" bindinng:"required"`
	Preview   string    `json:"preview"`
	Content   string    `json:"content" binding:"required"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserEntity struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
