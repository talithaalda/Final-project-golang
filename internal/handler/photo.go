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
	GetPhotoByUserID(ctx *gin.Context)
}

type photoHandlerImpl struct {
	photoService service.PhotoService
}

func NewPhotoHandler(photoService service.PhotoService) PhotoHandler {
	return &photoHandlerImpl{photoService: photoService}
}
// GetPhotos godoc
// @Summary Retrieve list of photos
// @Description Retrieve a list of all photos.
// @Tags photos
// @Accept json
// @Produce json
// @Success	200	{object} model.Photo
// @Success 200 {object} pkg.ErrorResponse "No photo found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /photos [get]
func (p *photoHandlerImpl) GetPhotos(ctx *gin.Context) {
	photos, err := p.photoService.GetPhotos(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	if len(photos) == 0 {
        ctx.JSON(http.StatusOK, gin.H{"message": "No photo found"})
        return
    }
	ctx.JSON(http.StatusOK, photos)
}
// GetPhotoByID godoc
// @Summary Retrieve photo by ID
// @Description Retrieve a photo by its ID
// @Tags photos
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Photo ID"
// @Success 200 {object} model.UpdatePhoto "photo"
// @Failure 400 {object} pkg.ErrorResponse "Bad request"
// @Failure 404 {object} pkg.ErrorResponse "Photo not found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /photos/{id} [get]
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
// DeletePhotoByID godoc
// @Summary Delete photo by ID
// @Description Delete a photo by its ID.
// @Tags photos
// @Accept json
// @Produce json
// @Param id path int true "Photo ID"
// @Security Bearer
// @Success	200	{object} model.UpdatePhoto
// @Failure 404 {object} pkg.ErrorResponse "Photo not found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /photos/{id} [delete]
func (p *photoHandlerImpl) DeletePhotoByID(ctx *gin.Context) {
	// Get photo ID from path parameter
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

	// Delete photo by ID
	photo, err := p.photoService.DeletePhotoByID(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	// Check if the photo exists
	if photo.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Photo not found"})
		return
	}

	// Check if the user ID from the middleware matches the user ID in the photo data
	if int(photo.UserID) != int(userIdInt) {
        ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "You are not authorized to delete this photo"})
        return
    }

	ctx.JSON(http.StatusOK, map[string]any{
		"photo":    photo,
		"message": "Your photo has been successfully deleted",
	})
}
// CreatePhoto godoc
// @Summary Create a new photo
// @Description Create a new photo.
// @Tags photos
// @Accept json
// @Produce json
// @Param photo body model.InputPhoto true "Photo data"
// @Security Bearer
// @Success	200	{object} model.CreatePhoto
// @Failure 400 {object} pkg.ErrorResponse "Bad request"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /photos [post]
func (p *photoHandlerImpl) CreatePhoto(ctx *gin.Context) {
	photo := model.InputPhoto{}
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
// EditPhoto godoc
// @Summary Update photo information
// @Description Update information of a photo.
// @Tags photos
// @Accept json
// @Produce json
// @Param id path int true "Photo ID"
// @Param photo body model.InputPhoto true "Photo data"
// @Security Bearer
// @Success	200	{object} model.UpdatePhoto
// @Failure 400 {object} pkg.ErrorResponse "Bad request"
// @Failure 401 {object} pkg.ErrorResponse "Unauthorized"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /photos/{id} [put]
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
	inputPhoto := model.InputPhoto{
		Title: photo.Title,
		Caption: photo.Caption,
		PhotoURL: photo.PhotoURL,
	}
    // Call service to edit photo data
    updatedPhoto, err := p.photoService.EditPhoto(ctx, uint64(id), inputPhoto)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }

    // Return updated photo data
    ctx.JSON(http.StatusOK, updatedPhoto)
}
// GetPhotoByUserID godoc
// @Summary Retrieve photos by user ID
// @Description Retrieve a list of photos by user ID.
// @Tags photos
// @Accept json
// @Produce json
// @Security Bearer
// @Param user_id query int true "User ID"
// @Success	200	{object} model.Photo
// @Failure 400 {object} pkg.ErrorResponse "Bad request"
// @Failure 404 {object} pkg.ErrorResponse "Photo not found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /photos [get]
func (s *photoHandlerImpl) GetPhotoByUserID(ctx *gin.Context) {
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

    // Panggil service untuk mendapatkan daftar photo berdasarkan user ID
    photos, err := s.photoService.GetPhotoByUserID(ctx, userID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }

    // Periksa jika tidak ada photo yang ditemukan
    if len(photos) == 0 {
        ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Photo not found"})
        return
    }

    // Kirim respons dengan daftar photo yang ditemukan
    ctx.JSON(http.StatusOK, photos)
}