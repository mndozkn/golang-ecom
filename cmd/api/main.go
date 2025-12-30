package main

import (
	"go-crud/internal/delivery/http"
	"go-crud/internal/delivery/http/middleware"
	"go-crud/internal/domain"
	"go-crud/internal/repository"
	"go-crud/internal/service"
	"go-crud/pkg/config"
	"go-crud/pkg/database"
	"log"

	"github.com/gofiber/fiber/v2"
)

// @title Go E-Commerce Backend API
// @version 1.0
// @description Fiber tabanlı profesyonel e-ticaret API dokümantasyonu.
// @termsOfService http://swagger.io/terms/

// @contact.name API Destek
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description JWT tokenınızı şu formatta girin: Bearer <token>
func main() {
	cfg := config.LoadConfig()

	db, err := database.NewPostgresDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	db.AutoMigrate(&domain.Category{}, &domain.Product{}, &domain.User{})

	// Repositories
	catRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Services
	catService := service.NewCategoryService(catRepo)
	productService := service.NewProductService(productRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	orderService := service.NewOrderService(orderRepo, productRepo)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	http.SetupRoutes(app, productService, authService, orderService, catService, cfg.JWTSecret)

	log.Fatal(app.Listen(":" + cfg.ServerPort))
}
