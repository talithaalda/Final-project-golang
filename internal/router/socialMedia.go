package router

import (
	"final_project/internal/handler"
	"final_project/internal/middleware"

	"github.com/gin-gonic/gin"
)

type SocialMediaRouter interface {
	Mount()
}

type socialMediaRouterImpl struct {
	v       *gin.RouterGroup
	handler handler.SocialMediaHandler
}

func NewSocialMediaRouter(v *gin.RouterGroup, handler handler.SocialMediaHandler) SocialMediaRouter {
	return &socialMediaRouterImpl{v: v, handler: handler}
}

func (c *socialMediaRouterImpl) Mount() {
	// Menambahkan middleware autentikasi ke setiap permintaan komentar
	c.v.Use(middleware.CheckAuthBearer)

	// Menampilkan daftar komentar berdasarkan ID foto
	c.v.GET("/:id", c.handler.GetSocialMediaByID)
	c.v.GET("", c.handler.GetSocialMediasByUserID)
	// Membuat komentar baru
	c.v.POST("", c.handler.CreateSocialMedia)
	// c.v.GET("", c.handler.GetSocialMedias)

	// Memperbarui komentar berdasarkan ID
	c.v.PUT("/:id", c.handler.UpdateSocialMedia)

	// Menghapus komentar berdasarkan ID
	c.v.DELETE("/:id", c.handler.DeleteSocialMedia)
}
