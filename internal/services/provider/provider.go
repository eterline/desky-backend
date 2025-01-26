package provider

import (
	"github.com/eterline/desky-backend/internal/configuration"
	"github.com/eterline/desky-backend/internal/models"
)

// ========================== Proxmox VE provider ==========================

type ProxmoxProvider struct {
	sessions []configuration.PVEInstance
}

func NewProxmoxProvider(c configuration.ServicesParameters) *ProxmoxProvider {
	return &ProxmoxProvider{
		sessions: c.PVE,
	}
}

func (pd *ProxmoxProvider) Get() models.ExportCredentialsStack {

	var stack models.ExportCredentialsStack

	if len(pd.sessions) > 0 {
		for _, s := range pd.sessions {
			stack = append(stack, models.ExporterCredentials{
				API:       s.API,
				Username:  s.Username,
				SecretKey: s.Secret,
				Extra: map[string]string{
					"node": s.Node,
				},
			})
		}
	}

	return stack
}

func (pd *ProxmoxProvider) Delete(query int) error {
	panic("not implemented")
}

func (pd *ProxmoxProvider) Append(models.ExporterCredentials) error {
	panic("not implemented")
}

// ========================== Docker provider ==========================

type DockerProvider struct {
	sessions []configuration.DockerInstance
}

func NewDockerProvider(c configuration.ServicesParameters) *DockerProvider {
	return &DockerProvider{
		sessions: c.Docker,
	}
}

func (pd *DockerProvider) Get() models.ExportCredentialsStack {

	var stack models.ExportCredentialsStack

	if len(pd.sessions) > 0 {
		for _, s := range pd.sessions {
			stack = append(stack, models.ExporterCredentials{
				API: s.API,
				Extra: map[string]string{
					"name": s.Name,
				},
			})
		}
	}

	return stack
}

func (pd *DockerProvider) Delete(query int) error {
	panic("not implemented")
}

func (pd *DockerProvider) Append(models.ExporterCredentials) error {
	panic("not implemented")
}
