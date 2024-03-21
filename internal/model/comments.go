package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Message   string    `json:"message"`
	UserID    uint64    `json:"user_id"`
	PhotoID   uint64    `json:"photo_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`
}
type CreateComment struct {
	ID        uint64    `json:"id" `
	Message   string    `json:"message" binding:"required"`
	UserID    uint64    `json:"user_id"`
	PhotoID   uint64    `json:"photo_id" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}
type GetCommentByID struct {
	ID        uint64    `json:"id" `
	Message   string    `json:"message" binding:"required"`
	PhotoID   uint64    `json:"photo_id" binding:"required"`
	UserID    uint64    `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	User      struct {
		ID       uint64 `json:"Id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Photo struct {
		ID       uint64 `json:"id"`
		Title    string `json:"title"`
		Caption  string `json:"caption"`
		PhotoURL string `json:"photo_url"`
		UserID   uint64 `json:"user_id"`
	} `json:"photo,omitempty" gorm:"foreignKey:PhotoID"`
	
}
type UpdateComment struct {
	ID        uint64    `json:"id" `
	Message   string    `json:"message" binding:"required"`
	UserID    uint64    `json:"user_id"`
	PhotoID   uint64    `json:"photo_id" binding:"required"`
	UpdatedAt time.Time `json:"updated_at"`
}