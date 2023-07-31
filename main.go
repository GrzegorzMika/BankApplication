package main

import (
	"context"
	"log"

	"BankApplication/internal/api"
	"BankApplication/internal/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource      = "postgres://postgres:password@localhost:54322/bank_tests?sslmode=disable"
	serverAddress = "localhost:8080"
)

func main() {
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatalln("Cannot connect to db: ", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatalln("Cannot start server: ", err)
	}
}
