package ip

type IpVersion string

const (
	IpV4 IpVersion = "IpV4"
	IpV6 IpVersion = "IpV6"
)

type Provider interface {
	Request(IpVersion) (IP, error)
	IsOnline(IpVersion) bool
}

type IP struct {
	Value  string
	Source string
}

type ProviderNotAvailableError struct {
	ProviderName string
}

func (e *ProviderNotAvailableError) Error() string {
	return "The provider " + e.ProviderName + " is not available"
}
