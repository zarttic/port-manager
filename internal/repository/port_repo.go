package repository

import (
	"database/sql"
	"fmt"
	"time"

	"port-manager/internal/model"
)

// PortRepository handles port data persistence
type PortRepository struct {
	db *sql.DB
}

// NewPortRepository creates a new PortRepository
func NewPortRepository(db *sql.DB) *PortRepository {
	return &PortRepository{db: db}
}

// Insert inserts a new port usage record
func (r *PortRepository) Insert(usage *model.PortUsage) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO port_usage (port, protocol, pid, process_name, process_path, start_time)
		VALUES (?, ?, ?, ?, ?, ?)
	`, usage.Port, usage.Protocol, usage.PID, usage.ProcessName, usage.ProcessPath, usage.StartTime)
	if err != nil {
		return 0, fmt.Errorf("failed to insert port usage: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateEndTime updates the end time of a port usage record
func (r *PortRepository) UpdateEndTime(id int64, endTime time.Time, duration int64) error {
	_, err := r.db.Exec(`
		UPDATE port_usage
		SET end_time = ?, duration = ?
		WHERE id = ?
	`, endTime, duration, id)
	if err != nil {
		return fmt.Errorf("failed to update end time: %w", err)
	}

	return nil
}

// GetActiveByPort gets active port usage by port number
func (r *PortRepository) GetActiveByPort(port int) (*model.PortUsage, error) {
	var usage model.PortUsage
	err := r.db.QueryRow(`
		SELECT id, port, protocol, pid, process_name, process_path, start_time
		FROM port_usage
		WHERE port = ? AND end_time IS NULL
		ORDER BY start_time DESC
		LIMIT 1
	`, port).Scan(
		&usage.ID, &usage.Port, &usage.Protocol, &usage.PID,
		&usage.ProcessName, &usage.ProcessPath, &usage.StartTime,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active port: %w", err)
	}

	return &usage, nil
}

// GetActiveByPID gets active port usage by PID
func (r *PortRepository) GetActiveByPID(pid int) ([]model.PortUsage, error) {
	rows, err := r.db.Query(`
		SELECT id, port, protocol, pid, process_name, process_path, start_time
		FROM port_usage
		WHERE pid = ? AND end_time IS NULL
	`, pid)
	if err != nil {
		return nil, fmt.Errorf("failed to get active ports by PID: %w", err)
	}
	defer rows.Close()

	var usages []model.PortUsage
	for rows.Next() {
		var usage model.PortUsage
		err := rows.Scan(
			&usage.ID, &usage.Port, &usage.Protocol, &usage.PID,
			&usage.ProcessName, &usage.ProcessPath, &usage.StartTime,
		)
		if err != nil {
			return nil, err
		}
		usages = append(usages, usage)
	}

	return usages, nil
}

// GetByPort gets port usage history
func (r *PortRepository) GetByPort(port int, limit int) ([]model.PortUsage, error) {
	rows, err := r.db.Query(`
		SELECT id, port, protocol, pid, process_name, process_path, start_time, end_time, duration
		FROM port_usage
		WHERE port = ?
		ORDER BY start_time DESC
		LIMIT ?
	`, port, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get port history: %w", err)
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

// CleanupOldData removes old records
func (r *PortRepository) CleanupOldData(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	_, err := r.db.Exec(`
		DELETE FROM port_usage
		WHERE end_time < ? AND end_time IS NOT NULL
	`, cutoff)
	if err != nil {
		return fmt.Errorf("failed to cleanup old data: %w", err)
	}

	return nil
}
