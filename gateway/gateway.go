package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var targetURL = "http://localhost:8080"

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Gateway listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}

	req, err := http.NewRequest(r.Method, targetURL+r.URL.Path, r.Body)
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}
