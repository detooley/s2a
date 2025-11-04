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
