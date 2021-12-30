package ip

import (
	"encoding/json"
	"net/http"
	"net/url"

	"matthias-kutz.com/hetzner-ddns/pkg/request"
)

type ipify struct {
	IP string `json:"ip"`
}

type Ipify struct {
}

const (
	HttpsScheme = "https"

	IpifyHost                 = "api.ipify.org"
	IpifyFormatQueryParam     = "format"
	IpifyQueryParamFormatJson = "json"
)

func (i Ipify) IsOnline() bool {
	ipifyUrl := url.URL{
		Scheme: HttpsScheme,
		Host:   IpifyHost,
	}

	_, err := http.Get(ipifyUrl.String())
	return err == nil
}

func (i Ipify) Request() (IP, error) {
	requestUrl := url.URL{
		Scheme:   HttpsScheme,
		Host:     IpifyHost,
		RawQuery: IpifyFormatQueryParam + "=" + IpifyQueryParamFormatJson,
	}

	respBody, err := request.Request(http.MethodGet, requestUrl,
		map[string]string{}, []byte{})

	if err != nil {
		return IP{}, &ProviderNotAvailableError{ProviderName: IpifyHost}
	}

	var ipify ipify
	json.Unmarshal(respBody, &ipify)

	ip := IP{
		Value:  ipify.IP,
		Source: IpifyHost,
	}

	return ip, nil
}
