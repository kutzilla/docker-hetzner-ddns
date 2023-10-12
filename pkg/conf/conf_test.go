package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Find a way to reset flags for multiple tests
//func TestReadOptionalEnvs(t *testing.T) {
//	assert := assert.New(t)
//
//	os.Setenv(EnvZoneName, "example.com")
//	os.Setenv(EnvApiToken, "abcdefghi1234567890")
//	os.Setenv(EnvRecordType, "A")
//	os.Setenv(EnvRecordName, "www")
//	os.Setenv(EnvTimeToLive, "43200")
//	os.Setenv(EnvCronExpression, "*/10 * * * *")
//
//	dynDnsConf := Read()
//	// Check optional args
//	assert.NotEqual(DefaultRecordName, dynDnsConf.RecordConf[DefaultRecordKey].RecordName)
//	assert.Equal("www", dynDnsConf.RecordConf[DefaultRecordKey].RecordName)
//
//	assert.NotEqual(DefaultTimeToLive, dynDnsConf.RecordConf[DefaultRecordKey].TTL)
//	assert.Equal(43200, dynDnsConf.RecordConf[DefaultRecordKey].TTL)
//
//	assert.NotEqual(DefaultCronExpression, dynDnsConf.CronConf.CronExpression)
//	assert.Equal("*/10 * * * *", dynDnsConf.CronConf.CronExpression)
//
//	t.Cleanup(func() {
//		flags
//	})
//}

func TestReadMultipleDomains(t *testing.T) {
	assert := assert.New(t)

	os.Setenv(EnvZoneName, "example.com")
	os.Setenv(EnvApiToken, "abcdefghi1234567890")

	os.Setenv("RECORD_NAME_TEST1", "test1")
	os.Setenv("RECORD_NAME_TEST1_TYPE", "A")
	os.Setenv("RECORD_NAME_TEST1_TTL", "4711")

	os.Setenv("RECORD_NAME_TEST2", "test2")
	os.Setenv("RECORD_NAME_TEST2_TYPE", "A")
	os.Setenv("RECORD_NAME_TEST2_TTL", "1337")

	os.Setenv(EnvCronExpression, "*/10 * * * *")

	dynDnsConf := Read()
	// Check optional args
	assert.NotEqual(DefaultRecordName, dynDnsConf.RecordConf["RECORD_NAME_TEST1"].RecordName)
	assert.Equal("test1", dynDnsConf.RecordConf["RECORD_NAME_TEST1"].RecordName)

	assert.NotEqual(DefaultTimeToLive, dynDnsConf.RecordConf["RECORD_NAME_TEST1"].TTL)
	assert.Equal(4711, dynDnsConf.RecordConf["RECORD_NAME_TEST1"].TTL)

	assert.NotEqual(DefaultRecordName, dynDnsConf.RecordConf["RECORD_NAME_TEST2"].RecordName)
	assert.Equal("test2", dynDnsConf.RecordConf["RECORD_NAME_TEST2"].RecordName)

	assert.NotEqual(DefaultTimeToLive, dynDnsConf.RecordConf["RECORD_NAME_TEST2"].TTL)
	assert.Equal(1337, dynDnsConf.RecordConf["RECORD_NAME_TEST2"].TTL)

	assert.NotEqual(DefaultCronExpression, dynDnsConf.CronConf.CronExpression)
	assert.Equal("*/10 * * * *", dynDnsConf.CronConf.CronExpression)
}
