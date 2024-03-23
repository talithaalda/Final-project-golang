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
	// Menambahkan middleware autentikasi ke setiap permintaan komentar
	c.v.Use(middleware.CheckAuthBearer)

	// Menampilkan daftar komentar berdasarkan ID foto
	c.v.GET("/:id", c.handler.GetCommentByID)
	c.v.GET("", c.handler.GetCommentsByPhotoID)
	// Membuat komentar baru
	c.v.POST("", c.handler.CreateComment)
	// c.v.GET("", c.handler.GetComments)

	// Memperbarui komentar berdasarkan ID
	c.v.PUT("/:id", c.handler.UpdateComment)

	// Menghapus komentar berdasarkan ID
	c.v.DELETE("/:id", c.handler.DeleteComment)
}
