package service

import (
	"github.com/simonecolaci/subito-assignment/internal/entity"
	"github.com/simonecolaci/subito-assignment/internal/repository"
)

type ProductService interface {
	Create(p *entity.Product) (*entity.Product, error)
	GetByID(id int64) (*entity.Product, error)
	List() ([]*entity.Product, error)
	Update(id int64, p *entity.Product) (*entity.Product, error)
	Delete(id int64) error
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{productRepository: productRepository}
}

func (s *productService) Create(p *entity.Product) (*entity.Product, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}

	return s.productRepository.Create(p)
}

func (s *productService) GetByID(id int64) (*entity.Product, error) {
	return s.productRepository.GetByID(id)
}

func (s *productService) List() ([]*entity.Product, error) {
	return s.productRepository.List()
}

func (s *productService) Update(id int64, p *entity.Product) (*entity.Product, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	p.ID = id
	return s.productRepository.Update(p)
}

func (s *productService) Delete(id int64) error {
	return s.productRepository.Delete(id)
}
