package tsdrstruct

// Define the structure of the JSON file
type file struct {
	TransactionList []transactionList `json:"transactionList"`
	Oversized       bool              `json:"oversized"`
	Size            int               `json:"size"`
}
type transactionList struct {
	Trademarks []trademarks `json:"trademarks"`
	SearchId   string       `json:"searchId"`
}
type trademarks struct {
	Status status   `json:"status"`
	GsList []gsList `json:"gsList"`
}
type status struct {
	SerialNumber         int    `json:"serialNumber"`
	UsRegistrationNumber string `json:"usRegistrationNumber"`
	MarkElement          string `json:"markElement"`
}
type gsList struct {
	Description    string `json:"description"`
	PrimeClassCode string `json:"primeClassCode"`
}
