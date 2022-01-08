package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/namsral/flag"
)

func TestReadEnvs(t *testing.T) {
	flag.ResetForTesting(nil)
	assert := assert.New(t)

	os.Setenv(EnvZoneName, "example.com")
	os.Setenv(EnvApiToken, "abcdefghi1234567890")
	os.Setenv(EnvRecordType, "A")

	dynDnsConf, err := Read()
	// Check mandatory args
	assert.Equal("example.com", dynDnsConf.DnsConf.ZoneName)
	assert.Equal("abcdefghi1234567890", dynDnsConf.DnsConf.ApiToken)
	assert.Equal("A", dynDnsConf.RecordConf.RecordType)
	// Check optional args
	assert.Equal(DefaultRecordName, dynDnsConf.RecordConf.RecordName)
	assert.Equal(DefaultTimeToLive, dynDnsConf.RecordConf.TTL)
	assert.Equal(DefaultCronExpression, dynDnsConf.CronConf.CronExpression)
	// Check error is nil
	assert.Nil(err)
}

func TestReadOptionalEnvs(t *testing.T) {
	flag.ResetForTesting(nil)
	assert := assert.New(t)

	os.Setenv(EnvZoneName, "example.com")
	os.Setenv(EnvApiToken, "abcdefghi1234567890")
	os.Setenv(EnvRecordType, "A")
	os.Setenv(EnvRecordName, "www")
	os.Setenv(EnvTimeToLive, "43200")
	os.Setenv(EnvCronExpression, "*/10 * * * *")

	dynDnsConf, _ := Read()
	// Check optional args
	assert.NotEqual(DefaultRecordName, dynDnsConf.RecordConf.RecordName)
	assert.Equal("www", dynDnsConf.RecordConf.RecordName)

	assert.NotEqual(DefaultTimeToLive, dynDnsConf.RecordConf.TTL)
	assert.Equal(43200, dynDnsConf.RecordConf.TTL)

	assert.NotEqual(DefaultCronExpression, dynDnsConf.CronConf.CronExpression)
	assert.Equal("*/10 * * * *", dynDnsConf.CronConf.CronExpression)
}

func TestAllMissingArgs(t *testing.T) {
	flag.ResetForTesting(nil)
	assert := assert.New(t)

	// No args
	_, err := Read()
	assert.NotNil(err)
}

func TestFewMissingArgs(t *testing.T) {
	flag.ResetForTesting(nil)
	assert := assert.New(t)

	// Missing mandatory args
	os.Setenv(EnvZoneName, "example.com")
	_, err := Read()
	assert.NotNil(err)
}
