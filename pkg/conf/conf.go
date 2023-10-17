package conf

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/namsral/flag"
)

const (
	EnvZoneName           = "ZONE_NAME"
	EnvApiToken           = "API_TOKEN"
	EnvRecordType         = "RECORD_TYPE"
	EnvRecordName         = "RECORD_NAME"
	EnvCronExpression     = "CRON_EXPRESSION"
	EnvTimeToLive         = "TTL"
	EnvMultipleDomainMode = "MULTIPLE_DOMAIN_MODE"

	DescZoneName           = "The DNS zone that DDNS updates should be applied to."
	DescApiToken           = "Your Hetzner API token."
	DescRecordType         = "The record type of your zone. If your zone uses an IPv4 address use `A`. Use `AAAA` if it uses an IPv6 address."
	DescRecordName         = "The name of the DNS-record that DDNS updates should be applied to. This could be `sub` if you like to update the subdomain `sub.example.com` of `example.com`. The default value is `@`"
	DescCronExpression     = "The cron expression of the DDNS update interval. The default is every 5 minutes - `*/5 * * * *`"
	DescTimeToLive         = "Time to live of the record"
	DescMultipleDomainMode = "Sets the mode to update a single record or multiple records at once"

	DefaultRecordKey          = "default"
	DefaultRecordName         = "@"
	DefaultRecordType         = "A"
	DefaultCronExpression     = "*/5 * * * *"
	DefaultTimeToLive         = 86400
	DefaultMultipleDomainMode = false

	IPv4           = "IPv4"
	IPv6           = "IPv6"
	IPv6RecordType = "AAAA"
)

type RecordConfig struct {
	mode    bool
	records map[string]*RecordConf
}

type DynDnsConf struct {
	DnsConf      DnsConf
	RecordConf   RecordConfig
	ProviderConf ProviderConf
	CronConf     CronConf
}

type DnsConf struct {
	ApiToken string
	ZoneName string
}

type RecordConf struct {
	RecordType string
	RecordName string
	TTL        int
}

type CronConf struct {
	CronExpression string
}

type ProviderConf struct {
	IpVersion string
}

type ArgumentMissingError struct {
	argumentName string
}

func (e *ArgumentMissingError) Error() string {
	return "The mandatory argument " + e.argumentName + " is missing"
}

func setupRecordConfig(recordConf RecordConfig) {
	flag.BoolVar(&recordConf.mode, EnvMultipleDomainMode, DefaultMultipleDomainMode, DescMultipleDomainMode)

	if recordConf.mode {
		envPrefix := fmt.Sprintf("%s_", EnvRecordName)
		for _, envRecord := range os.Environ() {
			if strings.HasPrefix(envRecord, envPrefix) {
				envKey := strings.Split(envRecord, "=")[0]

				if strings.HasSuffix(envKey, "_TTL") {
					continue
				}

				if _, exists := recordConf.records[envKey]; exists {
					continue
				}

				var record = &RecordConf{
					RecordType: DefaultRecordType, // assume its ipv4, fix later after arg parse if not
				}
				recordConf.records[envKey] = record

				flag.StringVar(&record.RecordName, envKey, DefaultRecordName, DescRecordName)
				flag.IntVar(&record.TTL, fmt.Sprintf("%s_TTL", envKey), DefaultTimeToLive, DescTimeToLive)
			}
		}
	}

	if !recordConf.mode {
		var record = &RecordConf{
			RecordType: DefaultRecordType,
		}
		flag.StringVar(&record.RecordName, EnvRecordName, DefaultRecordName, DescRecordName)
		flag.IntVar(&record.TTL, EnvTimeToLive, DefaultTimeToLive, DescTimeToLive)
		recordConf.records[DefaultRecordKey] = record
	}
}

func setRecordType(recordConf RecordConfig, recordType string) {
	for _, record := range recordConf.records {
		record.RecordType = recordType
	}
}

func Read() DynDnsConf {
	// Mandatory flags
	var zoneName, apiToken, recordType string
	flag.StringVar(&zoneName, EnvZoneName, zoneName, DescZoneName)
	flag.StringVar(&apiToken, EnvApiToken, apiToken, DescApiToken)
	flag.StringVar(&recordType, EnvRecordType, recordType, DescRecordType)

	records := RecordConfig{
		mode:    DefaultMultipleDomainMode,
		records: make(map[string]*RecordConf),
	}
	setupRecordConfig(records)

	var cronExpression = DefaultCronExpression
	flag.StringVar(&cronExpression, EnvCronExpression, cronExpression, DescCronExpression)

	// Parse flags
	flag.Parse()

	// Computed confs
	var ipVersion = IPv4
	if recordType == IPv6RecordType {
		ipVersion = IPv6
		setRecordType(records, recordType)
	}

	dynDnsConf := DynDnsConf{
		DnsConf:    DnsConf{ApiToken: apiToken, ZoneName: zoneName},
		RecordConf: records,
		ProviderConf: ProviderConf{
			IpVersion: ipVersion,
		},
		CronConf: CronConf{CronExpression: cronExpression},
	}

	validatedConf, err := validate(dynDnsConf)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	return validatedConf
}

func PrintUsage() {
	flag.Usage()
}

func validate(dynDnsConf DynDnsConf) (DynDnsConf, error) {
	// Check api token
	if dynDnsConf.DnsConf.ApiToken == "" {
		return dynDnsConf, &ArgumentMissingError{
			argumentName: EnvApiToken,
		}
	}

	// Check zone name
	if dynDnsConf.DnsConf.ZoneName == "" {
		return dynDnsConf, &ArgumentMissingError{
			argumentName: EnvZoneName,
		}
	}

	// Check record type
	for _, record := range dynDnsConf.RecordConf.records {
		if record.RecordType == "" {
			return dynDnsConf, &ArgumentMissingError{
				argumentName: EnvRecordType,
			}
		}
	}

	return dynDnsConf, nil
}
