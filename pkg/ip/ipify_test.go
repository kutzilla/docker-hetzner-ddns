package ip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsOnline(t *testing.T) {
	assert := assert.New(t)

	ipify := Ipify{}
	assert.True(ipify.IsOnline())
}

func TestRequest(t *testing.T) {
	assert := assert.New(t)

	ipify := Ipify{}
	ip, _ := ipify.Request()
	assert.NotNil(ip.Value)
	assert.NotNil(ip.Source)
	assert.Equal("api.ipify.org", ip.Source)
}
