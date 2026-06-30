package main

type Party struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	Company string `json:"company"`
	Address string `json:"address"`
}

type FormData struct {
	Party1                   Party  `json:"party1"`
	Party2                   Party  `json:"party2"`
	EffectiveDate            string `json:"effectiveDate"`
	MndaTermType             string `json:"mndaTermType"`
	MndaTermYears            int    `json:"mndaTermYears"`
	ConfidentialityTermType  string `json:"confidentialityTermType"`
	ConfidentialityTermYears int    `json:"confidentialityTermYears"`
	Purpose                  string `json:"purpose"`
	GoverningLaw             string `json:"governingLaw"`
	Jurisdiction             string `json:"jurisdiction"`
	Modifications            string `json:"modifications"`
}
