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

func (c *commentHandlerImpl) GetCommentsByPhotoID(ctx *gin.Context) {
    photoID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, pkg.ErrorResponse{Message: "Invalid comment ID"})
        return
    }

    comments, err := c.commentService.GetCommentsByPhotoID(ctx, photoID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }
	if len(comments) == 0 {
		ctx.JSON(http.StatusNotFound, pkg.ErrorResponse{Message: "comment not found"})
		return
	}
    ctx.JSON(http.StatusOK, comments)
}

func (c *commentHandlerImpl) CreateComment(ctx *gin.Context) {
    comment := model.CreateComment{}
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

    // Call service to edit photo data
    updatedComment, err := c.commentService.UpdateComment(ctx, uint64(id), comment)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{Message: err.Error()})
        return
    }

    // Return updated photo data
    ctx.JSON(http.StatusOK, updatedComment)
}

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
	ctx.JSON(http.StatusOK, comments)
}