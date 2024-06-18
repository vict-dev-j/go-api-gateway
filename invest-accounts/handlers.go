package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func GetInvestAccounts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM invest_accounts.public.invest_accounts")
	if err != nil {
		log.Println("Error querying invest account:", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	defer rows.Close()

	var investAccounts []InvestAccount
	for rows.Next() {
		var c InvestAccount
		err := rows.Scan(&c.ID, &c.OwnerId, &c.ClientSurveyNumber, &c.Share, &c.InvestedAmountOfMoney, &c.FreeAmountOfMoney)
		if err != nil {
			log.Println("Error scanning invest account row:", err)
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		investAccounts = append(investAccounts, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(investAccounts)
}

func GetInvestAccount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var c InvestAccount
	err := db.QueryRow("SELECT id, owner_id, client_survey_number, share, invested_amount_of_money, free_amount_of_money FROM invest_accounts.public.invest_accounts WHERE id = $1", id).Scan(
		&c.ID, &c.OwnerId, &c.ClientSurveyNumber, &c.Share, &c.InvestedAmountOfMoney, &c.FreeAmountOfMoney,
	)
	if err != nil {
		log.Println("Error querying invest account by ID:", err)
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Invest account not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func CreateInvestAccount(w http.ResponseWriter, r *http.Request) {
	var newAccount InvestAccount
	err := json.NewDecoder(r.Body).Decode(&newAccount)
	if err != nil {
		log.Println("Error decoding request body:", err)
		respondWithError(w, http.StatusBadRequest, "Bad request")
		return
	}

	err = db.QueryRow("INSERT INTO invest_accounts.public.invest_accounts(owner_id, client_survey_number, share, invested_amount_of_money, free_amount_of_money) VALUES($1, $2, $3, $4, $5) RETURNING id", newAccount.OwnerId, newAccount.ClientSurveyNumber, newAccount.Share, newAccount.InvestedAmountOfMoney, newAccount.FreeAmountOfMoney).Scan(&newAccount.ID)
	if err != nil {
		log.Println("Error inserting new customer:", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAccount)
}

func UpdateInvestAccount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var updatedAccount InvestAccount
	err := json.NewDecoder(r.Body).Decode(&updatedAccount)
	if err != nil {
		log.Println("Error decoding request body:", err)
		respondWithError(w, http.StatusBadRequest, "Bad request")
		return
	}

	_, err = db.Exec("UPDATE invest_accounts.public.invest_accounts SET owner_id=$1, client_survey_number=$2, share=$3, invested_amount_of_money=$4, free_amount_of_money=$5 WHERE id=$6", updatedAccount.OwnerId, updatedAccount.ClientSurveyNumber, updatedAccount.Share, updatedAccount.InvestedAmountOfMoney, updatedAccount.FreeAmountOfMoney, id)
	if err != nil {
		log.Println("Error updating customer:", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteInvestAccount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM invest_accounts.public.invest_accounts WHERE id = $1", id)
	if err != nil {
		log.Println("Error deleting invest account:", err)
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
