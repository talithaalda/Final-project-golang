package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Photo struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `json:"photo_url"`
	UserID    uint64    `json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`
	Comments []Comment `json:"comments,omitempty"`
	User      struct {
		ID       uint64 `json:"-"`
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"user,omitempty" gorm:"foreignKey:UserID"`
	// User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	
}
type CreatePhoto struct {
	ID        uint64    `json:"id" `
	Title string `json:"title" binding:"required"`
	PhotoURL string    `json:"photo_url" binding:"required"`
	Caption string    `json:"caption" `
	UserID    uint64    `json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
}
type GetPhoto struct {
	ID        uint64    `json:"id" `
	Title string `json:"title" binding:"required"`
	PhotoURL string    `json:"photo_url" binding:"required"`
	Caption string    `json:"caption" `
	UserID    uint64    `json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type UpdatePhoto struct {
	ID        uint64    `json:"id" `
	Title 		string `json:"title" binding:"required"`
	PhotoURL string    `json:"photo_url" binding:"required"`
	Caption string    `json:"caption" binding:"required"`
	UserID    uint64    `json:"user_id"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// type DefaultColumn struct {
// 	CreatedAt time.Time      `json:"created_at"`
// 	UpdatedAt time.Time      `json:"updated_at"`
// 	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
// }



// https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/
// https://gin-gonic.com/docs/examples/binding-and-validation/


func (u CreatePhoto) Validate() error {
	// check username
	if u.Title == "" {
		return errors.New("invalid title")
	}
	if u.PhotoURL == "" {
		return errors.New("invalid url")
	}
	return nil
}
func (u UpdatePhoto) Validate() error {
	// check username
	if u.Title == "" {
		return errors.New("invalid title")
	}
	if u.PhotoURL == "" {
		return errors.New("invalid url")
	}
	return nil
}