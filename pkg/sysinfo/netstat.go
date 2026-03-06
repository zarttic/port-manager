package sysinfo

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"

	"port-manager/internal/model"
)

// NetStat provides port scanning functionality
type NetStat struct{}

// NewNetStat creates a new NetStat instance
func NewNetStat() *NetStat {
	return &NetStat{}
}

// ScanPorts scans all listening and established ports
func (n *NetStat) ScanPorts() ([]model.PortInfo, error) {
	var ports []model.PortInfo

	// Get all network connections
	connections, err := net.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("failed to get network connections: %w", err)
	}

	for _, conn := range connections {
		// Skip empty connections
		if conn.Laddr.Port == 0 {
			continue
		}

		// Get process name
		processName := ""
		if conn.Pid > 0 {
			if p, err := process.NewProcess(int32(conn.Pid)); err == nil {
				processName, _ = p.Name()
			}
		}

		// Determine protocol
		protocol := "tcp"
		if conn.Type == syscall.SOCK_DGRAM {
			protocol = "udp"
		}

		portInfo := model.PortInfo{
			Port:        int(conn.Laddr.Port),
			Protocol:    protocol,
			State:       conn.Status,
			PID:         int(conn.Pid),
			ProcessName: processName,
			LocalAddr:   conn.Laddr.IP,
			RemoteAddr:  conn.Raddr.IP,
		}

		// Get process create time
		if conn.Pid > 0 {
			if p, err := process.NewProcess(int32(conn.Pid)); err == nil {
				if createTime, err := p.CreateTime(); err == nil {
					portInfo.CreateTime = time.Unix(0, createTime*int64(time.Millisecond))
				}
			}
		}

		ports = append(ports, portInfo)
	}

	return ports, nil
}

// ScanPort scans a specific port
func (n *NetStat) ScanPort(port int) (*model.PortInfo, error) {
	connections, err := net.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("failed to get network connections: %w", err)
	}

	for _, conn := range connections {
		if int(conn.Laddr.Port) == port {
			// Get process name
			processName := ""
			if conn.Pid > 0 {
				if p, err := process.NewProcess(int32(conn.Pid)); err == nil {
					processName, _ = p.Name()
				}
			}

			protocol := "tcp"
			if conn.Type == syscall.SOCK_DGRAM {
				protocol = "udp"
			}

			return &model.PortInfo{
				Port:        int(conn.Laddr.Port),
				Protocol:    protocol,
				State:       conn.Status,
				PID:         int(conn.Pid),
				ProcessName: processName,
				LocalAddr:   conn.Laddr.IP,
				RemoteAddr:  conn.Raddr.IP,
			}, nil
		}
	}

	return nil, fmt.Errorf("port %d not found", port)
}

// GetProcessPorts gets all ports used by a process
func (n *NetStat) GetProcessPorts(pid int) ([]model.PortInfo, error) {
	var ports []model.PortInfo

	connections, err := net.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("failed to get network connections: %w", err)
	}

	for _, conn := range connections {
		if int(conn.Pid) == pid && conn.Laddr.Port > 0 {
			processName := ""
			if p, err := process.NewProcess(int32(conn.Pid)); err == nil {
				processName, _ = p.Name()
			}

			protocol := "tcp"
			if conn.Type == syscall.SOCK_DGRAM {
				protocol = "udp"
			}

			ports = append(ports, model.PortInfo{
				Port:        int(conn.Laddr.Port),
				Protocol:    protocol,
				State:       conn.Status,
				PID:         int(conn.Pid),
				ProcessName: processName,
				LocalAddr:   conn.Laddr.IP,
				RemoteAddr:  conn.Raddr.IP,
			})
		}
	}

	return ports, nil
}
