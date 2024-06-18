package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	postgres *sql.DB
	r        *mux.Router
	err      error
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	initDB()
}

func teardown() {
	clearTestData()
	db.Close()
}

func TestGetInvestAccounts(t *testing.T) {
	clearTestData()
	mockInvestAccounts := []InvestAccount{
		{OwnerId: 1, ClientSurveyNumber: 123, Share: "ABC", InvestedAmountOfMoney: 1000.0, FreeAmountOfMoney: 500.0},
		{OwnerId: 2, ClientSurveyNumber: 456, Share: "DEF", InvestedAmountOfMoney: 2000.0, FreeAmountOfMoney: 1000.0},
	}
	insertMockInvestAccounts(mockInvestAccounts)

	req, err := http.NewRequest("GET", "/invest_accounts", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response []InvestAccount
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error decoding JSON response: %v", err)
	}

	if len(response) != len(mockInvestAccounts) {
		t.Errorf("Expected %d invest accounts, got %d", len(mockInvestAccounts), len(response))
	}
}

func TestGetInvestAccount(t *testing.T) {
	clearTestData()
	mockAccount := InvestAccount{
		OwnerId:               1,
		ClientSurveyNumber:    123,
		Share:                 "ABC",
		InvestedAmountOfMoney: 1000.0,
		FreeAmountOfMoney:     500.0,
	}
	insertMockInvestAccount(mockAccount)

	req, err := http.NewRequest("GET", fmt.Sprintf("/invest_accounts/%d", mockAccount.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response InvestAccount
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error decoding JSON response: %v", err)
	}

	if response.ID != mockAccount.ID {
		t.Errorf("Expected ID %d, got %d", mockAccount.ID, response.ID)
	}
}

func TestCreateInvestAccount(t *testing.T) {
	clearTestData()
	mockAccount := InvestAccount{
		OwnerId:               1,
		ClientSurveyNumber:    123,
		Share:                 "ABC",
		InvestedAmountOfMoney: 1000.0,
		FreeAmountOfMoney:     500.0,
	}
	body, err := json.Marshal(mockAccount)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/invest_accounts", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response InvestAccount
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error decoding JSON response: %v", err)
	}

	if response.OwnerId != mockAccount.OwnerId {
		t.Errorf("Expected owner ID %d, got %d", mockAccount.OwnerId, response.OwnerId)
	}
}

func TestUpdateInvestAccount(t *testing.T) {
	clearTestData()
	mockAccount := InvestAccount{
		OwnerId:               1,
		ClientSurveyNumber:    123,
		Share:                 "ABC",
		InvestedAmountOfMoney: 1000.0,
		FreeAmountOfMoney:     500.0,
	}
	insertMockInvestAccount(mockAccount)

	mockAccountToUpdate := InvestAccount{
		OwnerId:               2,
		ClientSurveyNumber:    456,
		Share:                 "DEF",
		InvestedAmountOfMoney: 2000.0,
		FreeAmountOfMoney:     1000.0,
	}
	body, err := json.Marshal(mockAccountToUpdate)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("/invest_accounts/%d", mockAccount.ID), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var updatedAccount InvestAccount
	err = db.QueryRow("SELECT owner_id, client_survey_number, share, invested_amount_of_money, free_amount_of_money FROM invest_accounts.public.invest_accounts WHERE id = $1", mockAccount.ID).Scan(
		&updatedAccount.OwnerId, &updatedAccount.ClientSurveyNumber, &updatedAccount.Share, &updatedAccount.InvestedAmountOfMoney, &updatedAccount.FreeAmountOfMoney,
	)
	if err != nil {
		t.Fatalf("Failed to retrieve updated account: %v", err)
	}

	if updatedAccount.OwnerId != mockAccountToUpdate.OwnerId {
		t.Errorf("Expected owner ID %d, got %d", mockAccountToUpdate.OwnerId, updatedAccount.OwnerId)
	}
}

func TestDeleteInvestAccount(t *testing.T) {
	clearTestData()
	mockAccount := InvestAccount{
		OwnerId:               1,
		ClientSurveyNumber:    123,
		Share:                 "ABC",
		InvestedAmountOfMoney: 1000.0,
		FreeAmountOfMoney:     500.0,
	}
	insertMockInvestAccount(mockAccount)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/invest_accounts/%d", mockAccount.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM invest_accounts.public.invest_accounts WHERE id = $1", mockAccount.ID).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to check if account was deleted: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected account with ID %d to be deleted, but it still exists", mockAccount.ID)
	}
}

func insertMockInvestAccounts(accounts []InvestAccount) {
	for _, account := range accounts {
		_, err := db.Exec("INSERT INTO invest_accounts.public.invest_accounts (owner_id, client_survey_number, share, invested_amount_of_money, free_amount_of_money) VALUES ($1, $2, $3, $4, $5)",
			account.OwnerId, account.ClientSurveyNumber, account.Share, account.InvestedAmountOfMoney, account.FreeAmountOfMoney)
		if err != nil {
			log.Fatalf("Failed to insert mock account: %v", err)
		}
	}
}

func insertMockInvestAccount(account InvestAccount) {
	_, err := db.Exec("INSERT INTO invest_accounts.public.invest_accounts (owner_id, client_survey_number, share, invested_amount_of_money, free_amount_of_money) VALUES ($1, $2, $3, $4, $5)",
		account.OwnerId, account.ClientSurveyNumber, account.Share, account.InvestedAmountOfMoney, account.FreeAmountOfMoney)
	if err != nil {
		log.Fatalf("Failed to insert mock account: %v", err)
	}
}

func clearTestData() {
	_, err := db.Exec("DELETE FROM invest_accounts.public.invest_accounts")
	if err != nil {
		log.Printf("Failed to clear test data: %v", err)
	}
}
