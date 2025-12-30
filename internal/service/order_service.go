package service

import (
	"context"
	"go-crud/internal/domain"
	"go-crud/pkg/utils"

	"go.uber.org/zap"
)

type orderService struct {
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
}

func NewOrderService(or domain.OrderRepository, pr domain.ProductRepository) domain.OrderService {
	return &orderService{
		orderRepo:   or,
		productRepo: pr,
	}
}

func (s *orderService) PlaceOrder(ctx context.Context, userID uint, items []domain.OrderItem) (*domain.Order, error) {
	var totalAmount float64

	// Her ürün için stok kontrolü ve fiyat hesaplaması
	for i := range items {
		product, err := s.productRepo.GetByID(ctx, int(items[i].ProductID))
		if err != nil {
			return nil, utils.NewError(404, "Ürün bulunamadı")
		}

		if product.Stock < items[i].Quantity {
			return nil, utils.NewError(400, product.Name+" için yetersiz stok")
		}

		items[i].Price = product.Price
		totalAmount += product.Price * float64(items[i].Quantity)

		items[i].SellerID = product.SellerID
	}

	order := &domain.Order{
		UserID:      userID,
		Items:       items,
		TotalAmount: totalAmount,
		Status:      domain.StatusPending, // Iota Enum: 0
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		utils.Log.Error("Sipariş oluşturma hatası", zap.Error(err))
		return nil, err
	}

	utils.Log.Info("Yeni Sipariş Oluşturuldu",
		zap.Uint("order_id", order.ID),
		zap.Uint("user_id", userID),
	)

	return order, nil
}

// GetOrdersByUserID: Kullanıcının geçmiş siparişlerini getirir (Preload içerir)
func (s *orderService) GetOrdersByUserID(ctx context.Context, userID uint) ([]domain.Order, error) {
	return s.orderRepo.GetByUserID(ctx, userID)
}

func (s *orderService) GetSellerStats(ctx context.Context, sellerID uint) (map[string]interface{}, error) {
	return s.orderRepo.GetSellerStats(ctx, sellerID)
}

func (s *orderService) UpdateStatus(ctx context.Context, orderID uint, status domain.OrderStatus) error {
	err := s.orderRepo.UpdateStatus(ctx, orderID, status)
	if err == nil {
		utils.Log.Info("Sipariş Durumu Güncellendi",
			zap.Uint("order_id", orderID),
			zap.Int("new_status", int(status)),
		)
	}
	return err
}

func (s *orderService) CancelOrder(ctx context.Context, orderID uint, userID uint, role string) error {
	// 1. Siparişi kontrol et
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return utils.NewError(404, "Sipariş bulunamadı")
	}

	if role == "buyer" && order.UserID != userID {
		utils.Log.Warn("Yetkisiz iptal denemesi", zap.Uint("user_id", userID), zap.Uint("order_id", orderID))
		return utils.NewError(403, "Bu siparişi iptal etme yetkiniz yok")
	}

	if order.Status != domain.StatusPending {
		return utils.NewError(400, "Kargolanmış veya tamamlanmış sipariş iptal edilemez")
	}

	if err := s.orderRepo.CancelOrderWithStockRestore(ctx, orderID); err != nil {
		utils.Log.Error("İptal işlemi sırasında hata", zap.Error(err), zap.Uint("order_id", orderID))
		return err
	}

	utils.Log.Info("Sipariş İptal Edildi ve Stoklar İade Edildi",
		zap.Uint("order_id", orderID),
		zap.Uint("cancelled_by", userID),
	)

	return nil
}
