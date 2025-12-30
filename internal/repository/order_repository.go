package repository

import (
	"context"
	"go-crud/internal/domain"
	"go-crud/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, item := range order.Items {
			var product domain.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				return utils.NewError(fiber.StatusNotFound, "Ürün bulunamadı")
			}

			if product.Stock < item.Quantity {
				return utils.NewError(fiber.StatusBadRequest, product.Name+" için yetersiz stok")
			}

			if err := tx.Model(&product).Update("stock", product.Stock-item.Quantity).Error; err != nil {
				return err
			}

			order.Items[i].SellerID = product.SellerID
			order.Items[i].Price = product.Price
			order.TotalAmount += product.Price * float64(item.Quantity)
		}

		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Order, error) {
	var orders []domain.Order

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Items").
		Preload("Items.Product").
		Order("created_at DESC").
		Find(&orders).Error

	return orders, err
}

func (r *orderRepository) GetBySellerID(ctx context.Context, sellerID uint) ([]domain.OrderItem, error) {
	var items []domain.OrderItem
	err := r.db.WithContext(ctx).Where("seller_id = ?", sellerID).Preload("Product").Find(&items).Error
	return items, err
}

func (r *orderRepository) GetSellerStats(ctx context.Context, sellerID uint) (map[string]interface{}, error) {
	var result struct {
		TotalRevenue float64 `json:"total_revenue"`
		TotalOrders  int64   `json:"total_orders"`
		TotalItems   int64   `json:"total_items"`
	}

	err := r.db.WithContext(ctx).
		Model(&domain.OrderItem{}).
		Select("SUM(price * quantity) as total_revenue, COUNT(DISTINCT order_id) as total_orders, SUM(quantity) as total_items").
		Where("seller_id = ?", sellerID).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_revenue": result.TotalRevenue,
		"total_orders":  result.TotalOrders,
		"total_items":   result.TotalItems,
	}, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, orderID uint, status domain.OrderStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *orderRepository) GetByID(ctx context.Context, orderID uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).First(&order, orderID).Error
	return &order, err
}

func (r *orderRepository) CancelOrderWithStockRestore(ctx context.Context, orderID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order domain.Order
		if err := tx.Preload("Items").First(&order, orderID).Error; err != nil {
			return err
		}

		for _, item := range order.Items {
			if err := tx.Model(&domain.Product{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		return tx.Model(&order).Update("status", domain.StatusCancelled).Error
	})
}
