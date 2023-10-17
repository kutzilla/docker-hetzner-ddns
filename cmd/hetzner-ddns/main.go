package main

import (
	"matthias-kutz.com/hetzner-ddns/pkg/conf"
	"matthias-kutz.com/hetzner-ddns/pkg/ddns"
	"matthias-kutz.com/hetzner-ddns/pkg/dns"
	"matthias-kutz.com/hetzner-ddns/pkg/ip"
)

func main() {
	dynDnsConf := conf.Read()

	dnsProvider := dns.Hetzner{
		ApiToken: dynDnsConf.DnsConf.ApiToken,
	}

	ipProvider := ip.Ipify{
		IpVersion: dynDnsConf.ProviderConf.IpVersion,
	}

	ddnsParameter := ddns.Parameter{
		ZoneName: dynDnsConf.DnsConf.ZoneName,
		Records:  ddns.ConvertFromConfig(dynDnsConf.RecordConf),
	}

	ddnsService := ddns.Service{
		DnsProvider: dnsProvider,
		IpProvider:  ipProvider,
		Parameter:   ddnsParameter,
	}

	ddnsScheduler := ddns.Scheduler{
		CronExpression: dynDnsConf.CronConf.CronExpression,
		Service:        ddnsService,
	}

	ddnsScheduler.Start()
}
