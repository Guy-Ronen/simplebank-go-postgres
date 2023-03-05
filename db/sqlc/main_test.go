package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	util "github.com/guy-ronen/simplebank/util"
	_ "github.com/lib/pq"
)

// dbDriver is the database driver

// testQueries is the global queries object for all the tests
var testQueries *Queries


// testDB is the global db object for all the tests
var testDB *sql.DB

func TestMain(m *testing.M) {
	// read the config file
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// open a connection to the test db
	testDB, err = sql.Open(config.DBDRIVER, config.DBSOURCE)

	// if there is an error, log it and exit
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// create a new instance of the queries object
	testQueries = New(testDB)

	// run the tests
	os.Exit(m.Run())
}