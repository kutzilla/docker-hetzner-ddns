package hetzner

import (
	"encoding/json"
	"net/http"
	"net/url"

	"matthias-kutz.com/hetzner-ddns/pkg/ipify"
	"matthias-kutz.com/hetzner-ddns/pkg/request"
)

const (
	HttpsScheme                  = "https"
	HetznerHost                  = "dns.hetzner.com"
	HetznerZonesPath             = "api/v1/zones"
	HetznerRecordsPath           = "api/v1/records"
	HetznerRecordsZoneQueryParam = "zone_id"
	HetznerAuthApiTokenHeader    = "Auth-API-Token"
	HetznerContentTypeHeader     = "Content-Type"

	DefaultRecordName = "@"

	ContentTypeApplicationJson = "application/json"
)

type NotFoundError struct {
}

func (e *NotFoundError) Error() string {
	return "The zone was not found"
}

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

func RequestZones(apiToken string) Zones {

	requestUrl := url.URL{
		Scheme: HttpsScheme,
		Host:   HetznerHost,
		Path:   HetznerZonesPath,
	}

	// Request zones
	respBody := request.Request(http.MethodGet, requestUrl,
		map[string]string{HetznerAuthApiTokenHeader: apiToken}, []byte{})

	// Unmarshal zones
	var zones Zones
	json.Unmarshal(respBody, &zones)

	return zones
}

func FindZoneByName(zones Zones, zoneName string) (Zone, error) {
	var foundZone Zone
	for _, v := range zones.Zone {
		if v.Name == zoneName {
			foundZone = v
			return foundZone, nil
		}
	}
	return foundZone, &NotFoundError{}
}

func RequestZoneRecords(zone Zone, apiToken string) Records {

	requestUrl := url.URL{
		Scheme:   HttpsScheme,
		Host:     HetznerHost,
		Path:     HetznerRecordsPath,
		RawQuery: HetznerRecordsZoneQueryParam + "=" + zone.Id,
	}

	respBody := request.Request(http.MethodGet, requestUrl,
		map[string]string{HetznerAuthApiTokenHeader: apiToken}, []byte{})

	var records Records
	json.Unmarshal(respBody, &records)

	return records
}

func UpdateDnsRecord(dnsRecord Record, ipify ipify.Ipify, apiToken string) Record {

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

	respBody := request.Request(http.MethodPut, requestUrl, map[string]string{HetznerAuthApiTokenHeader: apiToken,
		HetznerContentTypeHeader: ContentTypeApplicationJson}, requestBody)

	var recordUpdateResponse RecordUpdateResponse
	json.Unmarshal(respBody, &recordUpdateResponse)

	return recordUpdateResponse.Record
}

func FindDnsRecord(records Records, recordType string, dnsRecordName string) (Record, error) {
	var dnsRecord Record
	for _, v := range records.Record {
		if v.Name == dnsRecordName && v.Type == recordType {
			dnsRecord = v
			return dnsRecord, nil
		}
	}
	return dnsRecord, &NotFoundError{}
}
