package idm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// get enties from the ID Manual
func GetEntries(search string) map[string]any {
	results := query(search)
	extracted := extract(results)
	//fmt.Println(j["docs"].(map[string]any))
	return extracted
}

// query the id manual
func query(search string) string {
	url := "https://idm-tmng.uspto.gov/idm2-services/search/public?" + search
	// Create a new HTTP client
	client := &http.Client{}
	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()
	// Handle the response
	fmt.Println("Status:", resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error processing response:", err)
	}
	//fmt.Println(string(body))
	return string(body)
}

// Unmarshall the JSON file
func extract(jsonData string) map[string]any {
	//var data map[string]interface{}
	var data map[string]any
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Printf("could not unmarshal json: %s\n", err)
	}
	return data
}
