package main

type Customer struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Tel        string `json:"tel"`
	DebitCard  string `json:"debit_card"`
	CreditCard string `json:"credit_card"`
}
