package service

import (
	"fmt"

	"github.com/simonecolaci/subito-assignment/internal/entity"
	"github.com/simonecolaci/subito-assignment/internal/repository"
)

type OrderService interface {
	Create(o *entity.NewOrder) (*entity.Order, error)
	GetByID(id int64) (*entity.Order, error)
	List() ([]*entity.Order, error)
}

type orderService struct {
	orderRepository   repository.OrderRepository
	productRepository repository.ProductRepository
}

func NewOrderService(orderRepository repository.OrderRepository, productRepository repository.ProductRepository) OrderService {
	return &orderService{
		orderRepository:   orderRepository,
		productRepository: productRepository,
	}
}

func (s *orderService) Create(o *entity.NewOrder) (*entity.Order, error) {
	if err := o.Validate(); err != nil {
		return nil, err
	}

	order := &entity.Order{
		Items: make([]entity.OrderItem, 0, len(o.Items)),
	}

	for _, item := range o.Items {
		product, err := s.productRepository.GetByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product with id %d: %w", item.ProductID, entity.ErrNotFound)
		}

		orderItem := entity.OrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   product.Price,
			VATRate:     product.VATRate,
			Price:       product.Price * int64(item.Quantity),
			VAT:         product.VATAmount() * int64(item.Quantity),
		}

		order.Items = append(order.Items, orderItem)
	}

	order.TotalPrice += order.CalculateTotalPrice()
	order.TotalVAT += order.CalculateTotalVAT()

	return s.orderRepository.Create(order)
}

func (s *orderService) GetByID(id int64) (*entity.Order, error) {
	return s.orderRepository.GetByID(id)
}

func (s *orderService) List() ([]*entity.Order, error) {
	return s.orderRepository.List()
}
