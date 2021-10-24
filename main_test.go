package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetArgsWithEnvs(t *testing.T) {
	assert := assert.New(t)

	var zoneName, apiToken, recordType string

	os.Setenv(EnvZoneName, "example.com")
	os.Setenv(EnvApiToken, "abcdefghi1234567890")
	os.Setenv(EnvRecordType, "A")

	// Clear Args for the test
	os.Args = []string{}

	setArgs(&zoneName, &apiToken, &recordType)

	assert.Equal("example.com", zoneName)
	assert.Equal("abcdefghi1234567890", apiToken)
	assert.Equal("A", recordType)
}

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

	testZone, err := findZoneByName(zones, "test.com")
	assert.Equal("test.com", testZone.Name)
	assert.Nil(err)

	exampleZone, err := findZoneByName(zones, "example.com")
	assert.Equal("example.com", exampleZone.Name)
	assert.Nil(err)

	notAvailableZone, err := findZoneByName(zones, "not-avaible.com")
	assert.NotEqual("not-avaible.com", notAvailableZone.Name)
	assert.NotNil(err)

}
