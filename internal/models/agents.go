package models

var ExporterList = []string{"host", "cpu", "ram", "load", "temperature", "ports", "partitions"}

type SessionCredentials struct {
	Hostname string `json:"hostname"`
	ID       string `json:"id"`
	URL      string `json:"url"`
}

type FetchedResponse struct {
	ID   string         `json:"id"`
	Data map[string]any `json:"data"`
}

type FetchedResponseSingle struct {
	ID   string `json:"id"`
	Data any    `json:"data"`
}
