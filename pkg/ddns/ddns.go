package ddns

import (
	"fmt"

	"matthias-kutz.com/hetzner-ddns/pkg/dns"
	"matthias-kutz.com/hetzner-ddns/pkg/ip"
)

type Service struct {
	IpProvider  ip.Provider
	DnsProvider dns.Provider
	Parameter   Parameter
}

type Parameter struct {
	ZoneName   string
	RecordName string
	RecordType string
	TTL        int
}

func (service Service) Run() {

	// Check if online
	if !service.IpProvider.IsOnline() {
		fmt.Println("No connection to IP provider")
		return
	}

	// Request IP from ip provider
	ip, err := service.IpProvider.Request()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Request zones from dns provider
	zone, err := service.DnsProvider.RequestZone(service.Parameter.ZoneName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Request record from dns provider
	record, err := service.DnsProvider.RequestRecord(zone, service.Parameter.RecordName, service.Parameter.RecordType)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create record with IP from ip provider to update record at the dns provider
	updateRecord := dns.Record{
		Id:     record.Id,
		ZoneId: zone.Id,
		Name:   service.Parameter.RecordName,
		Type:   service.Parameter.RecordType,
		TTL:    service.Parameter.TTL,
		Value:  ip.Value,
	}

	// Check if record value has to be updated
	if ip.Value == record.Value {
		fmt.Println("No DNS update required for", service.Parameter.ZoneName, "with IP", ip.Value)
		return
	}

	fmt.Println("DNS update required for", service.Parameter.ZoneName, "with current IP", record.Value, "to IP", ip.Value)
	// Update the record at the dns provider
	updateRecord, err = service.DnsProvider.UpdateZoneRecord(zone, updateRecord)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Updated DNS for", service.Parameter.ZoneName, "from IP", record.Value, "to IP", updateRecord.Value)

}
