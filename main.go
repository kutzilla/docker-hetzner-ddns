package main

import (
	"fmt"

	"matthias-kutz.com/hetzner-ddns/pkg/args"
	"matthias-kutz.com/hetzner-ddns/pkg/cron"
	"matthias-kutz.com/hetzner-ddns/pkg/ddns"
	"matthias-kutz.com/hetzner-ddns/pkg/hetzner"
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

	// Request zones
	zones := hetzner.RequestZones(apiToken)
	// Find zone by the given name
	fmt.Println("Requesting zone:", zoneName)
	zone, _ := hetzner.FindZoneByName(zones, zoneName)

	// Create the DNS Parameter
	dnsParameter := ddns.CreateDynDnsParameters(zone, apiToken, recordType, recordName)

	// Create the Hetzner DynDNS Cron Job
	hetznerDynDnsJob := ddns.CreateHetznerDynDnsCronJobBy(dnsParameter)

	// Start the Cron Job with the expression
	cron.StartCronScheduler(cronExpression, hetznerDynDnsJob)
}
