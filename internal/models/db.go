package models

import (
	"encoding/json"
	"fmt"
)

// Apps service repository tables ===========================

type AppsTopicT struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex,unique"`
}

type AppsInstancesT struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Icon        string
	Description string
	Link        string

	TopicID uint
	Topic   AppsTopicT `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// User service repository tables ===========================

type DeskyUserT struct {
	ID       uint   `gorm:"primaryKey"`
	Login    string `gorm:"uniqueIndex"`
	Password string
}

func NewDeskyUserT(login, password string) *DeskyUserT {
	return &DeskyUserT{
		Login:    login,
		Password: password,
	}
}

// Exports service repository tables ===========================

type ExporterInfoT struct {
	ID    uint `gorm:"primarykey"`
	Type  string
	API   string `gorm:"uniqueIndex"`
	Extra string
}

func (t *ExporterInfoT) ResolveType() ExporterTypeString {
	return ExporterTypeString(t.Type)
}

func (t *ExporterInfoT) ResolveExtra() map[ExporterExtraField]any {
	extra := make(map[ExporterExtraField]any)
	json.Unmarshal([]byte(t.Extra), &extra)
	return extra
}

// SSHLander service repository tables ===========================

type SSHCredentialsT struct {
	ID uint `gorm:"primaryKey"`

	OperationSystemID uint
	OperationSystem   SSHSystemTypesT `gorm:"foreignKey:OperationSystemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Username string
	Host     string
	Port     uint16

	SecurityID uint
	Security   SSHSecureT `gorm:"foreignKey:SecurityID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (c *SSHCredentialsT) ValueUser() string {
	return c.Username
}

func (c *SSHCredentialsT) UsePrivateKey() bool {
	return c.Security.PrivateKeyUse
}

func (c *SSHCredentialsT) Password() string {
	return c.Security.Password
}

func (c *SSHCredentialsT) PrivateKey() []byte {
	return []byte(c.Security.PrivateKey)
}

func (c *SSHCredentialsT) Socket() string {
	return fmt.Sprintf("%s:%v", c.Host, c.Port)
}

type SSHSystemTypesT struct {
	ID         uint      `gorm:"primaryKey"`
	SystemType SSHtypeOS `gorm:"uniqueIndex"`
}

func MakeSSHSystemTypesT(value string) SSHSystemTypesT {
	return SSHSystemTypesT{
		SystemType: StringSSHtypeOS(value),
	}
}

type SSHSecureT struct {
	ID            uint `gorm:"primaryKey"`
	PrivateKeyUse bool
	Password      string
	PrivateKey    string
}

func MakeSSHSecureT(
	password string,
	keyUse bool,
	key string,
) SSHSecureT {
	return SSHSecureT{
		PrivateKeyUse: keyUse,
		Password:      password,
		PrivateKey:    key,
	}
}
