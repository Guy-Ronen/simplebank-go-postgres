package main

import (
	"database/sql"
	"log"

	"github.com/guy-ronen/simplebank/api"
	db "github.com/guy-ronen/simplebank/db/sqlc"
	util "github.com/guy-ronen/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	// read the config file
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	// connect to the database
	conn, err := sql.Open(config.DBDRIVER, config.DBSOURCE)

	// check if there is an error
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	// close the connection when the function exits
	store := db.NewStore(conn)

	// create a new server
	server := api.NewServer(store)

	// start the server
	err = server.Start(config.ServerAddress)

	// check if there is an error
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
