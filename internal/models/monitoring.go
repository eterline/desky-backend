package models

import "github.com/eterline/desky-backend/internal/services/system"

type StatsResponse struct {
	RAM  *system.RAMInfo     `json:"memory"`
	CPU  *system.CPUInfo     `json:"cpu"`
	Temp []system.SensorInfo `json:"temperature"`
	Load *system.AverageLoad `json:"load"`
}
