package ip

import (
	"encoding/json"
	"net/http"
	"net/url"

	"matthias-kutz.com/hetzner-ddns/pkg/conf"
	"matthias-kutz.com/hetzner-ddns/pkg/request"
)

type ipify struct {
	IP string `json:"ip"`
}

type Ipify struct {
	IpVersion string
}

const (
	HttpsScheme = "https"

	IpifyIpV4Host             = "api.ipify.org"
	IpifyIpV6Host             = "api64.ipify.org"
	IpifyFormatQueryParam     = "format"
	IpifyQueryParamFormatJson = "json"
)

func (i Ipify) getHost() string {
	if i.IpVersion == conf.IPv6 {
		return IpifyIpV6Host
	}
	return IpifyIpV4Host
}

func (i Ipify) IsOnline() bool {
	ipifyUrl := url.URL{
		Scheme: HttpsScheme,
		Host:   i.getHost(),
	}

	_, err := http.Get(ipifyUrl.String())
	return err == nil
}

func (i Ipify) Request() (IP, error) {
	requestUrl := url.URL{
		Scheme:   HttpsScheme,
		Host:     i.getHost(),
		RawQuery: IpifyFormatQueryParam + "=" + IpifyQueryParamFormatJson,
	}

	respBody, err := request.Request(http.MethodGet, requestUrl,
		map[string]string{}, []byte{})

	if err != nil {
		return IP{}, &ProviderNotAvailableError{ProviderName: i.getHost()}
	}

	var ipify ipify
	json.Unmarshal(respBody, &ipify)

	ip := IP{
		Value:  ipify.IP,
		Source: i.getHost(),
	}

	return ip, nil
}
