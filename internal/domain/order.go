package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// @Description Sipariş durumları: 0:Pending, 1:Shipped, 2:Delivered, 3:Cancelled
type OrderStatus int

const (
	StatusPending   OrderStatus = iota // 0: Beklemede
	StatusShipped                      // 1: Kargolandı
	StatusDelivered                    // 2: Teslim Edildi
	StatusCancelled                    // 3: İptal Edildi
)

func (s OrderStatus) String() string {
	return [...]string{"Pending", "Shipped", "Delivered", "Cancelled"}[s]
}

type Order struct {
	ID          uint           `gorm:"primaryKey" json:"id" example:"100"`
	UserID      uint           `json:"user_id" example:"5"`
	TotalAmount float64        `json:"total_amount" example:"1500.50"`
	Status      OrderStatus    `json:"status" gorm:"default:0" example:"0" enums:"0,1,2,3"`
	Items       []OrderItem    `json:"items"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OrderID   uint           `json:"order_id"`
	ProductID uint           `json:"product_id"`
	Quantity  int            `json:"quantity"`
	Price     float64        `json:"price"`
	SellerID  uint           `json:"seller_id"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByUserID(ctx context.Context, userID uint) ([]Order, error)
	GetBySellerID(ctx context.Context, sellerID uint) ([]OrderItem, error)
	GetSellerStats(ctx context.Context, sellerID uint) (map[string]interface{}, error)
	UpdateStatus(ctx context.Context, orderID uint, status OrderStatus) error
	GetByID(ctx context.Context, id uint) (*Order, error)
	CancelOrderWithStockRestore(ctx context.Context, orderID uint) error
}

type OrderService interface {
	PlaceOrder(ctx context.Context, userID uint, items []OrderItem) (*Order, error)
	GetOrdersByUserID(ctx context.Context, userID uint) ([]Order, error)
	GetSellerStats(ctx context.Context, sellerID uint) (map[string]interface{}, error)
	UpdateStatus(ctx context.Context, orderID uint, status OrderStatus) error
	CancelOrder(ctx context.Context, orderID uint, userID uint, role string) error
}
