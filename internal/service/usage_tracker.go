package service

import (
	"context"
	"log"
	"sync"
	"time"

	"port-manager/internal/model"
	"port-manager/internal/repository"
)

// UsageTracker tracks port usage over time
type UsageTracker struct {
	portRepo  *repository.PortRepository
	usageRepo *repository.UsageRepository
	tracking  map[int]int64 // port -> usage ID
	mu        sync.RWMutex
	stopCh    chan struct{}
}

// NewUsageTracker creates a new UsageTracker
func NewUsageTracker(portRepo *repository.PortRepository, usageRepo *repository.UsageRepository) *UsageTracker {
	return &UsageTracker{
		portRepo:  portRepo,
		usageRepo: usageRepo,
		tracking:  make(map[int]int64),
		stopCh:    make(chan struct{}),
	}
}

// StartTracking starts background tracking
func (t *UsageTracker) StartTracking(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Initial cleanup
	if err := t.portRepo.CleanupOldData(30); err != nil {
		log.Printf("Failed to cleanup old data: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.stopCh:
			return
		case <-ticker.C:
			// Periodic tasks (if needed)
		}
	}
}

// StopTracking stops background tracking
func (t *UsageTracker) StopTracking() {
	close(t.stopCh)
}

// RecordPortUsage records port usage
func (t *UsageTracker) RecordPortUsage(ports []model.PortInfo) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Track new ports
	for _, port := range ports {
		if _, tracked := t.tracking[port.Port]; !tracked {
			// New port usage
			usage := &model.PortUsage{
				Port:        port.Port,
				Protocol:    port.Protocol,
				PID:         port.PID,
				ProcessName: port.ProcessName,
				StartTime:   time.Now(),
			}

			id, err := t.portRepo.Insert(usage)
			if err != nil {
				log.Printf("Failed to record port usage: %v", err)
				continue
			}

			t.tracking[port.Port] = id
		}
	}

	// Check for closed ports
	for port, id := range t.tracking {
		found := false
		for _, p := range ports {
			if p.Port == port {
				found = true
				break
			}
		}

		if !found {
			// Port closed
			endTime := time.Now()

			// Get start time (simplified - in production, fetch from DB)
			duration := int64(0)

			err := t.portRepo.UpdateEndTime(id, endTime, duration)
			if err != nil {
				log.Printf("Failed to update port usage end time: %v", err)
			}

			delete(t.tracking, port)
		}
	}
}

// GetPortStats gets statistics for a port
func (t *UsageTracker) GetPortStats(port int) (*model.PortStats, error) {
	return t.usageRepo.GetPortStats(port)
}

// GetTopUsedPorts gets most used ports
func (t *UsageTracker) GetTopUsedPorts(limit int) ([]model.PortStats, error) {
	return t.usageRepo.GetTopUsedPorts(limit)
}

// GetUsageHistory gets usage history
func (t *UsageTracker) GetUsageHistory(start, end time.Time) ([]model.PortUsage, error) {
	return t.usageRepo.GetUsageHistory(start, end)
}
