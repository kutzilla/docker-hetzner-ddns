package ipify

import (
	"encoding/json"
	"net/http"
	"net/url"

	"matthias-kutz.com/hetzner-ddns/pkg/request"
)

type Ipify struct {
	IP string `json:"ip"`
}

const (
	HttpsScheme = "https"

	IpifyHost                 = "api.ipify.org"
	IpifyFormatQueryParam     = "format"
	IpifyQueryParamFormatJson = "json"
)

func IsOnline() bool {
	ipifyUrl := url.URL{
		Scheme: HttpsScheme,
		Host:   IpifyHost,
	}

	_, err := http.Get(ipifyUrl.String())
	return err == nil
}

func RequestIpify() Ipify {
	requestUrl := url.URL{
		Scheme:   HttpsScheme,
		Host:     IpifyHost,
		RawQuery: IpifyFormatQueryParam + "=" + IpifyQueryParamFormatJson,
	}

	respBody := request.Request(http.MethodGet, requestUrl,
		map[string]string{}, []byte{})

	var ipify Ipify
	json.Unmarshal(respBody, &ipify)

	return ipify
}
