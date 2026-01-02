package main

import (
	"fmt"
	"go-crud/internal/delivery/http"
	"go-crud/internal/delivery/http/middleware"
	"go-crud/internal/domain"
	"go-crud/internal/infrastructure/mailer"
	"go-crud/internal/worker"
	"go-crud/pkg/config"
	"go-crud/pkg/database"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "go-crud/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
)

// @title Go E-Commerce Backend API
// @version 1.0

// @host localhost:3000
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description JWT tokenınızı şu formatta girin: Bearer <token>
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env dosyası yüklenemedi")
	}
	cfg := config.LoadConfig()

	redisClient, err := database.ConnectRedis()
	if err != nil {
		log.Fatalf("Redis başlatılamadı: %v", err)
	}
	defer redisClient.Close()

	redisOpt := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	mailService := mailer.NewMailtrapService()
	distributor := worker.NewRedisTaskDistributor(redisOpt)

	processor := worker.NewTaskProcessor(redisOpt, mailService)
	go func() {
		if err := processor.Start(); err != nil {
			log.Fatal("Worker başlatılamadı")
		}
	}()

	db, err := database.NewPostgresDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	db.AutoMigrate(&domain.Category{}, &domain.Product{}, &domain.User{})

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	http.SetupRoutes(app, cfg.JWTSecret, db, redisClient, distributor)

	log.Fatal(app.Listen(":" + cfg.ServerPort))
}
