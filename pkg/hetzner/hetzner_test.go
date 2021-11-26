package hetzner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindZoneByName(t *testing.T) {
	assert := assert.New(t)

	zones := Zones{
		Zone: []Zone{
			{
				Name: "example.com",
			},
			{
				Name: "test.com",
			},
		},
	}

	testZone, err := FindZoneByName(zones, "test.com")
	assert.Equal("test.com", testZone.Name)
	assert.Nil(err)

	exampleZone, err := FindZoneByName(zones, "example.com")
	assert.Equal("example.com", exampleZone.Name)
	assert.Nil(err)

	notAvailableZone, err := FindZoneByName(zones, "not-avaible.com")
	assert.NotEqual("not-avaible.com", notAvailableZone.Name)
	assert.NotNil(err)

}
