package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"final_project/internal/handler"
	"final_project/internal/infrastructure"
	"final_project/internal/model"
	"final_project/internal/repository"
	"final_project/internal/router"
	"final_project/internal/service"
	"final_project/pkg"
	"final_project/pkg/helper"

	"github.com/gin-gonic/gin"

	_ "final_project/docs"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			GO DTS USER API DUCUMENTATION
// @version		2.0
// @description	Final Project My Gram
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.url	http://www.swagger.io/support
// @contact.email	support@swagger.io
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:3000
// @BasePath		/
// @schemes		http
// @SecurityDefinitions.apiKey Bearer
// @in header
// @name Authorization
func main() {
	// requirement technical:
	// [x] middleware untuk recover ketika panic
	// [x] mengecheck basic auth
	server()
}

// Product:
// authorization menggunakan jwt
// authentication bisa dilakukan dengan login
// ketika user login, akan memunculkan JWT ketika success

func server() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	g := gin.Default()
	g.Use(gin.Recovery())

	// /public => generate JWT public
	g.GET("/public", func(ctx *gin.Context) {
		now := time.Now()

		claim := model.StandardClaim{
			Jti: fmt.Sprintf("%v", time.Now().UnixNano()),
			Iss: "go-middleware",
			Aud: "golang-006",
			Sub: "public-token",
			Exp: uint64(now.Add(time.Hour).Unix()),
			Iat: uint64(now.Unix()),
			Nbf: uint64(now.Unix()),
		}
		token, err := helper.GenerateToken(claim)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
				Message: "error generating public token",
				Errors:  []string{err.Error()},
			})
			return
		}
		ctx.JSON(http.StatusOK, map[string]any{"token": token})
	})

	usersGroup := g.Group("/users")
	// usersGroup.Use(middleware.CheckAuthBasic)
	// usersGroup.Use(middleware.CheckAuthBearer)

	// dependency injection
	// dig by uber
	// wire

	// https://s8sg.medium.com/solid-principle-in-go-e1a624290346
	gorm := infrastructure.NewGormPostgres()
	userRepo := repository.NewUserQuery(gorm)
	userSvc := service.NewUserService(userRepo)
	userHdl := handler.NewUserHandler(userSvc)
	userRouter := router.NewUserRouter(usersGroup, userHdl)
	userRouter.Mount()
	photosGroup := g.Group("/photos")
	photoRepo := repository.NewPhotoQuery(gorm) 
	photoSvc := service.NewPhotoService(photoRepo, userRepo)   
	photoHdl := handler.NewPhotoHandler(photoSvc)   
	photoRouter := router.NewPhotoRouter(photosGroup, photoHdl)
	photoRouter.Mount()
	commentsGroup := g.Group("/comments")
	commentRepo := repository.NewCommentQuery(gorm) 
	commentSvc := service.NewCommentService(commentRepo, userRepo, photoRepo)   
	commentHdl := handler.NewCommentHandler(commentSvc)   
	commentRouter := router.NewCommentRouter(commentsGroup, commentHdl)
	commentRouter.Mount()
	socialMediasGroup := g.Group("/socialmedias")
	socialMediaRepo := repository.NewSocialMediaQuery(gorm) 
	socialMediaSvc := service.NewSocialMediaService(socialMediaRepo, userRepo)   
	socialMediaHdl := handler.NewSocialMediaHandler(socialMediaSvc)   
	socialMediaRouter := router.NewSocialMediaRouter(socialMediasGroup, socialMediaHdl)
	socialMediaRouter.Mount()

	// mount
	
	// swagger
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	g.Run(":3000")
}

