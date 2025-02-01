package models

var ExporterList = []string{"host", "cpu", "ram", "load", "temperature", "ports", "partitions"}

type SessionCredentials struct {
	Hostname string `json:"hostname"`
	ID       string `json:"id"`
	Valid    bool   `json:"valid"`
	URL      string `json:"url"`
}

type FetchedResponse struct {
	SessionCredentials
	Data map[string]any `json:"data"`
}
