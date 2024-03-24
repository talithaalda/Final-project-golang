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

type CommentHandler interface {
    GetCommentsByPhotoID(ctx *gin.Context)
    CreateComment(ctx *gin.Context)
    UpdateComment(ctx *gin.Context)
    DeleteComment(ctx *gin.Context)
	GetCommentByID(ctx *gin.Context)
	GetComments(ctx *gin.Context)
}

type commentHandlerImpl struct {
    commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) CommentHandler {
    return &commentHandlerImpl{commentService: commentService}
}
// GetCommentsByPhotoID retrieves comments by photo ID.
// @Summary Retrieve comments by photo ID
// @Description Retrieve a list of comments based on the given photo ID
// @Tags comments
// @Accept json
// @security Bearer
// @Produce json
// @Param photo_id query int true "Photo ID"
// @Success 200 {array} model.Comment "comments"
// @Success 200 {object} pkg.ErrorResponse "No comment found"
// @Failure 400 {object} pkg.ErrorResponse "Bad request"
// @Failure 404 {object} pkg.ErrorResponse "Comment not found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /comments [get]
func (s *commentHandlerImpl) GetCommentsByPhotoID(ctx *gin.Context) {
    // Ambil nilai photo_id dari query string
    photoIDStr := ctx.Query("photo_id")
    
    // Periksa apakah nilai photo_id ada atau tidak
    if photoIDStr == "" {
        ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "Photo ID is required"})
        return
    }

    // Parse nilai photo_id menjadi uint64
    photoID, err := strconv.ParseUint(photoIDStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "Invalid photo ID"})
        return
    }

    // Panggil service untuk mendapatkan daftar comment berdasarkan photo ID
    comments, err := s.commentService.GetCommentsByPhotoID(ctx, photoID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }

    // Periksa jika tidak ada comment yang ditemukan
    if len(comments) == 0 {
        ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Comment not found"})
        return
    }

    // Kirim respons dengan daftar comment yang ditemukan
    ctx.JSON(http.StatusOK, comments)
}
// CreateComment creates a new comment.
// @Summary Create a new comment
// @Description Create a new comment for a photo
// @Tags comments
// @Accept json
// @Produce json
// @Security Bearer
// @Param photo_id body model.CreateCommentInput true "Comment data"
// @Success 201 {object} model.Comment "comment"
// @Failure 400 {object} pkg.ErrorResponse "Bad request"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /comments [post]
func (c *commentHandlerImpl) CreateComment(ctx *gin.Context) {
    comment := model.CreateCommentInput{}
    if err := ctx.ShouldBindJSON(&comment); err != nil {
        ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "Invalid request body"})
        return
    }
	
    createdComment, err := c.commentService.CreateComment(ctx, comment, uint64(ctx.MustGet(middleware.CLAIM_USER_ID).(float64)))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, createdComment)
}
// UpdateComment updates an existing comment.
// @Summary Update an existing comment
// @Description Update an existing comment
// @Tags comments
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Comment ID"
// @Param comment body model.UpdateCommentInput true "Comment data"
// @Success 200 {object} model.Comment "updatedComment"
// @Failure 400 {object} pkg.ErrorResponse "Bad request"
// @Failure 401 {object} pkg.ErrorResponse "Unauthorized"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /comments/{id} [put]
func (c *commentHandlerImpl) UpdateComment(ctx *gin.Context) {
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
	comment, err := c.commentService.GetCommentByID1(ctx, uint64(id))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }

    // Check if the photo exists
    if comment.ID == 0 {
        ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "Photo not found"})
        return
    }
	userIdInt, ok := userId.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid user id session"})
		return
	}
    // Check if the user ID from the middleware matches the user ID in the photo data
    if int(comment.UserID) != int(userIdInt) {
        ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "You are not authorized to edit this photo"})
        return
    }
    // Parse photo data from request body
    
    if err := ctx.ShouldBindJSON(&comment); err != nil {
        ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "Invalid request body"})
        return
    }
	inputComment := model.UpdateCommentInput{
		Message: comment.Message,
	}
    // Call service to edit photo data
    updatedComment, err := c.commentService.UpdateComment(ctx, uint64(id), inputComment)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }

    // Return updated photo data
    ctx.JSON(http.StatusOK, updatedComment)
}
// DeleteComment deletes an existing comment.
// @Summary Delete an existing comment
// @Description Delete an existing comment
// @Tags comments
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Comment ID"
// @Success 200 {object} model.Comment "comment"
// @Failure 400 {object} pkg.ErrorResponse "Bad request"
// @Failure 401 {object} pkg.ErrorResponse "Unauthorized"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /comments/{id} [delete]
func (c *commentHandlerImpl) DeleteComment(ctx *gin.Context) {
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

	// Delete comment by ID
	comment, err := c.commentService.DeleteCommentByID(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}

	// Check if the comment exists
	if comment.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "comment not found"})
		return
	}

	// Check if the user ID from the middleware matches the user ID in the comment data
	if int(comment.UserID) != int(userIdInt) {
        ctx.JSON(http.StatusUnauthorized, pkg.ErrorResponse{Message: "You are not authorized to delete this comment"})
        return
    }

	ctx.JSON(http.StatusOK, map[string]any{
		"comment":    comment,
		"message": "Your comment has been successfully deleted",
	})
}
// GetCommentByID retrieves a comment by its ID.
// @Summary Retrieve comment by ID
// @Description Retrieve a comment by its ID
// @Tags comments
// @Accept json
// @security Bearer
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} model.Comment "comment"
// @Failure 400 {object} pkg.ErrorResponse "Bad request"
// @Failure 404 {object} pkg.ErrorResponse "Comment not found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /comments/{id} [get]
func (c *commentHandlerImpl) GetCommentByID(ctx *gin.Context) {
	// get photo ID from path parameter
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "invalid photo ID"})
		return
	}

	comment, err := c.commentService.GetCommentByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	if comment.ID == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "comment not found"})
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

func (c *commentHandlerImpl) GetComments(ctx *gin.Context) {
	comments, err := c.commentService.GetComments(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
		return
	}
	if len(comments) == 0 {
        ctx.JSON(http.StatusOK, gin.H{"message": "No comment found"})
        return
    }
	ctx.JSON(http.StatusOK, comments)
}