package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"final_project/internal/model"
	"final_project/internal/repository"
	"final_project/pkg/helper"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUsersById(ctx context.Context, id uint64) (model.User, error)
	DeleteUsersById(ctx context.Context, id uint64) (model.User, error)
	EditUser(ctx context.Context, id uint64, user model.User) (model.User, error)
	// activity
	SignUp(ctx context.Context, userSignUp model.UserSignUp) (model.User, error)

	// misc
	GenerateUserAccessToken(ctx context.Context, user model.User) (token string, err error)
	CheckCredentials(ctx context.Context, email string, password string) (model.User, error)
}

type userServiceImpl struct {
	repo repository.UserQuery
}

func NewUserService(repo repository.UserQuery) UserService {
	return &userServiceImpl{repo: repo}
}

func (u *userServiceImpl) GetUsers(ctx context.Context) ([]model.User, error) {
	users, err := u.repo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, err
}

func (u *userServiceImpl) GetUsersById(ctx context.Context, id uint64) (model.User, error) {
	user, err := u.repo.GetUsersByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	return user, err
}


func (u *userServiceImpl) DeleteUsersById(ctx context.Context, id uint64) (model.User, error) {
	user, err := u.repo.GetUsersByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	// if user doesn't exist, return
	if user.ID == 0 {
		return model.User{}, nil
	}

	// delete user by id
	err = u.repo.DeleteUsersByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return user, err
}


func (u *userServiceImpl) EditUser(ctx context.Context, id uint64, user model.User) (model.User, error) {
    // Perform validation or additional checks here if necessary

    // Call repository to edit user
    updatedUser, err := u.repo.EditUser(ctx, id, user)
    if err != nil {
        return model.User{}, err
    }
    return updatedUser, nil
}

func (u *userServiceImpl) SignUp(ctx context.Context, userSignUp model.UserSignUp) (model.User, error) {
	// assumption: semua user adalah user baru
	user := model.User{
		Username: userSignUp.Username,
		Email:    userSignUp.Email,
		Age:      userSignUp.Age,
		// FirstName: userSignUp.FirstName,
		// LastName:  userSignUp.LastName,
	}

	// encryption password
	// hashing
	pass, err := helper.GenerateHash(userSignUp.Password)
	if err != nil {
		return model.User{}, err
	}
	user.Password = pass

	// store to db
	res, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return model.User{}, err
	}
	return res, err
}

func (u *userServiceImpl) GenerateUserAccessToken(ctx context.Context, user model.User) (token string, err error) {
	// generate claim
	now := time.Now()

	claim := model.StandardClaim{
		Jti: fmt.Sprintf("%v", time.Now().UnixNano()),
		Iss: "go-middleware",
		Aud: "golang-006",
		Sub: "access-token",
		Exp: uint64(now.Add(time.Hour).Unix()),
		Iat: uint64(now.Unix()),
		Nbf: uint64(now.Unix()),
	}

	userClaim := model.AccessClaim{
		StandardClaim: claim,
		UserID:        user.ID,
		Username:      user.Username,
		Age:           user.Age,
	}

	token, err = helper.GenerateToken(userClaim)
	return
}


func (u *userServiceImpl) CheckCredentials(ctx context.Context, email string, password string) (model.User, error) {
	// Retrieve user by email
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, err
	}

	// Check if user exists
	if user.ID == 0 {
		return model.User{}, errors.New("user not found")
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return model.User{}, err
	}

	// Credentials are correct, return user
	return user, nil
}