package router

import (
	"os"

	"github.com/chillman2101/gits-catalogue/internal/handler"
	"github.com/chillman2101/gits-catalogue/internal/middleware"
	"github.com/chillman2101/gits-catalogue/internal/repository"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	Auth      *handler.AuthHandler
	Author    *handler.AuthorHandler
	Publisher *handler.PublisherHandler
	Book      *handler.BookHandler
}

func Setup(h Handlers, userRepo repository.UserRepository) *gin.Engine {
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := gin.Default()
	jwtAuth := middleware.JWTAuth(userRepo)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
		auth.POST("/refresh", h.Auth.Refresh)
		auth.POST("/logout", jwtAuth, h.Auth.Logout)
	}

	api := r.Group("/api/v1")
	api.Use(jwtAuth)
	{
		authors := api.Group("/authors")
		{
			authors.GET("", h.Author.GetAll)
			authors.GET("/:id", h.Author.GetByID)
			authors.POST("", h.Author.Create)
			authors.PUT("/:id", h.Author.Update)
			authors.DELETE("/:id", h.Author.Delete)
		}

		publishers := api.Group("/publishers")
		{
			publishers.GET("", h.Publisher.GetAll)
			publishers.GET("/:id", h.Publisher.GetByID)
			publishers.POST("", h.Publisher.Create)
			publishers.PUT("/:id", h.Publisher.Update)
			publishers.DELETE("/:id", h.Publisher.Delete)
		}

		books := api.Group("/books")
		{
			books.GET("", h.Book.GetAll)
			books.GET("/:id", h.Book.GetByID)
			books.POST("", h.Book.Create)
			books.PUT("/:id", h.Book.Update)
			books.DELETE("/:id", h.Book.Delete)
		}
	}

	return r
}
