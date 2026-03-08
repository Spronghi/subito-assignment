package repository

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/simonecolaci/subito-assignment/internal/entity"
)

type ProductRepository interface {
	Populate() error
	Create(p *entity.Product) (*entity.Product, error)
	GetByID(id int64) (*entity.Product, error)
	List() ([]*entity.Product, error)
}

type SQLiteProductRepository struct {
	db *sql.DB
}

func NewSQLiteProductRepository(db *sql.DB) (*SQLiteProductRepository, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS products (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		name        TEXT    NOT NULL,
		description TEXT    NOT NULL DEFAULT '',
		price       INTEGER NOT NULL,
		vat_rate    REAL    NOT NULL DEFAULT 0.22,
		created_at  DATETIME NOT NULL,
		updated_at  DATETIME NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	return &SQLiteProductRepository{db: db}, nil
}

func (r *SQLiteProductRepository) Populate() error {
	products := []entity.Product{
		{Name: "Laptop", Description: "A powerful laptop", Price: 80000, VATRate: 0.22},
		{Name: "Smartphone", Description: "A modern smartphone", Price: 50000, VATRate: 0.22},
		{Name: "Headphones", Description: "Noise-cancelling headphones", Price: 20000, VATRate: 0.22},
	}
	for i := range products {
		if _, err := r.Create(&products[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *SQLiteProductRepository) List() ([]*entity.Product, error) {
	rows, err := r.db.Query(
		`SELECT id, name, description, price, vat_rate, created_at, updated_at FROM products ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("Failed to close rows", "err", err)
		}
	}()

	var products []*entity.Product
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *SQLiteProductRepository) GetByID(id int64) (*entity.Product, error) {
	row := r.db.QueryRow(
		`SELECT id, name, description, price, vat_rate, created_at, updated_at FROM products WHERE id = ?`, id,
	)
	p, err := scanProduct(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, entity.ErrNotFound
	}
	return p, err
}

func (r *SQLiteProductRepository) Create(p *entity.Product) (*entity.Product, error) {
	now := time.Now().UTC()
	result, err := r.db.Exec(
		`INSERT INTO products (name, description, price, vat_rate, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		p.Name, p.Description, p.Price, p.VATRate, now, now,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

type productScanner interface {
	Scan(dest ...any) error
}

func scanProduct(s productScanner) (*entity.Product, error) {
	var p entity.Product
	err := s.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.VATRate, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
