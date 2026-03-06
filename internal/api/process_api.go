package api

import (
	"port-manager/internal/service"
)

// ProcessAPI provides process-related API endpoints
type ProcessAPI struct {
	processMgr *service.ProcessManager
}

// NewProcessAPI creates a new ProcessAPI
func NewProcessAPI(processMgr *service.ProcessManager) *ProcessAPI {
	return &ProcessAPI{
		processMgr: processMgr,
	}
}

// GetProcess gets process information
func (a *ProcessAPI) GetProcess(pid int) (map[string]interface{}, error) {
	proc, err := a.processMgr.GetProcess(pid)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"pid":        proc.PID,
		"name":       proc.Name,
		"path":       proc.Path,
		"user":       proc.User,
		"memory":     proc.Memory,
		"cpu":        proc.CPU,
		"createTime": proc.CreateTime,
	}, nil
}

// KillProcess kills a process
func (a *ProcessAPI) KillProcess(pid int) error {
	return a.processMgr.KillProcess(pid)
}

// GetProcessPorts gets all ports used by a process
func (a *ProcessAPI) GetProcessPorts(pid int) ([]map[string]interface{}, error) {
	ports, err := a.processMgr.GetProcessPorts(pid)
	if err != nil {
		return nil, err
	}

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
		}
	}

	return result, nil
}
