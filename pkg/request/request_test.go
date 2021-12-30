package request

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	assert := assert.New(t)

	resp, _ := Request(http.MethodGet, url.URL{Scheme: "https", Host: "matthias-kutz.com"},
		map[string]string{}, []byte{})

	assert.NotNil(resp)

}
