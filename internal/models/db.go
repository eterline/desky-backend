package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

// Apps service repository tables ===========================

type AppsTopicT struct {
	gorm.Model

	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex,unique"`
}

type AppsInstancesT struct {
	gorm.Model

	ID          uint `gorm:"primaryKey"`
	Name        string
	Icon        string
	Description string
	Link        string

	TopicID uint       // External key
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
