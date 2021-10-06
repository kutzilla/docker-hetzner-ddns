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

type IPInfo struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
	Readme   string `json:"readme"`
}

const (
	OrignRecordName              = "@"
	IPv4RecordType               = "A"
	IPv6RecordType               = "AAAA"
	HttpsScheme                  = "https"
	HetznerAuthApiTokenHeader    = "Auth-API-Token"
	HetznerHost                  = "dns.hetzner.com"
	HetznerZonesPath             = "api/v1/zones"
	HetznerRecordsPath           = "api/v1/records"
	HetznerRecordsZoneQueryParam = "zone_id"
	IPInfoHost                   = "ipinfo.io"
)

func main() {

	zoneName := os.Args[1]
	apiToken := os.Args[2]
	recordType := os.Args[3]

	// Validating given record type
	if recordType != IPv4RecordType && recordType != IPv6RecordType {
		fmt.Println("Given record type does not match", IPv4RecordType, "or", IPv6RecordType)
		os.Exit(0)
	}

	// Request all zones
	fmt.Println("Requesting zone", zoneName)
	zones := requestZones(apiToken)

	// Find zone by the given name
	zone := findZoneByName(zones, zoneName)
	fmt.Println("Found zone:", zone)

	fmt.Println("Requesting records for zone:", zone)
	records := requestZoneRecords(zone, apiToken)
	fmt.Println("Found records:", records)

	fmt.Println("Searching origin record for type", recordType, "in", records)
	originRecord := findOrginRecord(records, recordType)
	fmt.Println("Found origin record", originRecord)

	fmt.Println("Requesting IPInfo")
	ipInfo := requestIpInfo()
	fmt.Println("Found IPInfo:", ipInfo)

	fmt.Println("Updating", recordType, "record for zone", zone)

}

func request(httpMethod string, url url.URL, headers map[string]string) []byte {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(httpMethod, url.String(), nil)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Headers
	for key, element := range headers {
		req.Header.Add(key, element)
	}

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
		Scheme: HttpsScheme,
		Host:   HetznerHost,
		Path:   HetznerZonesPath,
	}

	// Request zones
	respBody := request(http.MethodGet, requestUrl, map[string]string{"Auth-API-Token": apiToken})

	// Unmarshal zones
	var zones Zones
	json.Unmarshal(respBody, &zones)

	return zones
}

func requestZoneRecords(zone Zone, apiToken string) Records {

	requestUrl := url.URL{
		Scheme:   HttpsScheme,
		Host:     HetznerHost,
		Path:     HetznerRecordsPath,
		RawQuery: HetznerRecordsZoneQueryParam + "=" + zone.Id,
	}

	respBody := request(http.MethodGet, requestUrl, map[string]string{"Auth-API-Token": apiToken})

	var records Records
	json.Unmarshal(respBody, &records)

	return records
}

func requestIpInfo() IPInfo {

	requestUrl := url.URL{
		Scheme: HttpsScheme,
		Host:   IPInfoHost,
	}

	respBody := request(http.MethodGet, requestUrl, map[string]string{})

	var ipInfo IPInfo
	json.Unmarshal(respBody, &ipInfo)

	return ipInfo
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

func findOrginRecord(records Records, recordType string) Record {
	var originRecord Record
	for _, v := range records.Record {
		if v.Name == OrignRecordName && v.Type == recordType {
			originRecord = v
		}
	}
	return originRecord
}
