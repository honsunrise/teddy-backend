package types

import "github.com/zhsyourai/URCF-engine/models"

type PluginsWithCount struct {
	TotalCount int64           `json:"total_count"`
	Items      []models.Plugin `json:"items"`
}

type PluginCommandExecResult struct {
	Result string `json:"result"`
}
