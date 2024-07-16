package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Customer struct {
	Customer_id int
	First_name  string
	Last_name   string
	Phone       string
	Email       string
	Created_at  string
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

	// Read data from JSON file
	customers, err := readCustomersFromFile("Json/customers.json")
	if err != nil {
		log.Fatal(err)
	}

	// Insert each customer into the database
	for _, customer := range customers {
		pk := InsertLoanCustomer(db, customer)
		fmt.Printf("Inserted customer with ID = %d\n", pk)
	}

	router := mux.NewRouter()

	// Define API endpoints
	router.HandleFunc("/customers", getCustomersHandler).Methods("GET")
	router.HandleFunc("/customers/{id}", getCustomerByIDHandler).Methods("GET")
	router.HandleFunc("/customers/{id}", updateCustomerHandler).Methods("PUT")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", router))
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

func readCustomersFromFile(filename string) ([]Customer, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var customers []Customer
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&customers); err != nil {
		return nil, err
	}

	return customers, nil
}

func getCustomersHandler(w http.ResponseWriter, r *http.Request) {
	// connect databases
	connStr := "postgres://postgres:secret@localhost:5432/gopgtest?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// query from table customer
	rows, err := db.Query("SELECT customer_id, first_name, last_name, phone, email, created_at FROM customers")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.Customer_id, &customer.First_name, &customer.Last_name, &customer.Phone, &customer.Email, &customer.Created_at); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		customers = append(customers, customer)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func getCustomerByIDHandler(w http.ResponseWriter, r *http.Request) {
	connStr := "postgres://postgres:secret@localhost:5432/gopgtest?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get customer_id from URL parameters
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	// Query database for customer with given customer_id
	var customer Customer
	query := `SELECT customer_id, first_name, last_name, phone, email, created_at FROM customers WHERE customer_id = $1`
	err = db.QueryRow(query, id).Scan(&customer.Customer_id, &customer.First_name, &customer.Last_name, &customer.Phone, &customer.Email, &customer.Created_at)
	if err == sql.ErrNoRows {
		// Return JSON error response if no customer with the given ID exists
		errorResponse := map[string]string{"error": "Customer not found"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func updateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	connStr := "postgres://postgres:secret@localhost:5432/gopgtest?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get customer_id from URL parameters
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	// Decode JSON request body into a Customer struct
	var customer Customer
	err = json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update query
	query := `UPDATE customers SET first_name = $2, last_name = $3, phone = $4, email = $5 WHERE customer_id = $1`
	result, err := db.Exec(query, id, customer.First_name, customer.Last_name, customer.Phone, customer.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		// Return JSON error response if no customer with the given ID was found to update
		errorResponse := map[string]string{"error": "Customer ID not found or no update performed"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Return success message
	w.WriteHeader(http.StatusOK)
	successMessage := map[string]string{"message": fmt.Sprintf("Customer with ID %d updated successfully", id)}
	json.NewEncoder(w).Encode(successMessage)
}
