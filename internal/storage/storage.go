package storage

import (
	"context"
	"fmt"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type PostgresStorage struct {
	conn *pgx.Conn
}

func NewPostgresStorage(ctx context.Context, connStr string) (*PostgresStorage, error) {
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	// defer conn.Close(context.Background())

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	if err := CreatePostgresDB(ctx, conn); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		conn: conn,
	}, nil
}

func CreatePostgresDB(ctx context.Context, conn *pgx.Conn) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT,
		password TEXT,
		email TEXT
	);
	
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT,
		description TEXT,
		price INTEGER,
		quantity INTEGER
	);
	
	CREATE TABLE IF NOT EXISTS purchases (
		id SERIAL PRIMARY KEY,
		user_id INTEGER,
		product_id INTEGER,
		quantity INTEGER,
		timestamp TEXT
	);
	`

	_, err := conn.Exec(ctx, query)
	return err
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}

	return string(hashedPassword), nil
}

func VerifyPassword(savedHash, inputPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(savedHash), []byte(inputPassword)); err != nil {
		return err
	}

	return nil
}

func IsDataUnique(conn *pgx.Conn, login string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = $1`
	err := conn.QueryRow(ctx, query, login).Scan(&count)
	if err != nil {
		return err
	}

	if count != 0 {
		return fmt.Errorf("non unique data")
	}

	return nil
}
