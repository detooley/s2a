package structs

import (
	"time"
)

// Define the structure of the Pulls Table
type Pulls struct {
	SerNo    int       `db:"ser_no"`
	Mark     string    `db:"mark"`
	Country  string    `db:"country"`
	Examiner string    `db:"examiner"`
	PullDate time.Time `db:"pull_date"`
}

// Define the structure of the Id Manual Table
type IdManual struct {
	Id                 int       `db:"id"`
	IdTx               string    `db:"id_tx"`
	ClassId            string    `db:"class_id"`
	DescriptionTx      string    `db:"description_tx"`
	Notes              string    `db:"notes"`
	Version            string    `db:"version"`
	Tm5                string    `db:"TM5"`
	BeginEffectiveDate time.Time `db:"begin_effective_date"`
	Status             string    `db:"status"`
}

// Define the structure of search results
type IdmPageData struct {
	Search   string
	OrderBy  OrderBy
	NumFound int
	Time     time.Duration
	Options  Options
	Results  []IdManual
}

// Define the structure of IDM page data
type Options struct {
	ShowId     string
	ShowStatus string
	ShowDate   string
	ShowVer    string
	ShowTm5    string
	Submitted  string
}

// Define the structure for Order By
type OrderBy struct {
	Element   string
	Direction string
}

// Define the structure of the TSDR JSON file
type File struct {
	TransactionList []TransactionList `json:"transactionList"`
	Oversized       bool              `json:"oversized"`
	Size            int               `json:"size"`
}
type TransactionList struct {
	Trademarks []Trademarks `json:"trademarks"`
	SearchId   string       `json:"searchId"`
}
type Trademarks struct {
	Status  Status   `json:"status"`
	Parties Parties  `json:"parties"`
	GsList  []GsList `json:"gsList"`
}
type Status struct {
	Staff                Staff  `json:"staff"`
	SerialNumber         int    `json:"serialNumber"`
	UsRegistrationNumber string `json:"usRegistrationNumber"`
	MarkElement          string `json:"markElement"`
}
type Staff struct {
	Examiner Examiner `json:"examiner"`
}
type Examiner struct {
	Name string `json:"name"`
}
type Parties struct {
	OwnerGroups OwnerGroups `json:"ownerGroups"`
}
type OwnerGroups struct {
	Ten []Ten `json:"10"`
}
type Ten struct {
	Citizenship Citizenship `json:"citizenship"`
}
type Citizenship struct {
	StateCountry StateCountry `json:"stateCountry"`
}
type StateCountry struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
type GsList struct {
	Description    string `json:"description"`
	PrimeClassCode string `json:"primeClassCode"`
}
