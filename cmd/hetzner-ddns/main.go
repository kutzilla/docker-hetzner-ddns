package main

import (
	"matthias-kutz.com/hetzner-ddns/pkg/args"
	"matthias-kutz.com/hetzner-ddns/pkg/ddns"
	"matthias-kutz.com/hetzner-ddns/pkg/dns"
	"matthias-kutz.com/hetzner-ddns/pkg/ip"
)

func main() {
	// Set args by cli values or env variables
	var zoneName, apiToken, recordType string
	args.SetArgs(&zoneName, &apiToken, &recordType)

	// Validate args
	args.ValidateArgs(zoneName, apiToken, recordType)

	// Set optional args or set the default values
	var recordName, cronExpression string
	args.SetOptionalArgs(&recordName, &cronExpression)

	dnsProvider := dns.Hetzner{
		ApiToken: apiToken,
	}

	ipProvider := ip.Ipify{}

	ddnsParameter := ddns.Parameter{
		ZoneName:   zoneName,
		RecordName: recordName,
		RecordType: recordType,
		TTL:        0,
	}

	ddnsService := ddns.Service{
		DnsProvider: dnsProvider,
		IpProvider:  ipProvider,
		Parameter:   ddnsParameter,
	}

	ddnsScheduler := ddns.Scheduler{
		CronExpression: cronExpression,
		Service:        ddnsService,
	}

	ddnsScheduler.Start()
}
