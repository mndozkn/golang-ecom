package http

import (
	"go-crud/internal/delivery/http/middleware"
	"go-crud/internal/repository"
	"go-crud/internal/service"
	"go-crud/internal/worker"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, jwtSecret string, db *gorm.DB, redisClient *redis.Client, distributor worker.TaskDistributor) {
	api := app.Group("/api/v1")

	// Repositories
	catRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Services
	catService := service.NewCategoryService(catRepo)
	productService := service.NewProductService(productRepo)
	authService := service.NewAuthService(userRepo, redisClient, jwtSecret, distributor)
	orderService := service.NewOrderService(orderRepo, productRepo)
	userService := service.NewUserService(userRepo)

	// Handler Init
	productHandler := &ProductHandler{Service: productService}
	authHandler := &AuthHandler{Service: authService}
	orderHandler := &OrderHandler{Service: orderService}
	categoryHandler := &CategoryHandler{Service: catService}
	userHandler := &UserHandler{Service: userService}
	adminHandler := &AdminHandler{}

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)

	api.Get("/products", productHandler.GetAll)
	api.Get("/categories", categoryHandler.GetAll)

	userArea := api.Group("/user", middleware.Protected(jwtSecret))
	userArea.Get("/profile", func(c *fiber.Ctx) error { return c.SendString("Profil") })
	userArea.Get("/orders", orderHandler.GetMyOrders)
	userArea.Post("/orders/:id/cancel", orderHandler.CancelOrder)
	userArea.Put("/change-password", userHandler.ChangePassword)

	sellerArea := api.Group("/seller", middleware.Protected(jwtSecret), middleware.RoleCheck("seller"))
	sellerArea.Post("/products", productHandler.Create)
	sellerArea.Get("/my-products", productHandler.GetSellerProducts)
	sellerArea.Get("/orders", orderHandler.GetSellerOrders)
	sellerArea.Get("/dashboard", orderHandler.GetSellerDashboard)
	sellerArea.Patch("/orders/:id/status", orderHandler.UpdateStatus)

	adminArea := api.Group("/admin", middleware.Protected(jwtSecret), middleware.RoleCheck("admin"))
	adminArea.Delete("/users/:id", userHandler.DeleteUser)
	adminArea.Get("/stats", adminHandler.GetSystemStats)
	adminArea.Post("/categories", categoryHandler.Create)
	adminArea.Delete("/categories/:id", categoryHandler.Delete)

	app.Get("/swagger/*", swagger.HandlerDefault)
}
