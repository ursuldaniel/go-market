package main

// 2006-01-02 15:04:05 -0700 //

import (
	"context"
	"log"
	"os"

	"github.com/ursuldaniel/go-market/internal/server"
	"github.com/ursuldaniel/go-market/internal/storage"
)

func main() {
	store, err := storage.NewPostgresStorage(context.TODO(), os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewServer(os.Getenv("LISTEN_ADDR"), store)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
