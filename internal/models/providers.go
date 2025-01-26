package models

type ExportCredentialsStack []ExporterCredentials

type ExporterCredentials struct {
	API       string            `json:"api,omitempty"`
	Username  string            `json:"username,omitempty"`
	SecretKey string            `json:"secret,omitempty"`
	Extra     map[string]string `json:"extra,omitempty"`
}
