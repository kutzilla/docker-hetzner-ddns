package args

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

	SetArgs(&zoneName, &apiToken, &recordType)

	assert.Equal("example.com", zoneName)
	assert.Equal("abcdefghi1234567890", apiToken)
	assert.Equal("A", recordType)
}
