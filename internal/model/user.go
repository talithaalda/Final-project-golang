package model

import (
	"errors"
	"final_project/pkg/helper"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	// FirstName string    `json:"first_name"`
	// LastName  string    `json:"last_name"`
	Email     string         `json:"email"`
	Password  string         `json:"-"`
	DoB       time.Time      `json:"dob" gorm:"column:dob"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`
	Photos    []Photo   `json:"photos,omitempty"`
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
	// Comments  []Comment `json:"comments,omitempty"`
	// SocialMedias []SocialMedia `json:"social_medias,omitempty"`
}

// https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/
// https://gin-gonic.com/docs/examples/binding-and-validation/
type UserSignUp struct {
	Username string `json:"username" binding:"required"`
	// FirstName string `json:"first_name"`
	// LastName  string `json:"last_name"`
	Password string    `json:"password" binding:"required"`
	Email    string    `json:"email" binding:"required"`
	DoB      time.Time `json:"dob"  binding:"required"`
}

func (u UserSignUp) Validate() error {
	// check username
	if u.Username == "" {
		return errors.New("invalid username")
	}
	if u.Email == "" {
		return errors.New("invalid email")
	}
	if len(u.Password) < 6 {
		return errors.New("invalid password")
	}
	if !helper.IsValidEmail(u.Email) {
		return errors.New("invalid email format")
	}
	minAge := 8
	if time.Since(u.DoB).Hours()/24/365 < float64(minAge) {
		return fmt.Errorf("age must be at least %d years old", minAge)
	}
	return nil
}
