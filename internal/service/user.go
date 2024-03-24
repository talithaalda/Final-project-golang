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
	EditUser(ctx context.Context, id uint64, user model.UserUpdateInput) (model.UserUpdate, error)
	// activity
	SignUp(ctx context.Context, userSignUp model.UserSignUp) (*model.PrintUser, error)

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


func (u *userServiceImpl) EditUser(ctx context.Context, id uint64, user model.UserUpdateInput) (model.UserUpdate, error) {
	userModel, err := u.repo.GetUsersByID(ctx, id)
	if err != nil {
		return model.UserUpdate{}, err
	}
	// Check if email is not empty and user with this email already exists
	if user.Email != "" && user.Email != userModel.Email {
		userEmail, _ := u.repo.GetUserByEmail(ctx, user.Email)
		if userEmail.ID != 0 {
			return model.UserUpdate{}, errors.New("email already exists")
		}

		// Check if email is valid
		if !helper.IsValidEmail(user.Email) {
			return model.UserUpdate{}, errors.New("invalid email format")
		}

		// If email is being updated, update user email
		userModel.Email = user.Email
	}
	// Check if username is not empty and user with this username already exists
	if user.Username != "" && user.Username != userModel.Username {
		userUsername, _ := u.repo.GetUserByUsername(ctx, user.Username)
		if userUsername.ID != 0 {
			return model.UserUpdate{}, errors.New("username already exists")
		}

		// If username is being updated, update user username
		userModel.Username = user.Username
	}
	userModel.UpdatedAt = time.Now()
	// Call repository to edit user
	updatedUser, err := u.repo.EditUser(ctx, id, userModel)
	if err != nil {
		return model.UserUpdate{}, err
	}
	userUpdate := model.UserUpdate{
		ID:       updatedUser.ID,
		Username: updatedUser.Username,
		Email:    updatedUser.Email,
		Dob:      updatedUser.Dob.String(),
		UpdatedAt: updatedUser.UpdatedAt,
	}
	return userUpdate, nil
}
func (u *userServiceImpl) SignUp(ctx context.Context, userSignUp model.UserSignUp) (*model.PrintUser, error) {
	dob, err := time.Parse("2006-01-02", userSignUp.Dob)
	if err != nil {
		return nil, errors.New("invalid date of birth format")
	}

	// Hitung usia berdasarkan tanggal lahir
	today := time.Now()
	age := today.Year() - dob.Year()
	if today.Month() < dob.Month() || (today.Month() == dob.Month() && today.Day() < dob.Day()) {
		age--
	}

	// Periksa apakah usia >= 8 tahun
	if age < 8 {
		return nil, errors.New("age must be at least 8 years old")
	}
	user := model.User{
		Username: userSignUp.Username,
		Email:    userSignUp.Email,
		Dob:      dob,
	}
	// encryption password
	// hashing
	pass, err := helper.GenerateHash(userSignUp.Password)
	if err != nil {
		return nil, err
	}
	user.Password = pass
	// store to db
	createdUser, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	printUser := model.PrintUser{
		ID:       createdUser.ID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
		Dob:      createdUser.Dob,
	}
	
	return &printUser, err
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
		Dob:           user.Dob,
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