package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadOptionalEnvs(t *testing.T) {
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
