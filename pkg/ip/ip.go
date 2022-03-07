package ip

type Provider interface {
	Request() (IP, error)
	IsOnline() bool
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
