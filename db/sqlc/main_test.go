package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// dbDriver is the database driver
const (
	dbDriver = "postgres"
	dbSource =  "postgresql://guy.ronen:secret@localhost:5432/simple_bank?sslmode=disable"
)

// testQueries is the global queries object for all the tests
var testQueries *Queries


// testDB is the global db object for all the tests
var testDB *sql.DB

func TestMain(m *testing.M) {
	// connect to test db
	var err error

	// open a connection to the test db
	testDB, err = sql.Open(dbDriver, dbSource)

	// if there is an error, log it and exit
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// create a new instance of the queries object
	testQueries = New(testDB)

	// run the tests
	os.Exit(m.Run())
}