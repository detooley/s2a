package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// TSDR UTILITIES
func query(sns string) string {
	url := "https://tsdrapi.uspto.gov/ts/cd/caseMultiStatus/sn?ids=" + sns
	// Create a new HTTP client
	client := &http.Client{}
	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}
	// Add headers to the request
	req.Header.Add("USPTO-API-KEY", "ol9RJ5mOfoLwRPcCwSD98V9W2H6N7tYS")
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()
	// Handle the response
	fmt.Println("Status:", resp.Status)
	// Put an error handler here
	// ... process response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error processing response:", err)
		return ""
	}
	//fmt.Println(string(body))
	return string(body)
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

// Unmarshall the JSON file
func extract(jsonData string) file {
	//var data map[string]interface{}
	var fileInfo file
	err := json.Unmarshal([]byte(jsonData), &fileInfo)
	if err != nil {
		fmt.Printf("could not unmarshal json: %s\n", err)
		return fileInfo
	}
	return fileInfo
}

// Parse user input
func parseInput(sns string) []string {
	sns = strings.TrimSpace(sns)
	res := regexp.MustCompile(`,|\s+`)
	sns = res.ReplaceAllString(sns, " ")
	sns = res.ReplaceAllString(sns, ",")
	snsSlice := strings.Split(sns, ",")
	return snsSlice
}

// Retrieve file from API and return Go variable
func GetFile(sns string) trademarks {
	// Query the database
	jsonData := query(sns)
	// Extract JSON to a structure
	fileInfo := extract(jsonData)
	// Clean up variable names and return
	file := fileInfo.TransactionList[0].Trademarks[0]
	return file
}

// WEB SERVER
// Root
func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /request\n")
	io.WriteString(w, "go.s2a.app\n")
}

// tsdr
func getTsdr(w http.ResponseWriter, r *http.Request) {
	fmt.Print("got /tsdr\n")
	// Extract serial numbers
	sns := r.URL.Query().Get("sns")
	snsSlice := parseInput(sns)
	// Query tsdr
	var files []trademarks
	for _, sn := range snsSlice {
		files = append(files, GetFile(sn))
	}
	fmt.Print(files)
	// Set template function
	funcMap := template.FuncMap{
		"dec": func(i int) int { return i - 1 },
		"slz": func(i string) string { return strings.TrimLeft(i, "0") },
	}
	// Parse template
	tmplFile := "tsdr.html"
	tmpl, err := template.New(tmplFile).Funcs(funcMap).ParseFiles(tmplFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Execute template
	err = tmpl.Execute(w, files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//	panic(err)
	}
}

// MAIN
func main() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/tsdr", getTsdr)
	err := http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
