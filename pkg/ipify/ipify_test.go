package ipify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsOnline(t *testing.T) {
	assert := assert.New(t)
	assert.True(IsOnline())
}
