package main

import (
	"github.com/gorilla/mux"

	"log"
	"net/http"
)

func main() {
	initDB()

	router := mux.NewRouter()
	router.HandleFunc("/invest-account", GetInvestAccounts).Methods("GET")
	router.HandleFunc("/invest-account/{id}", GetInvestAccount).Methods("GET")
	router.HandleFunc("/invest-account", CreateInvestAccount).Methods("POST")
	router.HandleFunc("/invest-account/{id}", UpdateInvestAccount).Methods("PUT")
	router.HandleFunc("/invest-account/{id}", DeleteInvestAccount).Methods("DELETE")

	log.Println("Server started on port 8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}
