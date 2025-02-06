package exporters

import "errors"

type ExporterServiceError struct {
	err error
}

func (es *ExporterServiceError) Error() string {
	return es.err.Error()
}

var ErrExporterNotExists = &ExporterServiceError{
	err: errors.New("exporter not exists"),
}
