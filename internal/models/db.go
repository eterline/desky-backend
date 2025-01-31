package models

import "gorm.io/gorm"

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

// Widget service repository tables ===========================

type WidgetsT struct {
	gorm.Model

	ID          uint `gorm:"primaryKey"`
	Name        string
	Icon        string `gorm:"uniqueIndex"`
	Description string
	Topic       string
}

// Docker service repository tables ===========================
