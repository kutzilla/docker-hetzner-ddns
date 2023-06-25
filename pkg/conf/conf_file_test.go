package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfFile(t *testing.T) {
	assert := assert.New(t)

	dynDnsConf := ReadFile("conf_file_test.json")

	assert.Equal("example.com", dynDnsConf.DnsConf.ZoneName)
	assert.Equal("my-secret-api-token", dynDnsConf.DnsConf.ApiToken)

}
