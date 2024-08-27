package storage

import (
	"context"
	"sort"
	"time"

	"github.com/ursuldaniel/go-market/internal/domain/models"
)

func (s *PostgresStorage) AddProduct(name, description string, price, quantity int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Prepare(ctx, "insert product", "INSERT INTO products (name, description, price, quantity) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "insert product", name, description, price, quantity)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *PostgresStorage) GetAllProducts() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `SELECT * FROM products`
	rows, err := s.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	products := []models.Product{}
	for rows.Next() {
		product := models.Product{}
		if err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Quantity); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	sortById := func(i, j int) bool {
		return products[i].Id < products[j].Id
	}
	sort.Slice(products, sortById)

	return products, nil
}

func (s *PostgresStorage) GetProductById(productId int) (models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `SELECT * FROM products WHERE id = $1`
	rows, err := s.conn.Query(ctx, query, productId)
	if err != nil {
		return models.Product{}, err
	}
	defer rows.Close()

	product := models.Product{}
	for rows.Next() {
		if err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Quantity); err != nil {
			return models.Product{}, err
		}

		if err != nil {
			return models.Product{}, err
		}
	}

	return product, nil
}

func (s *PostgresStorage) UpdateProduct(productId int, name, description string, price, quantity int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Prepare(ctx, "update", "UPDATE products SET name = $1, description = $2, price = $3, quantity = $4 WHERE id = $5")
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "update", name, description, price, quantity, productId)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *PostgresStorage) DeleteProduct(productId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `DELETE FROM products WHERE id = $1`
	_, err := s.conn.Exec(ctx, query, productId)
	if err != nil {
		return err
	}

	return nil
}
