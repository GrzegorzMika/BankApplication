package db

import (
	"context"
	"log"
	"os"
	"testing"

	"BankApplication/internal/util"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatalln("Cannot load config file", err)
	}

	testDB, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalln("Cannot connect to db: ", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
