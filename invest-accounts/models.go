package main

type InvestAccount struct {
	ID                    int     `json:"id"`
	OwnerId               int     `json:"owner_id"`
	ClientSurveyNumber    int     `json:"client_survey_number"`
	Share                 string  `json:"share"`
	InvestedAmountOfMoney float64 `json:"invested_amount_of_money"`
	FreeAmountOfMoney     float64 `json:"free_amount_of_money"`
}
