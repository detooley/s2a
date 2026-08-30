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
	fmt.Printf(search)
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
		return ""
	}
	// Add browser headers to avoid HTTP 406 WAF blocks
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()
	// Handle non-200 responses safely
	fmt.Println("Status:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP Error: %s for query %s\n", resp.Status, search)
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error processing response:", err)
		return ""
	}
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
