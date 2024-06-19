package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/customer/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/customer/{id}", Handler).Methods("GET")

	mockServer := httptest.NewServer(router)
	defer mockServer.Close()

	CustomersURL = mockServer.URL
	req.Header.Set("Authorization", "Bearer mockToken")

	Handler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id": 1, "name": "V N"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
