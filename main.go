package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/guy-ronen/simplebank/api"
	db "github.com/guy-ronen/simplebank/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource =  "postgresql://guy.ronen:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress 		= "0.0.0.0:8080"
)

func main() {
	
	// connect to the database
	conn, err := sql.Open(dbDriver, dbSource)

	// check if there is an error
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	
	// close the connection when the function exits
	store := db.NewStore(conn)
	
	// create a new server
	server := api.NewServer(store)
	
	// start the server
	err = server.Start(serverAddress)

	// check if there is an error
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}