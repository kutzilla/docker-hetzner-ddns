package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleDomain(t *testing.T) {
	assert := assert.New(t)

	os.Setenv(EnvZoneName, "example.com")
	os.Setenv(EnvApiToken, "abcdefghi1234567890")
	os.Setenv(EnvRecordType, "A")

	dynDnsConf := Read()

	assert.Equal("@", dynDnsConf.RecordConf.records["default"].RecordName)
	assert.Equal("A", dynDnsConf.RecordConf.records["default"].RecordType)
}

func TestReadMultipleDomains(t *testing.T) {
	assert := assert.New(t)

	os.Setenv(EnvMultipleDomainMode, "true")

	os.Setenv(EnvZoneName, "example.com")
	os.Setenv(EnvApiToken, "abcdefghi1234567890")
	os.Setenv(EnvRecordType, "A")

	os.Setenv("RECORD_NAME_TEST1", "test1")
	os.Setenv("RECORD_NAME_TEST1_TTL", "4711")

	os.Setenv("RECORD_NAME_TEST2", "test2")
	os.Setenv("RECORD_NAME_TEST2_TTL", "1337")

	os.Setenv(EnvCronExpression, "*/10 * * * *")

	dynDnsConf := Read()

	assert.Equal(true, dynDnsConf.RecordConf.mode)

	// Check optional args
	assert.NotEqual(DefaultRecordName, dynDnsConf.RecordConf.records["RECORD_NAME_TEST1"].RecordName)
	assert.Equal("test1", dynDnsConf.RecordConf.records["RECORD_NAME_TEST1"].RecordName)

	assert.NotEqual(DefaultTimeToLive, dynDnsConf.RecordConf.records["RECORD_NAME_TEST1"].TTL)
	assert.Equal(4711, dynDnsConf.RecordConf.records["RECORD_NAME_TEST1"].TTL)

	assert.NotEqual(DefaultRecordName, dynDnsConf.RecordConf.records["RECORD_NAME_TEST2"].RecordName)
	assert.Equal("test2", dynDnsConf.RecordConf.records["RECORD_NAME_TEST2"].RecordName)

	assert.NotEqual(DefaultTimeToLive, dynDnsConf.RecordConf.records["RECORD_NAME_TEST2"].TTL)
	assert.Equal(1337, dynDnsConf.RecordConf.records["RECORD_NAME_TEST2"].TTL)

	assert.NotEqual(DefaultCronExpression, dynDnsConf.CronConf.CronExpression)
	assert.Equal("*/10 * * * *", dynDnsConf.CronConf.CronExpression)
}
