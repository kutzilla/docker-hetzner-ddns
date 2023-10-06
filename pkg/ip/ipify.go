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
}

const (
	HttpsScheme = "https"

	IpifyIpV4Host             = "api.ipify.org"
	IpifyIpV6Host             = "api64.ipify.org"
	IpifyFormatQueryParam     = "format"
	IpifyQueryParamFormatJson = "json"
)

func (i Ipify) determineHostByIpVersion(ipVersion IpVersion) string {
	if ipVersion == conf.IPv6 {
		return IpifyIpV6Host
	}
	return IpifyIpV4Host
}

func (i Ipify) IsOnline(ipVersion IpVersion) bool {
	ipifyUrl := url.URL{
		Scheme: HttpsScheme,
		Host:   i.determineHostByIpVersion(ipVersion),
	}

	_, err := http.Get(ipifyUrl.String())
	return err == nil
}

func (i Ipify) Request(ipVersion IpVersion) (IP, error) {
	requestUrl := url.URL{
		Scheme:   HttpsScheme,
		Host:     i.determineHostByIpVersion(ipVersion),
		RawQuery: IpifyFormatQueryParam + "=" + IpifyQueryParamFormatJson,
	}

	respBody, err := request.Request(http.MethodGet, requestUrl,
		map[string]string{}, []byte{})

	if err != nil {
		return IP{}, &ProviderNotAvailableError{ProviderName: i.determineHostByIpVersion(ipVersion)}
	}

	var ipify ipify
	json.Unmarshal(respBody, &ipify)

	ip := IP{
		Value:  ipify.IP,
		Source: i.determineHostByIpVersion(ipVersion),
	}

	return ip, nil
}
