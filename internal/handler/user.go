package handler

import (
	"net/http"
	"strconv"

	"final_project/internal/middleware"
	"final_project/internal/model"
	"final_project/internal/service"
	"final_project/pkg"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	// users
	GetUsers(ctx *gin.Context)
	GetUsersById(ctx *gin.Context)
	DeleteUsersById(ctx *gin.Context)
	EditUser(ctx *gin.Context)
	// activity
	UserSignUp(ctx *gin.Context)
	UserLogin(ctx *gin.Context)
}

type userHandlerImpl struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) UserHandler {
	return &userHandlerImpl{
		svc: svc,
	}
}

// GetUsers godoc
// @Summary Retrieve list of users
// @Description Retrieve a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} model.User "users"
// @Success 200 {object} pkg.ErrorResponse "No user found"
// @Failure 400 {object} pkg.ErrorResponse "Bad request"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /users [get]
func (u *userHandlerImpl) GetUsers(ctx *gin.Context) {
	users, err := u.svc.GetUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	if len(users) == 0 {
        ctx.JSON(http.StatusOK, gin.H{"message": "No user found"})
        return
    }
	ctx.JSON(http.StatusOK, users)
}

// GetUsersById godoc
// @Summary Retrieve user by ID
// @Description Retrieve a user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Success 200 {object} model.User "user"
// @Failure 400 {object} pkg.ErrorResponse "Bad request"
// @Failure 404 {object} pkg.ErrorResponse "User not found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /users/{id} [get]
func (u *userHandlerImpl) GetUsersById(ctx *gin.Context) {
	// get id user
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}
	user, err := u.svc.GetUsersById(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	if user.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "user not found"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}
// EditUser godoc
// @Summary Update user information
// @Description Update information of a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security Bearer
// @Param user body model.UserUpdateInput true "User data"
// @Success 200 {object} model.UserUpdate "updatedUser"
// @Failure 400 {object} pkg.ErrorResponse "Bad request"
// @Failure 401 {object} pkg.ErrorResponse "Unauthorized"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /users/{id} [put]
func (u *userHandlerImpl) EditUser(ctx *gin.Context) {
	
    id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}
	userId, ok := ctx.Get(middleware.CLAIM_USER_ID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "invalid user session"})
		return
	}
	userIdInt, ok := userId.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid user id session"})
		return
	}
	if id != int(userIdInt) {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "invalid user request"})
		return
	}
    // Parse user data from request body
    var user model.UserUpdateInput
    if err := ctx.ShouldBindJSON(&user); err != nil {
        ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "Invalid request body"})
        return
    }

    // Call service to edit user data
    updatedUser, err := u.svc.EditUser(ctx, uint64(id), user)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }

    // Return updated user data
    ctx.JSON(http.StatusOK, updatedUser)
}

// DeleteUsersById godoc
//
//		@Summary		Delete user by selected id
//		@Description	will delete user with given id from param
//		@Tags			users
//		@Accept			json
//		@Produce		json
// 		@Security Bearer
// 		@Param 			id path int true "User ID"
//		@Success		200	{object}	model.User
//		@Failure		400	{object}	pkg.ErrorResponse
//		@Failure		404	{object}	pkg.ErrorResponse
//		@Failure		500	{object}	pkg.ErrorResponse
//		@Router			/users/{id} [delete]
func (u *userHandlerImpl) DeleteUsersById(ctx *gin.Context) {
	// get id user
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}

	// check user id session from context
	userId, ok := ctx.Get(middleware.CLAIM_USER_ID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "invalid user session"})
		return
	}
	userIdInt, ok := userId.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid user id session"})
		return
	}
	if id != int(userIdInt) {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "invalid user request"})
		return
	}

	user, err := u.svc.DeleteUsersById(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	if user.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "user not found"})
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{
		"user":    user,
		"message": "Your account has been successfully deleted",
	})
}
// UserSignUp godoc
// @Summary User register
// @Description Register a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.UserSignUp true "User sign-up details"
// @Success 201 {object} model.User "user"
// @Failure 400 {object} pkg.ErrorResponse "Bad request"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /users/register [post]
func (u *userHandlerImpl) UserSignUp(ctx *gin.Context) {
	// binding sign-up body
	userSignUp := model.UserSignUp{}
	if err := ctx.Bind(&userSignUp); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	if err := userSignUp.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	user, err := u.svc.SignUp(ctx, userSignUp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, map[string]any{
		"user":  user,
	})
}
// UserLogin godoc
// @Summary User login
// @Description Log in a user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param userLogin body model.UserLogin true "User login details"
// @Success 200 {string} token "token"
// @Failure 400 {object} pkg.ErrorResponse "Bad request"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /users/login [post]
func (u *userHandlerImpl) UserLogin(ctx *gin.Context) {
    var userLogin model.UserLogin

    // Parsing data dari body permintaan ke struct UserLogin
    if err := ctx.Bind(&userLogin); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

    // Memeriksa kredensial pengguna
    user, err := u.svc.CheckCredentials(ctx, userLogin.Email, userLogin.Password)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }

    // Menghasilkan token akses untuk pengguna yang berhasil login
    token, err := u.svc.GenerateUserAccessToken(ctx, user)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }

    // Mengirimkan token akses sebagai respons ke klien
    ctx.JSON(http.StatusOK, gin.H{"token": token})
}
