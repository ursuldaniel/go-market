package storage

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/ursuldaniel/go-market/internal/domain/models"
)

func (s *PostgresStorage) MakePurchase(userID, productID, quantity int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var productQuantity int
	query := `SELECT quantity FROM products WHERE id = $1`
	err := s.conn.QueryRow(ctx, query, productID).Scan(&productQuantity)
	if err != nil {
		return err
	}

	if productQuantity-quantity < 0 {
		return fmt.Errorf("not enough products")
	}

	query = `UPDATE products SET quantity = $1 WHERE id = $2`
	_, err = s.conn.Exec(ctx, query, productQuantity-quantity, productID)
	if err != nil {
		return err
	}

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Prepare(ctx, "insert purchase", "INSERT INTO purchases (user_id, product_id, quantity, timestamp) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "insert purchase", userID, productID, quantity, (time.Now().String())[:19])
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *PostgresStorage) GetUserPurchases(userID int) ([]models.Purchase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `SELECT * FROM purchases WHERE user_id = $1`
	rows, err := s.conn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	purchases := []models.Purchase{}
	for rows.Next() {
		purchase := models.Purchase{}
		if err := rows.Scan(&purchase.Id, &purchase.UserId, &purchase.ProductId, &purchase.Quantity, &purchase.Timestamp); err != nil {
			return nil, err
		}

		purchases = append(purchases, purchase)
	}

	sortById := func(i, j int) bool {
		return purchases[i].Id < purchases[j].Id
	}
	sort.Slice(purchases, sortById)

	return purchases, nil
}

func (s *PostgresStorage) GetProductPurchases(productID int) ([]models.Purchase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `SELECT * FROM purchases WHERE product_id = $1`
	rows, err := s.conn.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}

	purchases := []models.Purchase{}
	for rows.Next() {
		purchase := models.Purchase{}
		if err := rows.Scan(&purchase.Id, &purchase.UserId, &purchase.ProductId, &purchase.Quantity, &purchase.Timestamp); err != nil {
			return nil, err
		}

		purchases = append(purchases, purchase)
	}

	sortById := func(i, j int) bool {
		return purchases[i].Id < purchases[j].Id
	}
	sort.Slice(purchases, sortById)

	return purchases, nil
}
