package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

var testQueries *Queries

const (
	dbSource = "postgres://postgres:password@localhost:54322/bank_tests?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := pgx.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatalln("Cannot connect to db: ", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
