package ddns

import (
	"fmt"

	"github.com/robfig/cron/v3"
	"matthias-kutz.com/hetzner-ddns/pkg/args"
	"matthias-kutz.com/hetzner-ddns/pkg/hetzner"
	"matthias-kutz.com/hetzner-ddns/pkg/ipify"
)

type DnsParameter struct {
	zone       hetzner.Zone
	apiToken   string
	recordName string
	recordType string
}

func CreateDynDnsParameters(zone hetzner.Zone, apiToken string, recordType string, recordName string) DnsParameter {
	return DnsParameter{
		zone:       zone,
		apiToken:   apiToken,
		recordType: recordType,
		recordName: recordName,
	}
}

func CreateHetznerDynDnsCronJobBy(dnsParameter DnsParameter) cron.FuncJob {
	return CreateHetznerDynDnsCronJob(dnsParameter.zone, dnsParameter.apiToken, dnsParameter.recordType, dnsParameter.recordName)
}

func CreateHetznerDynDnsCronJob(zone hetzner.Zone, apiToken string, recordType string, recordName string) cron.FuncJob {
	return cron.FuncJob(func() {
		if ipify.IsOnline() {
			records := hetzner.RequestZoneRecords(zone, apiToken)
			record, _ := hetzner.FindDnsRecord(records, recordType, recordName)
			ipify := ipify.RequestIpify()
			fmt.Println("Current public IP is:", ipify.IP)

			var domain string
			if recordName == args.DefaultRecordName {
				domain = zone.Name
			} else {
				domain = recordName + "." + zone.Name
			}

			if record.Value == ipify.IP {
				fmt.Println("No DNS update required for", recordName+
					domain, "to IP", record.Value)
			} else {
				fmt.Println("DNS update required for", domain,
					"with IP", record.Value)
				updatedDnsRecord := hetzner.UpdateDnsRecord(record, ipify, apiToken)
				fmt.Println("Updated DNS for", domain, "from IP",
					record.Value, "to IP", updatedDnsRecord.Value)
			}
		}
	})
}
