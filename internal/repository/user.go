package repository

import (
	"context"

	"final_project/internal/infrastructure"
	"final_project/internal/model"

	"gorm.io/gorm"
)

type UserQuery interface {
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUsersByID(ctx context.Context, id uint64) (model.User, error)
	EditUser(ctx context.Context, id uint64, user model.User) (model.User, error)
	DeleteUsersByID(ctx context.Context, id uint64) error
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type UserCommand interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
}

type userQueryImpl struct {
	db infrastructure.GormPostgres
}

func NewUserQuery(db infrastructure.GormPostgres) UserQuery {
	return &userQueryImpl{db: db}
}

func (u *userQueryImpl) GetUsers(ctx context.Context) ([]model.User, error) {
	db := u.db.GetConnection()
	users := []model.User{}
	if err := db.
		WithContext(ctx).
		Table("users").
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userQueryImpl) GetUsersByID(ctx context.Context, id uint64) (model.User, error) {
	db := u.db.GetConnection()
	users := model.User{}
	if err := db.
		WithContext(ctx).
		Table("users").
		Where("id = ?", id).
		Find(&users).Error; err != nil {
		// if user not found, return nil error
		if err == gorm.ErrRecordNotFound {
			return model.User{}, nil
		}

		return model.User{}, err
	}
	return users, nil
}

func (u *userQueryImpl) DeleteUsersByID(ctx context.Context, id uint64) error {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("users").
		Delete(&model.User{ID: id}).
		Error; err != nil {
		return err
	}
	return nil
}

func (u *userQueryImpl) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	db := u.db.GetConnection()
	if err := db.
		WithContext(ctx).
		Table("users").
		Save(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}
func (u *userQueryImpl) EditUser(ctx context.Context, id uint64, photo model.User) (model.User, error) {
	db := u.db.GetConnection()
	updatedUser := model.User{}
	if err := db.
		WithContext(ctx).
		Table("users").
		Where("id = ?", id).Updates(&photo).First(&updatedUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return model.User{}, nil
			}
		}
	return updatedUser, nil
}
func (u *userQueryImpl) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	db := u.db.GetConnection()
	user := model.User{}
	if err := db.WithContext(ctx).Where("email = ?", email).Find(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.User{}, nil
		}
		return model.User{}, err
	}
	return user, nil
}