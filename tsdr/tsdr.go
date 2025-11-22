package tsdr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"s2a.app/utils/structs"
)

// Retrieve file from API and return Go variable
func GetFiles(ids string) []structs.Trademarks {
	idsp := parseInput(ids)
	// Query the database
	var files []structs.Trademarks
	for _, id := range idsp {
		jsonData := query(id)
		fmt.Printf("%s \n\n", id)
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
func query(id string) string {
	var sorr string
	if len(id) == 8 {
		sorr = "sn"
	} else {
		sorr = "rn"
	}
	url := "https://tsdrapi.uspto.gov/ts/cd/caseMultiStatus/" + sorr + "?ids=" + id
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
func extract(jsonData string) (structs.File, error) {
	//var data map[string]interface{}
	var fileInfo structs.File
	err := json.Unmarshal([]byte(jsonData), &fileInfo)
	if err != nil {
		fmt.Printf("could not unmarshal json: %s\n", err)
		return fileInfo, errors.New("could not unmarshal json")
	}
	return fileInfo, nil
}

// Parse user input
func parseInput(ids string) []string {
	ids = strings.TrimSpace(ids)
	res := regexp.MustCompile(`,|\s+`)
	ids = res.ReplaceAllString(ids, " ")
	ids = res.ReplaceAllString(ids, ",")
	idList := strings.Split(ids, ",")
	//chunk := ch.Chunk(idList, 25)
	//fmt.Print(chunk)
	return idList
}
