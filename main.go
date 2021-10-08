package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"

	"github.com/robfig/cron/v3"
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
	Id       string `json:"id"`
	Type     string `json:"type"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
	ZoneId   string `json:"zone_id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	TTL      int    `json:"ttl"`
}

type RecordUpdateRequest struct {
	Type   string `json:"type"`
	ZoneId string `json:"zone_id"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	TTL    int    `json:"ttl"`
}

type RecordUpdateResponse struct {
	Record Record `json:"record"`
}

type Ipify struct {
	IP string `json:"ip"`
}

const (
	DnsRecordName = "@"

	HttpsScheme = "https"

	EnvZoneName   = "ZONE_NAME"
	EnvApiToken   = "API_TOKEN"
	EnvRecordType = "RECORD_TYPE"

	IPv4RecordType = "A"
	IPv6RecordType = "AAAA"

	HetznerHost                  = "dns.hetzner.com"
	HetznerZonesPath             = "api/v1/zones"
	HetznerRecordsPath           = "api/v1/records"
	HetznerRecordsZoneQueryParam = "zone_id"
	HetznerAuthApiTokenHeader    = "Auth-API-Token"
	HetznerContentTypeHeader     = "Content-Type"

	ContentTypeApplicationJson = "application/json"

	IpifyHost                 = "api.ipify.org"
	IpifyFormatQueryParam     = "format"
	IpifyQueryParamFormatJson = "json"

	DefaultDnsUpdateCronExpression = "*/5 * * * *"
)

func main() {
	// Set args by cli values or env variables
	var zoneName, apiToken, recordType string
	setArgs(&zoneName, &apiToken, &recordType)

	// Validate args
	validateArgs(zoneName, apiToken, recordType)

	// Request zones
	zones := requestZones(apiToken)
	// Find zone by the given name
	fmt.Println("Requesting zone:", zoneName)
	zone := findZoneByName(zones, zoneName)

	// Create the Hetzner DynDNS Cron Job
	hetznerDynDnsJob := createHetznerDynDnsCronJob(zone, apiToken, recordType)

	// Start the Cron Job with the expression
	startCronScheduler(DefaultDnsUpdateCronExpression, hetznerDynDnsJob)
}

func startCronScheduler(cronExpression string, cronJob cron.FuncJob) {
	waitUntilStopped := func() {
		var endWaiter sync.WaitGroup
		endWaiter.Add(1)
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt)
		go func() {
			<-signalChannel
			endWaiter.Done()
		}()
		endWaiter.Wait()
	}
	cron := cron.New()
	cron.AddJob(cronExpression, cronJob)
	cron.Start()
	fmt.Println("Started DynDNS")
	waitUntilStopped()
	fmt.Println("Stopped DynDNS")
	cron.Stop()
}

func createHetznerDynDnsCronJob(zone Zone, apiToken string, recordType string) cron.FuncJob {
	return cron.FuncJob(func() {
		if isOnline() {
			records := requestZoneRecords(zone, apiToken)
			dnsRecord := findDnsRecord(records, recordType)
			ipify := requestIpify()
			fmt.Println("Current public IP is:", ipify.IP)

			if dnsRecord.Value == ipify.IP {
				fmt.Println("No DNS update required for",
					zone.Name, "to IP", dnsRecord.Value)
			} else {
				fmt.Println("DNS update required for", zone.Name,
					"with IP", dnsRecord.Value)
				updatedDnsRecord := updateDnsRecord(dnsRecord, ipify, apiToken)
				fmt.Println("Updated DNS for", zone.Name, "from IP",
					dnsRecord.Value, "to IP", updatedDnsRecord.Value)
			}
		}
	})
}

func isOnline() bool {
	geoIpUrl := url.URL{
		Scheme: HttpsScheme,
		Host:   IpifyHost,
	}

	_, err := http.Get(geoIpUrl.String())
	return err == nil
}

func setArgs(zoneName *string, apiToken *string, recordType *string) {
	if len(os.Args) > 3 {
		*zoneName = os.Args[1]
		*apiToken = os.Args[2]
		*recordType = os.Args[3]
	} else {
		*zoneName = os.Getenv(EnvZoneName)
		*apiToken = os.Getenv(EnvApiToken)
		*recordType = os.Getenv(EnvRecordType)
	}
}

func validateArgs(zoneName string, apiToken string, recordType string) {
	if zoneName == "" || apiToken == "" || recordType == "" {
		fmt.Println("You need to set the environment variables",
			EnvZoneName, ",", EnvApiToken, "and", EnvRecordType,
			"or provide them as args in the shown order")
		os.Exit(1)
	}

	// Validating given record type
	if recordType != IPv4RecordType && recordType != IPv6RecordType {
		fmt.Println("Given record type does not match",
			IPv4RecordType, "or", IPv6RecordType)
		os.Exit(1)
	}
}

func request(httpMethod string, url url.URL, headers map[string]string, body []byte) []byte {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(httpMethod, url.String(), bytes.NewBuffer(body))

	if err != nil {
		fmt.Println("Failure : ", err)
		return []byte{}
	}

	// Headers
	for key, element := range headers {
		req.Header.Add(key, element)
	}

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
		return []byte{}
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
	respBody := request(http.MethodGet, requestUrl,
		map[string]string{HetznerAuthApiTokenHeader: apiToken}, []byte{})

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

	respBody := request(http.MethodGet, requestUrl,
		map[string]string{HetznerAuthApiTokenHeader: apiToken}, []byte{})

	var records Records
	json.Unmarshal(respBody, &records)

	return records
}

func requestIpify() Ipify {
	requestUrl := url.URL{
		Scheme:   HttpsScheme,
		Host:     IpifyHost,
		RawQuery: IpifyFormatQueryParam + "=" + IpifyQueryParamFormatJson,
	}

	respBody := request(http.MethodGet, requestUrl,
		map[string]string{}, []byte{})

	var ipify Ipify
	json.Unmarshal(respBody, &ipify)

	return ipify
}

func updateDnsRecord(dnsRecord Record, ipify Ipify, apiToken string) Record {

	requestUrl := url.URL{
		Scheme: HttpsScheme,
		Host:   HetznerHost,
		Path:   HetznerRecordsPath + "/" + dnsRecord.Id,
	}

	// Creating new DNS record with IP from ipify.org
	requestRecordUpdate := RecordUpdateRequest{
		ZoneId: dnsRecord.ZoneId,
		Type:   dnsRecord.Type,
		Name:   dnsRecord.Name,
		Value:  ipify.IP,
		TTL:    dnsRecord.TTL,
	}

	requestBody, _ := json.Marshal(requestRecordUpdate)

	respBody := request(http.MethodPut, requestUrl, map[string]string{HetznerAuthApiTokenHeader: apiToken, HetznerContentTypeHeader: ContentTypeApplicationJson}, requestBody)

	var recordUpdateResponse RecordUpdateResponse
	json.Unmarshal(respBody, &recordUpdateResponse)

	return recordUpdateResponse.Record
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

func findDnsRecord(records Records, recordType string) Record {
	var dnsRecord Record
	for _, v := range records.Record {
		if v.Name == DnsRecordName && v.Type == recordType {
			dnsRecord = v
		}
	}
	return dnsRecord
}
