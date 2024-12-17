package tsdr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"tsdr/tsdrstruct"
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

// Unmarshall the JSON file
func extract(jsonData string) tsdrstruct.File {
	//var data map[string]interface{}
	var fileInfo tsdrstruct.File
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
func GetFile(sns string) tsdrstruct.Trademarks {
	// Query the database
	jsonData := query(sns)
	// Extract JSON to a structure
	fileInfo := extract(jsonData)
	// Clean up variable names and return
	file := fileInfo.TransactionList[0].Trademarks[0]
	return file
}
