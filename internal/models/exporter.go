package models

import "encoding/json"

type ExporterTypeString string
type ExporterExtraField string

type ExporterForm interface {
	ValueType() ExporterTypeString
	ValueAPI() string
	ValueExtra() string
}

const (
	ExporterProxmoxType ExporterTypeString = "proxmox"
	ExporterDockerType  ExporterTypeString = "docker"
)

const (
	SecretField   ExporterExtraField = "secret"
	PasswordField ExporterExtraField = "password"
	TokenField    ExporterExtraField = "token"
	LoginField    ExporterExtraField = "login"

	NodeNameField  ExporterExtraField = "node"
	DockerEnvField ExporterExtraField = "environment"
)

type ExporterInfo struct {
	ID    uint                       `json:"id,omitempty"`
	Type  ExporterTypeString         `json:"type"`
	API   string                     `json:"api"`
	Extra map[ExporterExtraField]any `json:"extra,omitempty"`
}

// ============= Models for http body requests append controllers =============

type ProxmoxFormExport struct {
	API      string `json:"api" validate:"url"`
	NodeName string `json:"node-name" validate:"required"`
	Login    string `json:"login" validate:"required"`
	Password string `json:"token" validate:"required"`
}

func (form *ProxmoxFormExport) ValueType() ExporterTypeString {
	return ExporterProxmoxType
}

func (form *ProxmoxFormExport) ValueAPI() string {
	return form.API
}

func (form *ProxmoxFormExport) ValueExtra() string {
	extra := &map[ExporterExtraField]any{
		NodeNameField: form.NodeName,
		LoginField:    form.Login,
		PasswordField: form.Password,
	}
	return extraFieldEncoder(extra)
}

// ==================

type DockerFormExport struct {
	API     string `json:"api" validate:"url"`
	EnvName string `json:"environment" validate:"required"`
}

func (form *DockerFormExport) ValueType() ExporterTypeString {
	return ExporterDockerType
}

func (form *DockerFormExport) ValueAPI() string {
	return form.API
}

func (form *DockerFormExport) ValueExtra() string {
	extra := &map[ExporterExtraField]any{
		DockerEnvField: form.EnvName,
	}
	return extraFieldEncoder(extra)
}

// ==================

func extraFieldEncoder(v *map[ExporterExtraField]any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	return string(data)
}
