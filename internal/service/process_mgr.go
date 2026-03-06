package service

import (
	"port-manager/internal/model"
	"port-manager/pkg/sysinfo"
)

// ProcessManager provides process management functionality
type ProcessManager struct {
	procMgr *sysinfo.ProcessMgr
}

// NewProcessManager creates a new ProcessManager
func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		procMgr: sysinfo.NewProcessMgr(),
	}
}

// GetProcess gets process information
func (m *ProcessManager) GetProcess(pid int) (*model.ProcessInfo, error) {
	return m.procMgr.GetProcess(pid)
}

// KillProcess kills a process
func (m *ProcessManager) KillProcess(pid int) error {
	return m.procMgr.KillProcess(pid)
}

// GetProcessPorts gets all ports used by a process
func (m *ProcessManager) GetProcessPorts(pid int) ([]model.PortInfo, error) {
	netstat := sysinfo.NewNetStat()
	return netstat.GetProcessPorts(pid)
}

// ListProcesses lists all running processes
func (m *ProcessManager) ListProcesses() ([]model.ProcessInfo, error) {
	return m.procMgr.ListProcesses()
}
