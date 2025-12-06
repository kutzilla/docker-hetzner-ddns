package ddns

import (
	"log"
	"matthias-kutz.com/hetzner-ddns/pkg/conf"

	"matthias-kutz.com/hetzner-ddns/pkg/dns"
	"matthias-kutz.com/hetzner-ddns/pkg/ip"
)

type Service struct {
	IpProvider  ip.Provider
	DnsProvider dns.Provider
	Parameter   Parameter
}

type Parameter struct {
	ZoneName string
	Records  []Record
}

type Record struct {
	RecordName string
	RecordType string
	TTL        int
}

func ConvertFromConfig(config conf.RecordConfig) []Record {
	var recordList []Record

	for _, recordConf := range config {
		recordList = append(recordList, Record{
			RecordName: recordConf.RecordName,
			RecordType: recordConf.RecordType,
			TTL:        recordConf.TTL,
		})
	}

	return recordList
}

func (service Service) Run() {

	// Check if online
	if !service.IpProvider.IsOnline() {
		log.Println("No connection to IP provider")
		return
	}

	// Request IP from ip provider
	ip, err := service.IpProvider.Request()
	if err != nil {
		log.Println(err)
		return
	}

	// Request zones from dns provider
	zone, err := service.DnsProvider.RequestZone(service.Parameter.ZoneName)
	if err != nil {
		log.Println(err)
		return
	}

	for _, dnsRecord := range service.Parameter.Records {
		// Request record from dns provider
		record, err := service.DnsProvider.RequestRecord(zone, dnsRecord.RecordName, dnsRecord.RecordType)
		if err != nil {
			log.Println(err)
			continue
		}

		// Create record with IP from ip provider to update record at the dns provider
		updateRecord := dns.Record{
			Id:     record.Id,
			ZoneId: zone.Id,
			Name:   dnsRecord.RecordName,
			Type:   dnsRecord.RecordType,
			TTL:    dnsRecord.TTL,
			Value:  ip.Value,
		}

		// Check if record value has to be updated
		if ip.Value == record.Value {
			//log.Printf("no update required for %s.%s [%s]\n", dnsRecord.RecordName, service.Parameter.ZoneName, ip.Value)
			continue
		}

		log.Printf("updating %s.%s [%s] to %s\n", dnsRecord.RecordName, service.Parameter.ZoneName, record.Value, ip.Value)
		// Update the record at the dns provider
		updateRecord, err = service.DnsProvider.UpdateZoneRecord(zone, updateRecord)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("updated %s.%s [%s]\n", dnsRecord.RecordName, service.Parameter.ZoneName, updateRecord.Value)
	}
}
