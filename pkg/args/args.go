package args

import (
	"fmt"
	"os"
)

const (
	IPv4RecordType = "A"
	IPv6RecordType = "AAAA"

	EnvZoneName       = "ZONE_NAME"
	EnvApiToken       = "API_TOKEN"
	EnvRecordType     = "RECORD_TYPE"
	EnvRecordName     = "RECORD_NAME"
	EnvCronExpression = "CRON_EXPRESSION"

	DefaultRecordName     = "@"
	DefaultCronExpression = "*/5 * * * *"
)

func SetArgs(zoneName *string, apiToken *string, recordType *string) {
	if len(os.Args) > 3 {
		*zoneName = os.Args[1]
		*apiToken = os.Args[2]
		*recordType = os.Args[3]
	} else {
		*zoneName = os.Getenv(EnvZoneName)
		*apiToken = os.Getenv(EnvApiToken)
		*recordType = os.Getenv(EnvRecordType)
	}
}

func SetOptionalArgs(recordName *string, cronExpression *string) {
	if len(os.Args) > 4 {
		*recordName = os.Args[4]
	} else if os.Getenv(EnvRecordName) != "" {
		*recordName = os.Getenv(EnvRecordName)
	} else {
		*recordName = DefaultRecordName
	}

	if len(os.Args) > 5 {
		*cronExpression = os.Args[5]
	} else if os.Getenv(EnvCronExpression) != "" {
		*cronExpression = os.Getenv(EnvCronExpression)
	} else {
		*cronExpression = DefaultCronExpression
	}
}

func ValidateArgs(zoneName string, apiToken string, recordType string) {
	if zoneName == "" || apiToken == "" || recordType == "" {
		fmt.Println("You need to set the environment variables",
			EnvZoneName, ",", EnvApiToken, "and", EnvRecordType,
			"or provide them as args in the shown order")
		os.Exit(1)
	}

	// Validating given record type
	if recordType != IPv4RecordType && recordType != IPv6RecordType {
		fmt.Println("Given record type does not match",
			IPv4RecordType, "or", IPv6RecordType)
		os.Exit(1)
	}
}
