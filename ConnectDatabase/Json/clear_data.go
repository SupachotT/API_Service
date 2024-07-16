package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://postgres:secret@localhost:5432/gopgtest?sslmode=disable"

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Clear data from the customers table
	_, err = db.Exec("DELETE FROM customers")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data cleared from the customers table.")
}
