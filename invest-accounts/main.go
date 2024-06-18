package main

import (
	"github.com/gorilla/mux"

	"log"
	"net/http"
)

func main() {
	initDB()

	router := mux.NewRouter()
	router.HandleFunc("/customer", GetInvestAccounts).Methods("GET")
	router.HandleFunc("/customer/{id}", GetInvestAccount).Methods("GET")
	router.HandleFunc("/customer", CreateInvestAccount).Methods("POST")
	router.HandleFunc("/customer/{id}", UpdateInvestAccount).Methods("PUT")
	router.HandleFunc("/customer/{id}", DeleteInvestAccount).Methods("DELETE")

	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":8082", router))
}
