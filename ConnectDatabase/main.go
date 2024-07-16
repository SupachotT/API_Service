package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Customer struct {
	First_name string
	Last_name  string
	Phone      string
	Email      string
}

func main() {
	connStr := "postgres://postgres:secret@localhost:5432/gopgtest?sslmode=disable"
	// UserName, Password, HostName, dbName

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	createLoanCustomerTable(db)

	customer := Customer{"Benjamin", "Franklin", "034-591201", "admin@loomsoom.go.th"}
	pk := InsertLoanCustomer(db, customer)

	fmt.Printf("ID = %d\n", pk)
}

func createLoanCustomerTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS customers (
		customer_id SERIAL PRIMARY KEY,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		phone VARCHAR(15) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertLoanCustomer(db *sql.DB, customer Customer) int {
	query := `INSERT INTO Customers (first_name, last_name, phone, email)
		VALUES ($1, $2, $3, $4) RETURNING customer_id`

	var pk int
	err := db.QueryRow(query, customer.First_name, customer.Last_name, customer.Phone, customer.Email).Scan(&pk)
	if err != nil {
		log.Fatal(err)
	}
	return pk

}
