package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM customers.public.customers")
	if err != nil {
		log.Println("Error querying customers:", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.ID, &c.Name, &c.Age, &c.Tel, &c.DebitCard, &c.CreditCard)
		if err != nil {
			log.Println("Error scanning customer row:", err)
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		customers = append(customers, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var c Customer
	err := db.QueryRow("SELECT * FROM customers.public.customers WHERE id = $1", id).Scan(&c.ID, &c.Name, &c.Age, &c.Tel, &c.DebitCard, &c.CreditCard)
	if err != nil {
		log.Println("Error querying customer by ID:", err)
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Customer not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var newCustomer Customer
	err := json.NewDecoder(r.Body).Decode(&newCustomer)
	if err != nil {
		log.Println("Error decoding request body:", err)
		respondWithError(w, http.StatusBadRequest, "Bad request")
		return
	}

	err = db.QueryRow("INSERT INTO customers.public.customers(name, age, tel, debit_card, credit_card) VALUES($1, $2, $3, $4, $5) RETURNING id", newCustomer.Name, newCustomer.Age, newCustomer.Tel, newCustomer.DebitCard, newCustomer.CreditCard).Scan(&newCustomer.ID)
	if err != nil {
		log.Println("Error inserting new customer:", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCustomer)
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var updatedCustomer Customer
	err := json.NewDecoder(r.Body).Decode(&updatedCustomer)
	if err != nil {
		log.Println("Error decoding request body:", err)
		respondWithError(w, http.StatusBadRequest, "Bad request")
		return
	}

	_, err = db.Exec("UPDATE customers.public.customers SET name=$1, age=$2, tel=$3, debit_card=$4, credit_card=$5 WHERE id=$6", updatedCustomer.Name, updatedCustomer.Age, updatedCustomer.Tel, updatedCustomer.DebitCard, updatedCustomer.CreditCard, id)
	if err != nil {
		log.Println("Error updating customer:", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM customers.public.customers WHERE id = $1", id)
	if err != nil {
		log.Println("Error deleting customer:", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	errMsg := map[string]string{"error": message}
	json.NewEncoder(w).Encode(errMsg)
}
