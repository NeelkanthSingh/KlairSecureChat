package main

import (
	"awesomeProject/handlers"
	"awesomeProject/middleware"
	"context"
	"fmt"
	"github.com/NeelkanthSingh/Klair/db"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func main() {
	router := mux.NewRouter()

	router.Use(middleware.Authorize)
	router.HandleFunc("/", handlers.HomeHandler)

	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.IsAuthenticated)
	adminRouter.HandleFunc("/andy", handlers.AdminHandler)
	adminRouter.HandleFunc("/logout", handlers.LogoutHandler)

	router.HandleFunc("/token", handlers.SignUpHandler).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "You've requested the server")
	})

	//Create database connection
	connPool, err := pgxpool.NewWithConfig(context.Background(), db.Config())
	if err != nil {
		log.Fatal("Error creating connection to the database!!")
	}

	connection, err := connPool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Error acquiring connection from the database pools!!")
	}
	defer connection.Release()

	err = connection.Ping(context.Background())
	if err != nil {
		log.Fatal("Could not ping database")
	}
	fmt.Println("Connected to the database!!")

	db.CreateTableQuery(connPool)

	db.InsertQuery(connPool)

	defer connPool.Close()

	http.ListenAndServe(":8080", router)
}
