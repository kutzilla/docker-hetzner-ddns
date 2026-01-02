package main

import (
	"github.com/kutzilla/hetzner-cloud-ddns/pkg/conf"
	"github.com/kutzilla/hetzner-cloud-ddns/pkg/ddns"
	"github.com/kutzilla/hetzner-cloud-ddns/pkg/dns"
	"github.com/kutzilla/hetzner-cloud-ddns/pkg/ip"
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
