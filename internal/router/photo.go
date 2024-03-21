package router

import (
	"final_project/internal/handler"
	"final_project/internal/middleware"

	"github.com/gin-gonic/gin"
)

type PhotoRouter interface {
	Mount()
}

type photoRouterImpl struct {
	v       *gin.RouterGroup
	handler handler.PhotoHandler
}

func NewPhotoRouter(v *gin.RouterGroup, handler handler.PhotoHandler) PhotoRouter {
	return &photoRouterImpl{v: v, handler: handler}
}

func (p *photoRouterImpl) Mount() {
	// Endpoint untuk menambahkan middleware autentikasi ke setiap permintaan foto
	p.v.Use(middleware.CheckAuthBearer)

	// Endpoint untuk menampilkan daftar foto
	p.v.GET("", p.handler.GetPhotos)

	// Endpoint untuk menampilkan detail foto berdasarkan ID
	p.v.GET("/:id", p.handler.GetPhotoByID)
	p.v.GET("/user_id/:id", p.handler.GetPhotoByUserID)

	// Endpoint untuk menghapus foto berdasarkan ID
	p.v.DELETE("/:id", p.handler.DeletePhotoByID)
	p.v.PUT("/:id", p.handler.EditPhoto)

	// Endpoint untuk membuat foto baru
	p.v.POST("", p.handler.CreatePhoto)
}
