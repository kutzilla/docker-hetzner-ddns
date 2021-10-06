package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type Zones struct {
	Zone []Zone `json:"zones"`
}

type Zone struct {
	Id             string `json:"id"`
	Created        string `json:"created"`
	Modified       string `json:"modified"`
	LegacyDnsHost  string `json:"legacy_dns_host"`
	Owner          string `json:"owner"`
	Name           string `json:"name"`
	Paused         bool   `json:"paused"`
	Permission     string `json:"permission"`
	Project        string `json:"project"`
	Registrar      string `json:"registrar"`
	Status         string `json:"status"`
	TTL            int    `json:"ttl"`
	Verified       string `json:"verified"`
	RecordsCount   int    `json:"records_count"`
	IsSecondaryDns bool   `json:"is_secondary_dns"`
}

type Records struct {
	Record []Record `json:"records"`
}

type Record struct {
	Id   string `json:"id"`
	Type string `json:"type"`

	Created  string `json:"created"`
	Modified string `json:"modified"`
	ZoneId   string `json:"zone_id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	TTL      int    `json:"ttl"`
}

func main() {

	zoneName := os.Args[1]
	apiToken := os.Args[2]
	recordType := os.Args[3]

	// Request all zones
	fmt.Println("Requesting zone", zoneName)
	zones := requestZones(apiToken)

	// Find zone by the given name
	zone := findZoneByName(zones, zoneName)
	fmt.Println("Found zone:", zone)

	fmt.Println("Requesting records for zone:", zone)
	records := requestZoneRecords(zone, apiToken)
	fmt.Println("Found records:", records)

	fmt.Println("Updating", recordType, "record for zone", zone)

}

func request(httpMethod string, url url.URL, apiToken string) []byte {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(httpMethod, url.String(), nil)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Headers
	req.Header.Add("Auth-API-Token", apiToken)

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	return respBody
}

func requestZones(apiToken string) Zones {

	requestUrl := url.URL{
		Scheme: "https",
		Host:   "dns.hetzner.com",
		Path:   "api/v1/zones",
	}

	// Request zones
	respBody := request("GET", requestUrl, apiToken)

	// Unmarshal zones
	var zones Zones
	json.Unmarshal(respBody, &zones)

	return zones
}

func requestZoneRecords(zone Zone, apiToken string) Records {

	requestUrl := url.URL{
		Scheme:   "https",
		Host:     "dns.hetzner.com",
		Path:     "api/v1/records",
		RawQuery: "zone_id=" + zone.Id,
	}

	respBody := request("GET", requestUrl, apiToken)

	var records Records
	json.Unmarshal(respBody, &records)

	return records
}

func findZoneByName(zones Zones, zoneName string) Zone {
	var foundZone Zone
	for _, v := range zones.Zone {
		if v.Name == zoneName {
			foundZone = v
		}
	}
	return foundZone
}
