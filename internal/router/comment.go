package router

import (
	"final_project/internal/handler"
	"final_project/internal/middleware"

	"github.com/gin-gonic/gin"
)

type CommentRouter interface {
	Mount()
}

type commentRouterImpl struct {
	v       *gin.RouterGroup
	handler handler.CommentHandler
}

func NewCommentRouter(v *gin.RouterGroup, handler handler.CommentHandler) CommentRouter {
	return &commentRouterImpl{v: v, handler: handler}
}

func (c *commentRouterImpl) Mount() {
	c.v.Use(middleware.CheckAuthBearer)
	c.v.GET("/:id", c.handler.GetCommentByID)
	c.v.GET("", c.handler.GetCommentsByPhotoID)
	c.v.POST("", c.handler.CreateComment)
	// c.v.GET("", c.handler.GetComments)
	c.v.PUT("/:id", c.handler.UpdateComment)
	c.v.DELETE("/:id", c.handler.DeleteComment)
}
