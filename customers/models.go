package main

import "time"

type Customer struct {
	ID                            int       `json:"id"`
	Name                          string    `json:"name"`
	Surname                       string    `json:"surname"`
	Age                           int       `json:"age"`
	PhoneNumber                   string    `json:"phone_number"`
	DebitCard                     string    `json:"debit_card"`
	CreditCard                    string    `json:"credit_card"`
	DateOfBirth                   time.Time `json:"date_of_birth"`
	DateOfIssue                   time.Time `json:"date_of_issue"`
	IssuingAuthority              string    `json:"issuing_authority"`
	HasForeignCountryTaxLiability bool      `json:"has_foreign_country_tax_liability"`
}
