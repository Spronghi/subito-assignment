package repository

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/simonecolaci/subito-assignment/internal/entity"
)

type OrderRepository interface {
	Populate() error
	Create(o *entity.Order) (*entity.Order, error)
	GetByID(id int64) (*entity.Order, error)
	List() ([]*entity.Order, error)
}

type SQLiteOrderRepository struct {
	db *sql.DB
}

func NewSQLiteOrderRepository(db *sql.DB) (*SQLiteOrderRepository, error) {
	// TODO: move this on a separate migration package and run it at application startup, this way we can also handle schema migrations in the future
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS orders (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		total_price INTEGER NOT NULL,
		total_vat   INTEGER NOT NULL,
		created_at  DATETIME NOT NULL
	)`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS order_items (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		order_id     INTEGER NOT NULL REFERENCES orders(id),
		product_id   INTEGER NOT NULL,
		product_name TEXT    NOT NULL,
		quantity     INTEGER NOT NULL,
		unit_price   INTEGER NOT NULL,
		vat_rate     REAL    NOT NULL,
		price        INTEGER NOT NULL,
		vat          INTEGER NOT NULL
	)`)
	if err != nil {
		return nil, err
	}
	return &SQLiteOrderRepository{db: db}, nil
}

func (r *SQLiteOrderRepository) Populate() error {
	orders := []entity.Order{
		{
			Items: []entity.OrderItem{
				{ProductID: 1, ProductName: "Laptop", Quantity: 1, UnitPrice: 80000, VATRate: 0.22, Price: 80000, VAT: 17600},
				{ProductID: 2, ProductName: "Smartphone", Quantity: 2, UnitPrice: 50000, VATRate: 0.22, Price: 100000, VAT: 22000},
			},
			TotalPrice: 180000,
			TotalVAT:   39600,
			CreatedAt:  time.Now(),
		},
		{
			Items: []entity.OrderItem{
				{ProductID: 3, ProductName: "Headphones", Quantity: 1, UnitPrice: 20000, VATRate: 0.22, Price: 20000, VAT: 4400},
			},
			TotalPrice: 20000,
			TotalVAT:   4400,
			CreatedAt:  time.Now(),
		},
	}

	for i := range orders {
		if _, err := r.Create(&orders[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *SQLiteOrderRepository) GetByID(id int64) (*entity.Order, error) {
	row := r.db.QueryRow(
		`SELECT id, total_price, total_vat, created_at FROM orders WHERE id = ?`,
		id,
	)

	o, err := scanOrder(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, entity.ErrNotFound
	}

	o.Items, err = r.getOrderItemsByOrderID(o.ID)
	if err != nil {
		return nil, err
	}

	return o, err
}

func (r *SQLiteOrderRepository) getOrderItemsByOrderID(id int64) ([]entity.OrderItem, error) {
	rows, err := r.db.Query(
		`SELECT id, product_id, product_name, quantity, unit_price, vat_rate, price, vat FROM order_items WHERE order_id = ?`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("Failed to close rows", "err", err)
		}
	}()

	var items []entity.OrderItem
	for rows.Next() {
		item, err := scanOrderItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *SQLiteOrderRepository) Create(o *entity.Order) (*entity.Order, error) {
	now := time.Now()

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			slog.Error("Failed to rollback transaction", "err", err)
		}
	}()

	result, err := tx.Exec(
		`INSERT INTO orders (total_price, total_vat, created_at) VALUES (?, ?, ?)`,
		o.TotalPrice, o.TotalVAT, now,
	)
	if err != nil {
		return nil, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	createdItems, err := r.createOrderItems(o.Items, orderID, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	o.ID = orderID
	o.CreatedAt = now
	o.Items = createdItems

	return o, nil
}

func (r *SQLiteOrderRepository) createOrderItems(items []entity.OrderItem, orderID int64, tx *sql.Tx) ([]entity.OrderItem, error) {
	var createdItems []entity.OrderItem
	for _, item := range items {
		result, err := tx.Exec(
			`INSERT INTO order_items (order_id, product_id, product_name, quantity, unit_price, vat_rate, price, vat) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			orderID, item.ProductID, item.ProductName, item.Quantity, item.UnitPrice, item.VATRate, item.Price, item.VAT,
		)
		if err != nil {
			return nil, err
		}

		itemID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		item.ID = itemID
		createdItems = append(createdItems, item)
	}

	return createdItems, nil
}

func (r *SQLiteOrderRepository) List() ([]*entity.Order, error) {
	rows, err := r.db.Query(
		`SELECT id, total_price, total_vat, created_at FROM orders ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("Failed to close rows", "err", err)
		}
	}()

	var orders []*entity.Order
	for rows.Next() {
		o, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, o := range orders {
		o.Items, err = r.getOrderItemsByOrderID(o.ID)
		if err != nil {
			return nil, err
		}
	}

	return orders, nil
}

func scanOrder(s scanner) (*entity.Order, error) {
	var o entity.Order
	err := s.Scan(&o.ID, &o.TotalPrice, &o.TotalVAT, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func scanOrderItem(s scanner) (*entity.OrderItem, error) {
	var item entity.OrderItem
	err := s.Scan(&item.ID, &item.ProductID, &item.ProductName, &item.Quantity, &item.UnitPrice, &item.VATRate, &item.Price, &item.VAT)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
