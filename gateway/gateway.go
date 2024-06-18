package main

import (
	"fmt"
	"io"
	"net/http"
	_ "strings"

	"github.com/gorilla/mux"
)

var (
	CustomersURL      = "http://localhost:8080"
	InvestAccountsURL = "http://localhost:8082"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{service:customer|invest-account}{rest:.*}", Handler)
	fmt.Println("Gateway listening on :8081")
	err := http.ListenAndServe(":8081", router)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	rest := vars["rest"]

	var targetURL string
	switch service {
	case "customer":
		targetURL = CustomersURL + "/customer" + rest
	case "invest-account":
		targetURL = InvestAccountsURL + "/invest-account" + rest
	default:
		http.Error(w, "Path not supported", http.StatusNotFound)
		return
	}

	proxyRequest(w, r, targetURL)
}

func proxyRequest(w http.ResponseWriter, r *http.Request, targetURL string) {
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	copyHeaders(req.Header, r.Header)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error proxying request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Error copying response", http.StatusInternalServerError)
		return
	}
}

func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}
