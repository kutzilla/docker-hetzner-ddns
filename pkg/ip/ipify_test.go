package ip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsOnline(t *testing.T) {
	assert := assert.New(t)

	ipify := Ipify{}
	assert.True(ipify.IsOnline(IpV4))
}

func TestRequest(t *testing.T) {
	assert := assert.New(t)

	ipify := Ipify{}
	ip, _ := ipify.Request(IpV4)
	assert.NotNil(ip.Value)
	assert.NotNil(ip.Source)
	assert.Equal("api.ipify.org", ip.Source)
}

func TestRequestWithIPv6(t *testing.T) {
	assert := assert.New(t)

	ipify := Ipify{}
	ip, _ := ipify.Request(IpV6)
	assert.NotNil(ip.Value)
	assert.NotNil(ip.Source)
	assert.Equal("api64.ipify.org", ip.Source)
}
