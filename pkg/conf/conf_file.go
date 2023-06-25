package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ConfFile struct {
	ZoneName string   `json:"zoneName"`
	ApiToken string   `json:"apiToken"`
	Records  []record `json:"records"`
}

type record struct {
	RecordType     string
	RecordName     string
	TTL            int
	CronExpression string
}

func ReadFile(path string) DynDnsConf {
	// read the conf file
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// unmarshal the conf file content
	var confFile ConfFile
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &confFile)

	return DynDnsConf{
		DnsConf:      DnsConf{ApiToken: confFile.ApiToken, ZoneName: confFile.ZoneName},
		RecordConf:   RecordConf{RecordType: confFile.Records[0].RecordType, RecordName: confFile.Records[0].RecordName},
		ProviderConf: ProviderConf{IpVersion: DetermineIpVersion(confFile.Records[0].RecordType)},
		CronConf:     CronConf{CronExpression: confFile.Records[0].CronExpression},
	}
}
