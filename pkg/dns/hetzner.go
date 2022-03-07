package dns

import (
	"encoding/json"
	"net/http"
	"net/url"

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

type Hetzner struct {
	ApiToken string
}

type hetznerZonesResponse struct {
	Zones []hetznerZone `json:"zones"`
}

type hetznerZone struct {
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

type hetznerRecordsResponse struct {
	Records []hetznerRecord `json:"records"`
}

type hetznerRecord struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
	ZoneId   string `json:"zone_id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	TTL      int    `json:"ttl"`
}

type hetznerUpdateRequest struct {
	Type   string `json:"type"`
	ZoneId string `json:"zone_id"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	TTL    int    `json:"ttl"`
}

type hetznerUpdateResponse struct {
	Record hetznerRecord `json:"record"`
}

func (h Hetzner) RequestZone(zoneName string) (Zone, error) {

	requestUrl := url.URL{
		Scheme: HttpsScheme,
		Host:   HetznerHost,
		Path:   HetznerZonesPath,
	}

	// Request api
	respBody, err := request.Request(http.MethodGet, requestUrl,
		map[string]string{HetznerAuthApiTokenHeader: h.ApiToken}, []byte{})

	if err != nil {
		return Zone{}, &ZoneNotFoundError{
			ZoneName: zoneName,
		}
	}

	// Unmarshal the response
	var hetznerZonesResponse hetznerZonesResponse
	json.Unmarshal(respBody, &hetznerZonesResponse)

	// Convert the response
	var zones []Zone
	for _, z := range hetznerZonesResponse.Zones {
		n := Zone{
			Id:     z.Id,
			Name:   z.Name,
			Status: z.Status,
		}
		zones = append(zones, n)
	}

	// Find the zone by name
	for _, z := range zones {
		if z.Name == zoneName {
			return z, nil
		}
	}

	return Zone{}, &ZoneNotFoundError{
		ZoneName: zoneName,
	}
}

func (h Hetzner) RequestRecord(zone Zone, recordName string, recordType string) (Record, error) {

	requestUrl := url.URL{
		Scheme:   HttpsScheme,
		Host:     HetznerHost,
		Path:     HetznerRecordsPath,
		RawQuery: HetznerRecordsZoneQueryParam + "=" + zone.Id,
	}

	// Request api
	respBody, err := request.Request(http.MethodGet, requestUrl,
		map[string]string{HetznerAuthApiTokenHeader: h.ApiToken}, []byte{})

	if err != nil {
		return Record{}, &RecordNotFoundError{Zone: zone, RecordName: recordName}
	}

	// Unmashal the response
	var hetznerRecordsResponse hetznerRecordsResponse
	json.Unmarshal(respBody, &hetznerRecordsResponse)

	// Convert the response
	var records []Record
	for _, z := range hetznerRecordsResponse.Records {
		n := Record{
			Id:   z.Id,
			Type: z.Type,
			Name: z.Name,

			Value:  z.Value,
			TTL:    z.TTL,
			ZoneId: zone.Id,
		}
		records = append(records, n)
	}
	// Find the record by name and type
	for _, v := range records {
		if v.Name == recordName && v.Type == recordType {
			return v, nil
		}
	}
	return Record{}, &RecordNotFoundError{
		Zone:       zone,
		RecordName: recordName,
		RecordType: recordType,
	}
}

func (h Hetzner) UpdateZoneRecord(zone Zone, record Record) (Record, error) {

	requestUrl := url.URL{
		Scheme: HttpsScheme,
		Host:   HetznerHost,
		Path:   HetznerRecordsPath + "/" + record.Id,
	}

	hetznerUpdateRequest := hetznerUpdateRequest{
		ZoneId: record.ZoneId,
		Type:   record.Type,
		Name:   record.Name,
		Value:  record.Value,
		TTL:    record.TTL,
	}

	requestBody, _ := json.Marshal(hetznerUpdateRequest)

	respBody, err := request.Request(http.MethodPut, requestUrl, map[string]string{HetznerAuthApiTokenHeader: h.ApiToken,
		HetznerContentTypeHeader: ContentTypeApplicationJson}, requestBody)

	if err != nil {
		return Record{}, &UpdateNotPossibleError{
			Zone:   zone,
			Record: record,
		}
	}

	var hetznerUpdateResponse hetznerUpdateResponse
	json.Unmarshal(respBody, &hetznerUpdateResponse)

	return Record{
		Id:     hetznerUpdateResponse.Record.Id,
		Type:   hetznerUpdateResponse.Record.Type,
		ZoneId: hetznerUpdateResponse.Record.ZoneId,
	}, nil
}
