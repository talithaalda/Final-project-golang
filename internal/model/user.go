package model

import (
	"errors"
	"final_project/pkg/helper"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
    ID        uint64         `json:"id"`
    Username  string         `json:"username"`
    Email     string         `json:"email"`
    Password  string         `json:"-"`
    Age       int            `json:"age" gorm:"column:age"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`
    Photos    []Photo        `json:"photos,omitempty"`
	Comments  []Comment      `json:"comments,omitempty"`
}

// type DefaultColumn struct {
// 	CreatedAt time.Time      `json:"created_at"`
// 	UpdatedAt time.Time      `json:"updated_at"`
// 	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
// }

type UserMediaSocial struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"user_id"`
	Title     string `json:"title"`
	Url       string `json:"url"`
	Photos    []Photo   `json:"photos,omitempty"`
	Comments  []Comment `json:"comments,omitempty"`
	// SocialMedias []SocialMedia `json:"social_medias,omitempty"`
}

// https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/
// https://gin-gonic.com/docs/examples/binding-and-validation/
type UserSignUp struct {
    Username string    `json:"username" binding:"required"`
    Password string    `json:"password" binding:"required"`
    Email    string    `json:"email" binding:"required"`
    Age       int      `json:"age" binding:"required"`
}

func (u UserSignUp) Validate() error {
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if !helper.IsValidEmail(u.Email) {
		return errors.New("invalid email format")
	}
	 minAge := 8
    if u.Age < minAge {
        return fmt.Errorf("age must be at least %d years old", minAge)
    }
	return nil
}
type UserLogin struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
}
