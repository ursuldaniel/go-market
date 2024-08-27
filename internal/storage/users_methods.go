package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/ursuldaniel/go-market/internal/domain/models"
)

func (s *PostgresStorage) RegisterUser(username, password, email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := IsDataUnique(s.conn, username)
	if err != nil {
		return fmt.Errorf("non unique data")
	}

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Prepare(ctx, "insert user", "INSERT INTO users (username, password, email) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "insert user", username, hashedPassword, email)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *PostgresStorage) LoginUser(username, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `SELECT id, password FROM users WHERE username = $1`
	rows, err := s.conn.Query(ctx, query, username)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	var id int
	var savedHash string
	for rows.Next() {
		err := rows.Scan(
			&id,
			&savedHash,
		)

		if err != nil {
			return -1, err
		}
	}

	err = VerifyPassword(savedHash, password)
	if err != nil {
		return -1, fmt.Errorf("invalid data")
	}

	return id, nil
}

func (s *PostgresStorage) GetUserProfile(userId int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `SELECT * FROM users WHERE id = $1`
	rows, err := s.conn.Query(ctx, query, userId)
	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()

	user := models.User{}
	for rows.Next() {
		err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Password,
			&user.Email,
		)

		if err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}
