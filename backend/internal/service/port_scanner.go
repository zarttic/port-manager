package service

import (
	"context"
	"log"
	"sync"
	"time"

	"port-manager/internal/model"
	"port-manager/internal/repository"
	"port-manager/pkg/sysinfo"
)

// PortScanner provides port scanning functionality
type PortScanner struct {
	netstat *sysinfo.NetStat
	cache   map[int]model.PortInfo
	mu      sync.RWMutex
}

// NewPortScanner creates a new PortScanner
func NewPortScanner() *PortScanner {
	return &PortScanner{
		netstat: sysinfo.NewNetStat(),
		cache:   make(map[int]model.PortInfo),
	}
}

// ScanPorts scans all ports
func (s *PortScanner) ScanPorts() ([]model.PortInfo, error) {
	ports, err := s.netstat.ScanPorts()
	if err != nil {
		return nil, err
	}

	// Update cache
	s.mu.Lock()
	s.cache = make(map[int]model.PortInfo)
	for _, port := range ports {
		s.cache[port.Port] = port
	}
	s.mu.Unlock()

	return ports, nil
}

// ScanPort scans a specific port
func (s *PortScanner) ScanPort(port int) (*model.PortInfo, error) {
	// Check cache first
	s.mu.RLock()
	if cached, ok := s.cache[port]; ok {
		s.mu.RUnlock()
		return &cached, nil
	}
	s.mu.RUnlock()

	// Scan from system
	return s.netstat.ScanPort(port)
}

// GetCachedPorts gets cached port information
func (s *PortScanner) GetCachedPorts() []model.PortInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ports := make([]model.PortInfo, 0, len(s.cache))
	for _, port := range s.cache {
		ports = append(ports, port)
	}
	return ports
}
