package providers

import "fmt"

type ProvidersErrors struct {
	err error
}

func (e *ProvidersErrors) Error() string {
	return e.err.Error()
}

var (
	ErrUnknownService = &ProvidersErrors{
		err: fmt.Errorf("unknown service"),
	}
	ErrNotConfiguredExport = &ProvidersErrors{
		err: fmt.Errorf("service exporter not configured"),
	}
)
