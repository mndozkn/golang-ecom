package http

import (
	"go-crud/internal/delivery/http/middleware"
	"go-crud/internal/domain"
	
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App, ps domain.ProductService, as domain.AuthService, os domain.OrderService, cs domain.CategoryService, jwtSecret string) {
	api := app.Group("/api/v1")

	// Handler Init
	productHandler := &ProductHandler{Service: ps}
	authHandler := &AuthHandler{Service: as}
	orderHandler := &OrderHandler{Service: os}
	categoryHandler := &CategoryHandler{Service: cs}
	userHandler := &UserHandler{}
	adminHandler := &AdminHandler{}

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	api.Get("/products", productHandler.GetAll)
	api.Get("/categories", categoryHandler.GetAll)

	userArea := api.Group("/user", middleware.Protected(jwtSecret))
	userArea.Get("/profile", func(c *fiber.Ctx) error { return c.SendString("Profil") })
	userArea.Get("/orders", orderHandler.GetMyOrders)
	userArea.Post("/orders/:id/cancel", orderHandler.CancelOrder)

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
