package conf

import "github.com/namsral/flag"

const (
	EnvZoneName       = "ZONE_NAME"
	EnvApiToken       = "API_TOKEN"
	EnvRecordType     = "RECORD_TYPE"
	EnvRecordName     = "RECORD_NAME"
	EnvCronExpression = "CRON_EXPRESSION"
	EnvTimeToLive     = "TTL"

	DescZoneName       = "Name of the zone"
	DescApiToken       = ""
	DescRecordType     = ""
	DescRecordName     = ""
	DescCronExpression = ""
	DescTimeToLive     = ""

	DefaultRecordName     = "@"
	DefaultCronExpression = "*/5 * * * *"
	DefaultTimeToLive     = 86400
)

type DynDnsConf struct {
	DnsConf    DnsConf
	RecordConf RecordConf
	CronConf   CronConf
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

type ArgumentMissingError struct {
	argumentName string
}

func (e *ArgumentMissingError) Error() string {
	return "The argument " + e.argumentName + " is missing"
}

func Read() (DynDnsConf, error) {
	// Mandatory flags
	var zoneName, apiToken, recordType string
	flag.StringVar(&zoneName, EnvZoneName, zoneName, DescZoneName)
	flag.StringVar(&apiToken, EnvApiToken, apiToken, DescApiToken)
	flag.StringVar(&recordType, EnvRecordType, recordType, DescRecordType)

	// Optional flags
	var recordName = DefaultRecordName
	flag.StringVar(&recordName, EnvRecordName, recordName, DescRecordName)
	var cronExpression = DefaultCronExpression
	flag.StringVar(&cronExpression, EnvCronExpression, cronExpression, DescCronExpression)
	var ttl = DefaultTimeToLive
	flag.IntVar(&ttl, EnvTimeToLive, ttl, DescTimeToLive)

	// Parse flags
	flag.Parse()

	dynDnsConf := DynDnsConf{
		DnsConf: DnsConf{ApiToken: apiToken, ZoneName: zoneName},
		RecordConf: RecordConf{
			RecordType: recordType,
			RecordName: recordName,
			TTL:        ttl,
		},
		CronConf: CronConf{CronExpression: cronExpression},
	}

	return validate(dynDnsConf)
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
	if dynDnsConf.RecordConf.RecordType == "" {
		return dynDnsConf, &ArgumentMissingError{
			argumentName: EnvRecordType,
		}
	}

	return dynDnsConf, nil
}
