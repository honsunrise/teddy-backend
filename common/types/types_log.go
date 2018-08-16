package types

import "github.com/zhsyourai/URCF-engine/models"

type LogsWithCount struct {
	TotalCount int64        `json:"total_count"`
	Items      []models.Log `json:"items"`
}
