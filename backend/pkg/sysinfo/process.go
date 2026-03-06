package sysinfo

import (
	"fmt"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v3/process"

	"port-manager/internal/model"
)

// ProcessMgr provides process management functionality
type ProcessMgr struct{}

// NewProcessMgr creates a new ProcessMgr instance
func NewProcessMgr() *ProcessMgr {
	return &ProcessMgr{}
}

// GetProcess gets detailed information about a process
func (p *ProcessMgr) GetProcess(pid int) (*model.ProcessInfo, error) {
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, fmt.Errorf("process not found: %w", err)
	}

	name, _ := proc.Name()
	path, _ := proc.Exe()
	user, _ := proc.Username()

	// Memory info
	memInfo, _ := proc.MemoryInfo()
	var memory uint64
	if memInfo != nil {
		memory = memInfo.RSS
	}

	// CPU percent
	cpuPercent, _ := proc.CPUPercent()

	// Create time
	createTime, _ := proc.CreateTime()

	return &model.ProcessInfo{
		PID:        pid,
		Name:       name,
		Path:       path,
		User:       user,
		Memory:     memory,
		CPU:        cpuPercent,
		CreateTime: time.Unix(0, createTime*int64(time.Millisecond)),
	}, nil
}

// KillProcess terminates a process by PID
func (p *ProcessMgr) KillProcess(pid int) error {
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return fmt.Errorf("process not found: %w", err)
	}

	// Try graceful termination first
	if err := proc.Terminate(); err != nil {
		// If graceful termination fails, force kill
		if err := proc.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}
	}

	return nil
}

// ListProcesses lists all running processes
func (p *ProcessMgr) ListProcesses() ([]model.ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to list processes: %w", err)
	}

	var processes []model.ProcessInfo
	for _, proc := range procs {
		name, _ := proc.Name()
		pid := int(proc.Pid)

		// Memory info
		memInfo, _ := proc.MemoryInfo()
		var memory uint64
		if memInfo != nil {
			memory = memInfo.RSS
		}

		// CPU percent
		cpuPercent, _ := proc.CPUPercent()

		processes = append(processes, model.ProcessInfo{
			PID:    pid,
			Name:   name,
			Memory: memory,
			CPU:    cpuPercent,
		})
	}

	return processes, nil
}
