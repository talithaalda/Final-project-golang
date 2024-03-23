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

type SocialMediaHandler interface {
    GetSocialMediasByUserID(ctx *gin.Context)
    CreateSocialMedia(ctx *gin.Context)
    UpdateSocialMedia(ctx *gin.Context)
    DeleteSocialMedia(ctx *gin.Context)
	GetSocialMediaByID(ctx *gin.Context)
	GetSocialMedias(ctx *gin.Context)
}

type socialMediaHandlerImpl struct {
    socialMediaService service.SocialMediaService
}

func NewSocialMediaHandler(socialMediaService service.SocialMediaService) SocialMediaHandler {
    return &socialMediaHandlerImpl{socialMediaService: socialMediaService}
}

func (s *socialMediaHandlerImpl) GetSocialMediasByUserID(ctx *gin.Context) {
    // Ambil nilai user_id dari query string
    userIDStr := ctx.Query("user_id")
    
    // Periksa apakah nilai user_id ada atau tidak
    if userIDStr == "" {
        ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "User ID is required"})
        return
    }

    // Parse nilai user_id menjadi uint64
    userID, err := strconv.ParseUint(userIDStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "Invalid user ID"})
        return
    }

    // Panggil service untuk mendapatkan daftar social media berdasarkan user ID
    socialMedias, err := s.socialMediaService.GetSocialMediasByUserID(ctx, userID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }

    // Periksa jika tidak ada social media yang ditemukan
    if len(socialMedias) == 0 {
        ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Social media not found"})
        return
    }

    // Kirim respons dengan daftar social media yang ditemukan
    ctx.JSON(http.StatusOK, socialMedias)
}


func (s *socialMediaHandlerImpl) CreateSocialMedia(ctx *gin.Context) {
    socialMedia := model.CreateSocialMedia{}
    if err := ctx.ShouldBindJSON(&socialMedia); err != nil {
        ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "Invalid request body"})
        return
    }

    createdSocialMedia, err := s.socialMediaService.CreateSocialMedia(ctx, socialMedia, uint64(ctx.MustGet(middleware.CLAIM_USER_ID).(float64)))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, createdSocialMedia)
}

func (s *socialMediaHandlerImpl) UpdateSocialMedia(ctx *gin.Context) {
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
	socialMedia, err := s.socialMediaService.GetSocialMediaByID1(ctx, uint64(id))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }

    // Check if the social media exists
    if socialMedia.ID == 0 {
        ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Social media not found"})
        return
    }
	userIdInt, ok := userId.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid user id session"})
		return
	}
    // Check if the user ID from the middleware matches the user ID in the social media data
    if int(socialMedia.UserID) != int(userIdInt) {
        ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "You are not authorized to edit this social media"})
        return
    }
    // Parse social media data from request body
    
    if err := ctx.ShouldBindJSON(&socialMedia); err != nil {
        ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "Invalid request body"})
        return
    }

    // Call service to edit social media data
    updatedSocialMedia, err := s.socialMediaService.UpdateSocialMedia(ctx, uint64(id), socialMedia)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }

    // Return updated social media data
    ctx.JSON(http.StatusOK, updatedSocialMedia)
}

func (s *socialMediaHandlerImpl) DeleteSocialMedia(ctx *gin.Context) {
    id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid required param"})
		return
	}

	// Check user ID session from context
	userId, ok := ctx.Get(middleware.CLAIM_USER_ID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "Invalid user session"})
		return
	}
	userIdInt, ok := userId.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "Invalid user ID session"})
		return
	}

	// Delete socialMedia by ID
	socialMedia, err := s.socialMediaService.DeleteSocialMediaByID(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	// Check if the socialMedia exists
	if socialMedia.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "socialMedia not found"})
		return
	}

	// Check if the user ID from the middleware matches the user ID in the socialMedia data
	if int(socialMedia.UserID) != int(userIdInt) {
        ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "You are not authorized to delete this socialMedia"})
        return
    }

	ctx.JSON(http.StatusOK, map[string]any{
		"socialMedia":    socialMedia,
		"message": "Your social media has been successfully deleted",
	})
}
func (s *socialMediaHandlerImpl) GetSocialMediaByID(ctx *gin.Context) {
	// get social media ID from path parameter
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid social media ID"})
		return
	}

	socialMedia, err := s.socialMediaService.GetSocialMediaByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	if socialMedia.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "social media not found"})
		return
	}

	ctx.JSON(http.StatusOK, socialMedia)
}
func (s *socialMediaHandlerImpl) GetSocialMedias(ctx *gin.Context) {
	socialMedias, err := s.socialMediaService.GetSocialMedias(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	if len(socialMedias) == 0 {
        ctx.JSON(http.StatusOK, gin.H{"message": "No social media found"})
        return
    }
	ctx.JSON(http.StatusOK, socialMedias)
}