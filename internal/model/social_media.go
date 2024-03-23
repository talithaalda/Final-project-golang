package model

import (
	"time"

	"gorm.io/gorm"
)

type SocialMedia struct {
	ID        uint64         `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	SocialMediaURL        string         `json:"social_media_url"`
	UserID    uint64         `json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`
	User      struct {
		ID       uint64 `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
type CreateSocialMedia struct {
	ID        uint64         `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" binding:"required"`
	SocialMediaURL        string         `json:"social_media_url" binding:"required"`
	UserID    uint64         `json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
}
type UpdateSocialMedia struct {
	ID        uint64         `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" binding:"required"`
	SocialMediaURL        string         `json:"social_media_url " binding:"required"`
	UserID    uint64         `json:"user_id"`
	UpdatedAt time.Time  	 `json:"updated_at"`
}
