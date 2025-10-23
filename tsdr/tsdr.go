package tsdr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	ch "s2a.app/utils/chunk"
)

// Retrieve file from API and return Go variable
func GetFiles(sns string) []trademarks {
	snsp := parseInput(sns)
	// Query the database
	var files []trademarks
	for _, snsc := range snsp {
		jsonData := query(snsc)
		fmt.Printf("%s \n\n", snsc)
		// Extract JSON to a structure
		fileInfo, err := extract(jsonData)
		if err != nil {
		}
		// See if I can figure out a way to do this with append
		for _, file := range fileInfo.TransactionList {
			// Got to find a better way to check if file valid
			filet := file.Trademarks[0]
			files = append(files, filet)
			//fmt.Print("%s\n\n", file)
		}
	}
	return files
}

// Query TSDR
func query(sns string) string {
	url := "https://tsdrapi.uspto.gov/ts/cd/caseMultiStatus/sn?ids=" + sns
	// Create a new HTTP client
	client := &http.Client{}
	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	// Add headers to the request
	req.Header.Add("USPTO-API-KEY", "j37IXbK5pGQmOIX08rTQoP8lSxwFedQP")
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()
	// Handle the response
	fmt.Println("Status:", resp.Status)
	// Put an error handler here
	// ... process response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error processing response:", err)
	}
	//fmt.Println(string(body))
	return string(body)
}

// Unmarshall the JSON file
func extract(jsonData string) (file, error) {
	//var data map[string]interface{}
	var fileInfo file
	err := json.Unmarshal([]byte(jsonData), &fileInfo)
	if err != nil {
		fmt.Printf("could not unmarshal json: %s\n", err)
		return fileInfo, errors.New("could not unmarshal json")
	}
	return fileInfo, nil
}

// Parse user input
func parseInput(sns string) []string {
	sns = strings.TrimSpace(sns)
	res := regexp.MustCompile(`,|\s+`)
	sns = res.ReplaceAllString(sns, " ")
	sns = res.ReplaceAllString(sns, ",")
	sn := strings.Split(sns, ",")
	chunk := ch.Chunk(sn, 25)
	//fmt.Print(chunk)
	return chunk
}

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
	Status  status   `json:"status"`
	Parties parties  `json:"parties"`
	GsList  []gsList `json:"gsList"`
}
type status struct {
	Staff                staff  `json:"staff"`
	SerialNumber         int    `json:"serialNumber"`
	UsRegistrationNumber string `json:"usRegistrationNumber"`
	MarkElement          string `json:"markElement"`
}
type staff struct {
	Examiner examiner `json:"examiner"`
}
type examiner struct {
	Name string `json:"name"`
}
type parties struct {
	OwnerGroups ownerGroups `json:"ownerGroups"`
}
type ownerGroups struct {
	Ten []ten `json:"10"`
}
type ten struct {
	Citizenship citizenship `json:"citizenship"`
}
type citizenship struct {
	StateCountry stateCountry `json:"stateCountry"`
}
type stateCountry struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
type gsList struct {
	Description    string `json:"description"`
	PrimeClassCode string `json:"primeClassCode"`
}
