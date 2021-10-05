package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Zone struct {
	id   string `json:"id"`
	name string `json:"name"`
}

func main() {

	zoneName := os.Args[1]
	apiToken := os.Args[2]

	fmt.Println("Finding:", zoneName)

	zone := findZoneByName(zoneName, apiToken)
	fmt.Println("Found:", zone)
}

func findZoneByName(zoneName string, apiToken string) Zone {

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", "https://dns.hetzner.com/api/v1/zones", nil)

	// Headers
	req.Header.Add("Auth-API-Token", apiToken)

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response body:", string(respBody))

	var zones []Zone
	json.Unmarshal(respBody, &zones)
	fmt.Println("Zones:", zones)

	if len(zones) > 0 {
		return zones[0]
	} else {
		return Zone{}
	}
}
