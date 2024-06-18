package main

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestGetCustomer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("^SELECT id, name, surname, age, phone_number, debit_card, credit_card, date_of_birth, date_of_issue, issuing_authority, has_foreign_country_tax_liability FROM customers\\.public\\.customers WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "surname", "age", "phone_number", "debit_card", "credit_card", "date_of_birth", "date_of_issue", "issuing_authority", "has_foreign_country_tax_liability"}).
			AddRow(1, "Vi", "N", 20, "1234567890", "1234-5678-9101-1121", "5432-1098-7654-3210", time.Now(), time.Now(), "Authority XYZ", false))

	req, err := http.NewRequest("GET", "/customer/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/customer/{id}", func(w http.ResponseWriter, r *http.Request) {
		customerID := 1

		var c Customer
		err := db.QueryRow("SELECT id, name, surname, age, phone_number, debit_card, credit_card, date_of_birth, date_of_issue, issuing_authority, has_foreign_country_tax_liability FROM customers.public.customers WHERE id = $1", customerID).
			Scan(&c.ID, &c.Name, &c.Surname, &c.Age, &c.PhoneNumber, &c.DebitCard, &c.CreditCard, &c.DateOfBirth, &c.DateOfIssue, &c.IssuingAuthority, &c.HasForeignCountryTaxLiability)
		if err != nil {
			t.Errorf("Error querying customer by ID: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(c)
	}).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", status)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error verifying mock database expectations: %v", err)
	}
}

func TestGetCustomers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "surname", "age", "phone_number", "debit_card", "credit_card", "date_of_birth", "date_of_issue", "issuing_authority", "has_foreign_country_tax_liability"}).
		AddRow(1, "fdsg", "gfd", 30, "1234567890", "1234-5678-9101-1121", "5432-1098-7654-3210", time.Now(), time.Now(), "Authority XYZ", false).
		AddRow(2, "fghf", "ghjgh", 28, "9876543210", "5678-9101-1121-3141", "8765-4321-0987-6543", time.Now(), time.Now(), "Authority ABC", true)

	mock.ExpectQuery("^SELECT id, name, surname, age, phone_number, debit_card, credit_card, date_of_birth, date_of_issue, issuing_authority, has_foreign_country_tax_liability FROM customers\\.public\\.customers").
		WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/customer", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/customer", GetCustomers).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", status)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error verifying mock database expectations: %v", err)
	}
}

func TestCreateCustomer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	newCustomer := Customer{
		Name:                          "Victoria",
		Surname:                       "N",
		Age:                           56,
		PhoneNumber:                   "1234567890",
		DebitCard:                     "1234-5678-9101-1121",
		CreditCard:                    "5432-1098-7654-3210",
		DateOfBirth:                   time.Now(),
		DateOfIssue:                   time.Now(),
		IssuingAuthority:              "Authority XYZ",
		HasForeignCountryTaxLiability: false,
	}

	mock.ExpectExec("INSERT INTO customers.public.customers").
		WithArgs(newCustomer.Name, newCustomer.Surname, newCustomer.Age, newCustomer.PhoneNumber, newCustomer.DebitCard, newCustomer.CreditCard, newCustomer.DateOfBirth, newCustomer.DateOfIssue, newCustomer.IssuingAuthority, newCustomer.HasForeignCountryTaxLiability).
		WillReturnResult(sqlmock.NewResult(1, 1))

	reqBody, err := json.Marshal(newCustomer)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/customer", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/customer", CreateCustomer).Methods("POST")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status 201 Created, got %d", status)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error verifying mock database expectations: %v", err)
	}
}

func TestUpdateCustomer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	updatedCustomer := Customer{
		Name:                          "test name",
		Surname:                       "test",
		Age:                           25,
		PhoneNumber:                   "9876543210",
		DebitCard:                     "5678-9101-1121-3141",
		CreditCard:                    "8765-4321-0987-6543",
		DateOfBirth:                   time.Now(),
		DateOfIssue:                   time.Now(),
		IssuingAuthority:              "Authority ABC",
		HasForeignCountryTaxLiability: true,
	}

	mock.ExpectExec("UPDATE customers.public.customers SET").
		WithArgs(updatedCustomer.Name, updatedCustomer.Surname, updatedCustomer.Age, updatedCustomer.PhoneNumber, updatedCustomer.DebitCard, updatedCustomer.CreditCard, updatedCustomer.DateOfBirth, updatedCustomer.DateOfIssue, updatedCustomer.IssuingAuthority, updatedCustomer.HasForeignCountryTaxLiability, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	reqBody, err := json.Marshal(updatedCustomer)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/customer/1", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/customer/{id}", UpdateCustomer).Methods("PUT")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", status)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error verifying mock database expectations: %v", err)
	}
}

func TestDeleteCustomer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM customers.public.customers WHERE id = $1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Return a successful delete operation

	req, err := http.NewRequest("DELETE", "/customer/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/customer/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		customerID := vars["id"]

		id, err := strconv.Atoi(customerID)
		if err != nil {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			return
		}

		result, err := db.Exec("DELETE FROM customers.public.customers WHERE id = $1", id)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Customer not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("DELETE")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", status)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error verifying mock database expectations: %v", err)
	}
}
