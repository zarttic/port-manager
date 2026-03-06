package api

import (
	"port-manager/internal/service"
)

// PortAPI provides port-related API endpoints
type PortAPI struct {
	scanner      *service.PortScanner
	usageTracker *service.UsageTracker
}

// NewPortAPI creates a new PortAPI
func NewPortAPI(scanner *service.PortScanner, usageTracker *service.UsageTracker) *PortAPI {
	return &PortAPI{
		scanner:      scanner,
		usageTracker: usageTracker,
	}
}

// ScanPorts scans all ports
func (a *PortAPI) ScanPorts() ([]map[string]interface{}, error) {
	ports, err := a.scanner.ScanPorts()
	if err != nil {
		return nil, err
	}

	// Record usage
	go a.usageTracker.RecordPortUsage(ports)

	// Convert to map for JSON serialization
	result := make([]map[string]interface{}, len(ports))
	for i, port := range ports {
		result[i] = map[string]interface{}{
			"port":        port.Port,
			"protocol":    port.Protocol,
			"state":       port.State,
			"pid":         port.PID,
			"processName": port.ProcessName,
			"localAddr":   port.LocalAddr,
			"remoteAddr":  port.RemoteAddr,
			"createTime":  port.CreateTime,
		}
	}

	return result, nil
}

// ScanPort scans a specific port
func (a *PortAPI) ScanPort(port int) (map[string]interface{}, error) {
	portInfo, err := a.scanner.ScanPort(port)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"port":        portInfo.Port,
		"protocol":    portInfo.Protocol,
		"state":       portInfo.State,
		"pid":         portInfo.PID,
		"processName": portInfo.ProcessName,
		"localAddr":   portInfo.LocalAddr,
		"remoteAddr":  portInfo.RemoteAddr,
		"createTime":  portInfo.CreateTime,
	}, nil
}

// GetPortHistory gets port usage history
func (a *PortAPI) GetPortHistory(port int) ([]map[string]interface{}, error) {
	usages, err := a.usageTracker.GetPortHistory(port)
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
