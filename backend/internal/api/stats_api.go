package api

import (
	"time"

	"port-manager/internal/repository"
)

// StatsAPI provides statistics-related API endpoints
type StatsAPI struct {
	usageRepo *repository.UsageRepository
}

// NewStatsAPI creates a new StatsAPI
func NewStatsAPI(usageRepo *repository.UsageRepository) *StatsAPI {
	return &StatsAPI{
		usageRepo: usageRepo,
	}
}

// GetPortStats gets statistics for a specific port
func (a *StatsAPI) GetPortStats(port int) (map[string]interface{}, error) {
	stats, err := a.usageRepo.GetPortStats(port)
	if err != nil {
		return nil, err
	}

	topProcesses := make([]map[string]interface{}, len(stats.TopProcesses))
	for i, proc := range stats.TopProcesses {
		topProcesses[i] = map[string]interface{}{
			"processName": proc.ProcessName,
			"usageCount":  proc.UsageCount,
			"totalTime":   proc.TotalTime,
		}
	}

	return map[string]interface{}{
		"port":         stats.Port,
		"totalUsage":   stats.TotalUsage,
		"usageCount":   stats.UsageCount,
		"lastUsed":     stats.LastUsed,
		"topProcesses": topProcesses,
	}, nil
}

// GetTopUsedPorts gets most frequently used ports
func (a *StatsAPI) GetTopUsedPorts(limit int) ([]map[string]interface{}, error) {
	stats, err := a.usageRepo.GetTopUsedPorts(limit)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(stats))
	for i, s := range stats {
		result[i] = map[string]interface{}{
			"port":       s.Port,
			"totalUsage": s.TotalUsage,
			"usageCount": s.UsageCount,
			"lastUsed":   s.LastUsed,
		}
	}

	return result, nil
}

// GetUsageHistory gets usage history within time range
func (a *StatsAPI) GetUsageHistory(start, end string) ([]map[string]interface{}, error) {
	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return nil, err
	}

	usages, err := a.usageRepo.GetUsageHistory(startTime, endTime)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(usages))
	for i, usage := range usages {
		result[i] = map[string]interface{}{
			"id":          usage.ID,
			"port":        usage.Port,
			"protocol":    usage.Protocol,
			"pid":         usage.PID,
			"processName": usage.ProcessName,
			"startTime":   usage.StartTime,
			"endTime":     usage.EndTime,
			"duration":    usage.Duration,
		}
	}

	return result, nil
}
