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
	p.v.Use(middleware.CheckAuthBearer)
	// p.v.GET("", p.handler.GetPhotos)
	p.v.GET("/:id", p.handler.GetPhotoByID)
	p.v.GET("", p.handler.GetPhotoByUserID)
	p.v.DELETE("/:id", p.handler.DeletePhotoByID)
	p.v.PUT("/:id", p.handler.EditPhoto)
	p.v.POST("", p.handler.CreatePhoto)
}
