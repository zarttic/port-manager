package model

import "time"

// PortInfo represents port information
type PortInfo struct {
	Port        int       `json:"port"`
	Protocol    string    `json:"protocol"`   // tcp/udp
	State       string    `json:"state"`      // LISTEN/ESTABLISHED
	PID         int       `json:"pid"`
	ProcessName string    `json:"processName"`
	LocalAddr   string    `json:"localAddr"`
	RemoteAddr  string    `json:"remoteAddr"`
	CreateTime  time.Time `json:"createTime"`
}

// ProcessInfo represents process information
type ProcessInfo struct {
	PID        int       `json:"pid"`
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	User       string    `json:"user"`
	Memory     uint64    `json:"memory"`
	CPU        float64   `json:"cpu"`
	CreateTime time.Time `json:"createTime"`
}

// PortUsage represents port usage record
type PortUsage struct {
	ID          int64      `json:"id"`
	Port        int        `json:"port"`
	Protocol    string     `json:"protocol"`
	PID         int        `json:"pid"`
	ProcessName string     `json:"processName"`
	ProcessPath string     `json:"processPath"`
	StartTime   time.Time  `json:"startTime"`
	EndTime     *time.Time `json:"endTime"`
	Duration    int64      `json:"duration"` // seconds
}

// PortStats represents port statistics
type PortStats struct {
	Port         int             `json:"port"`
	TotalUsage   int64           `json:"totalUsage"` // Total usage time in seconds
	UsageCount   int             `json:"usageCount"` // Number of times used
	LastUsed     time.Time       `json:"lastUsed"`
	TopProcesses []ProcessUsage  `json:"topProcesses"`
}

// ProcessUsage represents process usage statistics
type ProcessUsage struct {
	ProcessName string `json:"processName"`
	UsageCount  int    `json:"usageCount"`
	TotalTime   int64  `json:"totalTime"` // Total time in seconds
}
