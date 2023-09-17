package main

import (
	"log"
	"fmt"
	"context"
	"net/http"
	"github.com/NeelkanthSingh/Klair/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "You've requested the server")
    })

	//Create database connection
	connPool,err := pgxpool.NewWithConfig(context.Background(), db.Config())
	if err!=nil {
		log.Fatal("Error creating connection to the database!!")
	} 

	connection, err := connPool.Acquire(context.Background())
	if err!=nil {
		log.Fatal("Error acquiring connection from the database pools!!")
	} 
	defer connection.Release()

	err = connection.Ping(context.Background())
	if err!=nil{
		log.Fatal("Could not ping database")
	}
	fmt.Println("Connected to the database!!")
	
	db.CreateTableQuery(connPool)
	
	db.InsertQuery(connPool)

	defer connPool.Close()

	http.ListenAndServe(":8080", router)
}
