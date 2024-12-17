package tsdrstruct

// Define the structure of the JSON file
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
	Status Status   `json:"status"`
	GsList []GsList `json:"gsList"`
}
type Status struct {
	SerialNumber         int    `json:"serialNumber"`
	UsRegistrationNumber string `json:"usRegistrationNumber"`
	MarkElement          string `json:"markElement"`
}
type GsList struct {
	Description    string `json:"description"`
	PrimeClassCode string `json:"primeClassCode"`
}
