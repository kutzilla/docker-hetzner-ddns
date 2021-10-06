package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Zones struct {
	Zone []Zone `json:"zones"`
}

type Zone struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func main() {

	zoneName := os.Args[1]
	apiToken := os.Args[2]
	recordType := os.Args[3]

	fmt.Println("Searching:", zoneName)

	zone := findZoneByName(zoneName, apiToken)

	fmt.Println("Updating", recordType, "record for zone", zone.Name)

}

func findZoneByName(zoneName string, apiToken string) Zone {

	// Request zones by API Token
	zones := requestZones(apiToken)

	// Find zone instance for given zone name
	zone := findZoneInZones(zoneName, zones)

	fmt.Println("Found zone for", zoneName, "with id", zone.Id)

	return zone
}

func requestZones(apiToken string) Zones {
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

	// Unmarshal zones
	var zones Zones
	json.Unmarshal(respBody, &zones)

	return zones
}

func findZoneInZones(zoneName string, zones Zones) Zone {
	var foundZone Zone
	for _, v := range zones.Zone {
		if v.Name == zoneName {
			foundZone = v
			return foundZone
		}
	}
	return foundZone
}
