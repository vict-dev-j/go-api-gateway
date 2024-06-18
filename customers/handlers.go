package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, surname, age, phone_number, debit_card, credit_card, date_of_birth, date_of_issue, issuing_authority, has_foreign_country_tax_liability FROM customers.public.customers")
	if err != nil {
		log.Println("Error querying customers:", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.ID, &c.Name, &c.Surname, &c.Age, &c.PhoneNumber, &c.DebitCard, &c.CreditCard, &c.DateOfBirth, &c.DateOfIssue, &c.IssuingAuthority, &c.HasForeignCountryTaxLiability)
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
	idStr := params["id"]

	idStr = strings.TrimSpace(idStr)

	if idStr == "" {
		log.Println("Empty customer ID")
		respondWithError(w, http.StatusBadRequest, "Customer ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid customer ID:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	var c Customer
	err = db.QueryRow("SELECT id, name, surname, age, phone_number, debit_card, credit_card, date_of_birth, date_of_issue, issuing_authority, has_foreign_country_tax_liability FROM customers.public.customers WHERE id = $1", id).Scan(&c.ID, &c.Name, &c.Surname, &c.Age, &c.PhoneNumber, &c.DebitCard, &c.CreditCard, &c.DateOfBirth, &c.DateOfIssue, &c.IssuingAuthority, &c.HasForeignCountryTaxLiability)
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
	if db == nil {
		log.Println("Database connection is not initialized")
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	var newCustomer Customer
	err := json.NewDecoder(r.Body).Decode(&newCustomer)
	if err != nil {
		log.Println("Error decoding request body:", err)
		respondWithError(w, http.StatusBadRequest, "Bad request")
		return
	}

	err = db.QueryRow("INSERT INTO customers.public.customers(name, surname, age, phone_number, debit_card, credit_card, date_of_birth, date_of_issue, issuing_authority, has_foreign_country_tax_liability) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		newCustomer.Name, newCustomer.Surname, newCustomer.Age, newCustomer.PhoneNumber, newCustomer.DebitCard, newCustomer.CreditCard, newCustomer.DateOfBirth, newCustomer.DateOfIssue, newCustomer.IssuingAuthority, newCustomer.HasForeignCountryTaxLiability).Scan(&newCustomer.ID)
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

	_, err = db.Exec("UPDATE customers.public.customers SET name=$1, surname=$2, age=$3, phone_number=$4, debit_card=$5, credit_card=$6, date_of_birth=$7, date_of_issue=$8, issuing_authority=$9, has_foreign_country_tax_liability=$10 WHERE id=$11",
		updatedCustomer.Name, updatedCustomer.Surname, updatedCustomer.Age, updatedCustomer.PhoneNumber, updatedCustomer.DebitCard, updatedCustomer.CreditCard, updatedCustomer.DateOfBirth, updatedCustomer.DateOfIssue, updatedCustomer.IssuingAuthority, updatedCustomer.HasForeignCountryTaxLiability, id)
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
