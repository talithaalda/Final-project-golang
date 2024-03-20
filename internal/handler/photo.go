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

type PhotoHandler interface {
	GetPhotos(ctx *gin.Context)
	GetPhotoByID(ctx *gin.Context)
	DeletePhotoByID(ctx *gin.Context)
	CreatePhoto(ctx *gin.Context)
	EditPhoto(ctx *gin.Context)
}

type photoHandlerImpl struct {
	photoService service.PhotoService
}

func NewPhotoHandler(photoService service.PhotoService) PhotoHandler {
	return &photoHandlerImpl{photoService: photoService}
}

func (p *photoHandlerImpl) GetPhotos(ctx *gin.Context) {
	photos, err := p.photoService.GetPhotos(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, photos)
}

func (p *photoHandlerImpl) GetPhotoByID(ctx *gin.Context) {
	// get photo ID from path parameter
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid photo ID"})
		return
	}

	photo, err := p.photoService.GetPhotoByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	if photo.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "photo not found"})
		return
	}

	ctx.JSON(http.StatusOK, photo)
}

func (p *photoHandlerImpl) DeletePhotoByID(ctx *gin.Context) {
	// get photo ID from path parameter
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid photo ID"})
		return
	}

	photo, err := p.photoService.DeletePhotoByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	if photo.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "photo not found"})
		return
	}

	ctx.JSON(http.StatusOK, photo)
}

func (p *photoHandlerImpl) CreatePhoto(ctx *gin.Context) {
	var photo model.Photo
	if err := ctx.BindJSON(&photo); err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	createdPhoto, err := p.photoService.CreatePhoto(ctx, photo, uint64(ctx.MustGet(middleware.CLAIM_USER_ID).(float64)))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdPhoto)
}
func (p *photoHandlerImpl) EditPhoto(ctx *gin.Context) {
	
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
	photo, err := p.photoService.GetPhotoByID(ctx, uint64(id))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }

    // Check if the photo exists
    if photo.ID == 0 {
        ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Photo not found"})
        return
    }
	userIdInt, ok := userId.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid user id session"})
		return
	}
    // Check if the user ID from the middleware matches the user ID in the photo data
    if int(photo.UserID) != int(userIdInt) {
        ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "You are not authorized to edit this photo"})
        return
    }
    // Parse photo data from request body
    
    if err := ctx.ShouldBindJSON(&photo); err != nil {
        ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "Invalid request body"})
        return
    }

    // Call service to edit photo data
    updatedPhoto, err := p.photoService.EditPhoto(ctx, uint64(id), photo)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }

    // Return updated photo data
    ctx.JSON(http.StatusOK, updatedPhoto)
}