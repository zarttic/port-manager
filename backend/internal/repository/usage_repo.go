package repository

import (
	"database/sql"
	"fmt"
	"time"

	"port-manager/internal/model"
)

// UsageRepository handles usage statistics
type UsageRepository struct {
	db *sql.DB
}

// NewUsageRepository creates a new UsageRepository
func NewUsageRepository(db *sql.DB) *UsageRepository {
	return &UsageRepository{db: db}
}

// GetPortStats gets statistics for a specific port
func (r *UsageRepository) GetPortStats(port int) (*model.PortStats, error) {
	stats := &model.PortStats{Port: port}

	// Get total usage time and count
	err := r.db.QueryRow(`
		SELECT
			COALESCE(SUM(duration), 0) as total_usage,
			COUNT(*) as usage_count,
			MAX(start_time) as last_used
		FROM port_usage
		WHERE port = ? AND end_time IS NOT NULL
	`, port).Scan(&stats.TotalUsage, &stats.UsageCount, &stats.LastUsed)
	if err != nil {
		return nil, fmt.Errorf("failed to get port stats: %w", err)
	}

	// Get top processes for this port
	rows, err := r.db.Query(`
		SELECT
			process_name,
			COUNT(*) as usage_count,
			SUM(duration) as total_time
		FROM port_usage
		WHERE port = ? AND end_time IS NOT NULL
		GROUP BY process_name
		ORDER BY total_time DESC
		LIMIT 5
	`, port)
	if err != nil {
		return nil, fmt.Errorf("failed to get top processes: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pu model.ProcessUsage
		err := rows.Scan(&pu.ProcessName, &pu.UsageCount, &pu.TotalTime)
		if err != nil {
			return nil, err
		}
		stats.TopProcesses = append(stats.TopProcesses, pu)
	}

	return stats, nil
}

// GetTopUsedPorts gets most frequently used ports
func (r *UsageRepository) GetTopUsedPorts(limit int) ([]model.PortStats, error) {
	rows, err := r.db.Query(`
		SELECT
			port,
			SUM(duration) as total_usage,
			COUNT(*) as usage_count,
			MAX(start_time) as last_used
		FROM port_usage
		WHERE end_time IS NOT NULL
		GROUP BY port
		ORDER BY total_usage DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top used ports: %w", err)
	}
	defer rows.Close()

	var statsList []model.PortStats
	for rows.Next() {
		var stats model.PortStats
		err := rows.Scan(&stats.Port, &stats.TotalUsage, &stats.UsageCount, &stats.LastUsed)
		if err != nil {
			return nil, err
		}
		statsList = append(statsList, stats)
	}

	return statsList, nil
}

// GetUsageHistory gets usage history within time range
func (r *UsageRepository) GetUsageHistory(start, end time.Time) ([]model.PortUsage, error) {
	rows, err := r.db.Query(`
		SELECT id, port, protocol, pid, process_name, process_path, start_time, end_time, duration
		FROM port_usage
		WHERE start_time >= ? AND start_time <= ?
		ORDER BY start_time DESC
	`, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage history: %w", err)
	}
	defer rows.Close()

	var usages []model.PortUsage
	for rows.Next() {
		var usage model.PortUsage
		err := rows.Scan(
			&usage.ID, &usage.Port, &usage.Protocol, &usage.PID,
			&usage.ProcessName, &usage.ProcessPath, &usage.StartTime,
			&usage.EndTime, &usage.Duration,
		)
		if err != nil {
			return nil, err
		}
		usages = append(usages, usage)
	}

	return usages, nil
}
