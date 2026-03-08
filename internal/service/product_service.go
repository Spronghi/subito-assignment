package service

import (
	"fmt"

	"github.com/simonecolaci/subito-assignment/internal/entity"
)

type ProductService interface {
	Create(p *entity.Product) (*entity.Product, error)
	GetByID(id int64) (*entity.Product, error)
	List() ([]*entity.Product, error)
}

type productService struct {
}

func NewProductService() ProductService {
	return &productService{}
}

func (s *productService) Create(p *entity.Product) (*entity.Product, error) {
	if err := validateProduct(p); err != nil {
		return nil, err
	}

	// TODO: call the repository
	p.ID = 1

	return p, nil
}

func (s *productService) GetByID(id int64) (*entity.Product, error) {
	if id == 0 {
		return nil, entity.ErrNotFound
	}

	// TODO: call the repository
	return &entity.Product{
		ID:          id,
		Name:        fmt.Sprintf("Product %d", id),
		Description: fmt.Sprintf("Description for product %d", id),
		Price:       10000 + id*1000,
		VATRate:     0.22,
	}, nil
}

func (s *productService) List() ([]*entity.Product, error) {
	// TODO: call the repository
	return []*entity.Product{
		{ID: 1, Name: "Product 1", Description: "Description for product 1", Price: 10000, VATRate: 0.22},
		{ID: 2, Name: "Product 2", Description: "Description for product 2", Price: 11000, VATRate: 0.22},
		{ID: 3, Name: "Product 3", Description: "Description for product 3", Price: 12000, VATRate: 0.22},
	}, nil
}

func validateProduct(p *entity.Product) error {
	if p.Name == "" {
		return fmt.Errorf("%w: name is required", entity.ErrInvalidInput)
	}
	if p.Price < 0 {
		return fmt.Errorf("%w: price must be >= 0", entity.ErrInvalidInput)
	}
	if p.VATRate < 0 {
		return fmt.Errorf("%w: vat_rate must be >= 0", entity.ErrInvalidInput)
	}
	return nil
}
