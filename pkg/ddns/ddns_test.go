package ddns

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"matthias-kutz.com/hetzner-ddns/pkg/dns"
	"matthias-kutz.com/hetzner-ddns/pkg/ip"
)

type MockIpProvider struct {
	mockValue  string
	mockSource string
}

type MockDnsProvider struct {
	mockZone   dns.Zone
	mockRecord dns.Record
}

func TestRun(t *testing.T) {
	assert := assert.New(t)

	ipProvider := &MockIpProvider{
		mockValue:  "1.2.3.4",
		mockSource: "MockIpProvider",
	}

	dnsProvider := &MockDnsProvider{
		mockZone: dns.Zone{
			Id:     "12345",
			Name:   "example.com",
			Status: "verified",
		},
		mockRecord: dns.Record{
			Id:     "67890",
			Name:   "@",
			Type:   "A",
			Value:  "4.3.2.1",
			ZoneId: "12345",
			TTL:    86400,
		},
	}

	ddnsService := Service{
		IpProvider:  ipProvider,
		DnsProvider: dnsProvider,
		Parameter: Parameter{
			ZoneName:   "example.com",
			RecordName: "@",
			RecordType: "A",
		},
	}

	// assert IP value is 4.3.2.1 before
	assert.Equal("4.3.2.1", dnsProvider.mockRecord.Value)
	ddnsService.Run()
	// assert IP value is value of mock IP Provider (1.2.3.4)
	assert.Equal("1.2.3.4", dnsProvider.mockRecord.Value)

}

func (i *MockIpProvider) Request() (ip.IP, error) {
	return ip.IP{
		Value:  i.mockValue,
		Source: i.mockSource,
	}, nil
}

func (i *MockIpProvider) IsOnline() bool {
	return true
}

func (d *MockDnsProvider) RequestZone(zoneName string) (dns.Zone, error) {
	return d.mockZone, nil
}

func (d *MockDnsProvider) RequestRecord(zone dns.Zone, recordName string, recordType string) (dns.Record, error) {
	return d.mockRecord, nil
}

func (d *MockDnsProvider) UpdateZoneRecord(zone dns.Zone, record dns.Record) (dns.Record, error) {
	d.mockRecord = dns.Record{
		Id:     record.Id,
		Name:   record.Name,
		Type:   record.Type,
		Value:  record.Value,
		ZoneId: record.ZoneId,
		TTL:    record.TTL,
	}
	return d.mockRecord, nil

}
