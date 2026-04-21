package main

import (
	"log"
	"os"

	_ "github.com/chillman2101/gits-catalogue/docs"
	"github.com/chillman2101/gits-catalogue/internal/config"
	"github.com/chillman2101/gits-catalogue/internal/handler"
	"github.com/chillman2101/gits-catalogue/internal/repository"
	"github.com/chillman2101/gits-catalogue/internal/router"
	"github.com/chillman2101/gits-catalogue/internal/service"
	"github.com/chillman2101/gits-catalogue/pkg/redis"
	"github.com/chillman2101/gits-catalogue/pkg/redis/cache"
	"github.com/joho/godotenv"
)

// @title           Gits Catalogue API
// @version         1.0
// @description     RESTful API for book catalogue management
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	cfg := config.Load()

	db, err := cfg.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	var cacheHelper *cache.CacheHelper
	redisClient, err := redis.NewClient(cfg.RedisAddr, cfg.RedisDB)
	if err != nil {
		log.Printf("warning: redis unavailable, running without cache: %v", err)
		cacheHelper = cache.NewCacheHelper(nil)
	} else {
		cacheHelper = cache.NewCacheHelper(redisClient.GetClient())
	}

	userRepo := repository.NewUserRepository(db)
	authorRepo := repository.NewAuthorRepository(db)
	publisherRepo := repository.NewPublisherRepository(db)
	bookRepo := repository.NewBookRepository(db)

	authSvc := service.NewAuthService(userRepo)
	authorSvc := service.NewAuthorService(authorRepo, cacheHelper)
	publisherSvc := service.NewPublisherService(publisherRepo, cacheHelper)
	bookSvc := service.NewBookService(bookRepo, cacheHelper)

	h := router.Handlers{
		Auth:      handler.NewAuthHandler(authSvc),
		Author:    handler.NewAuthorHandler(authorSvc),
		Publisher: handler.NewPublisherHandler(publisherSvc),
		Book:      handler.NewBookHandler(bookSvc),
	}

	r := router.Setup(h, userRepo)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("server running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
