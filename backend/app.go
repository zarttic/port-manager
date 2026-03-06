package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"port-manager/internal/api"
	"port-manager/internal/repository"
	"port-manager/internal/service"
)

// App struct
type App struct {
	ctx          context.Context
	db           *sql.DB
	portAPI      *api.PortAPI
	processAPI   *api.ProcessAPI
	statsAPI     *api.StatsAPI
	scanner      *service.PortScanner
	processMgr   *service.ProcessManager
	usageTracker *service.UsageTracker
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	if err := a.initDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repositories
	portRepo := repository.NewPortRepository(a.db)
	usageRepo := repository.NewUsageRepository(a.db)

	// Initialize services
	a.scanner = service.NewPortScanner()
	a.processMgr = service.NewProcessManager()
	a.usageTracker = service.NewUsageTracker(portRepo, usageRepo)

	// Initialize APIs
	a.portAPI = api.NewPortAPI(a.scanner, a.usageTracker)
	a.processAPI = api.NewProcessAPI(a.processMgr)
	a.statsAPI = api.NewStatsAPI(usageRepo)

	// Start background tracking
	go a.usageTracker.StartTracking(ctx)

	log.Println("Application started successfully")
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	// Stop tracking
	a.usageTracker.StopTracking()

	// Close database
	if a.db != nil {
		a.db.Close()
	}

	log.Println("Application shutdown complete")
}

// initDatabase initializes SQLite database
func (a *App) initDatabase() error {
	// Get user data directory
	dataDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config dir: %w", err)
	}

	// Create app-specific directory
	appDir := filepath.Join(dataDir, "port-manager")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("failed to create app dir: %w", err)
	}

	// Open database
	dbPath := filepath.Join(appDir, "port-manager.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)

	a.db = db

	// Run migrations
	return a.runMigrations()
}

// runMigrations executes database migrations
func (a *App) runMigrations() error {
	// Port usage table
	_, err := a.db.Exec(`
		CREATE TABLE IF NOT EXISTS port_usage (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			port INTEGER NOT NULL,
			protocol TEXT NOT NULL CHECK(protocol IN ('tcp', 'udp')),
			pid INTEGER NOT NULL,
			process_name TEXT NOT NULL,
			process_path TEXT,
			start_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			end_time DATETIME,
			duration INTEGER,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create port_usage table: %w", err)
	}

	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_port_usage_port ON port_usage(port);",
		"CREATE INDEX IF NOT EXISTS idx_port_usage_pid ON port_usage(pid);",
		"CREATE INDEX IF NOT EXISTS idx_port_usage_start_time ON port_usage(start_time);",
		"CREATE INDEX IF NOT EXISTS idx_port_usage_process ON port_usage(process_name);",
	}

	for _, indexSQL := range indexes {
		if _, err := a.db.Exec(indexSQL); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	// Scan history table
	_, err = a.db.Exec(`
		CREATE TABLE IF NOT EXISTS scan_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			scan_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			port_count INTEGER,
			duration_ms INTEGER
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create scan_history table: %w", err)
	}

	log.Println("Database migrations completed")
	return nil
}

// Exposed API methods

// ScanPorts scans all ports
func (a *App) ScanPorts() ([]map[string]interface{}, error) {
	return a.portAPI.ScanPorts()
}

// ScanPort scans a specific port
func (a *App) ScanPort(port int) (map[string]interface{}, error) {
	return a.portAPI.ScanPort(port)
}

// GetPortHistory gets port usage history
func (a *App) GetPortHistory(port int) ([]map[string]interface{}, error) {
	return a.portAPI.GetPortHistory(port)
}

// GetProcess gets process information
func (a *App) GetProcess(pid int) (map[string]interface{}, error) {
	return a.processAPI.GetProcess(pid)
}

// KillProcess kills a process by PID
func (a *App) KillProcess(pid int) error {
	return a.processAPI.KillProcess(pid)
}

// GetProcessPorts gets all ports used by a process
func (a *App) GetProcessPorts(pid int) ([]map[string]interface{}, error) {
	return a.processAPI.GetProcessPorts(pid)
}

// GetPortStats gets statistics for a port
func (a *App) GetPortStats(port int) (map[string]interface{}, error) {
	return a.statsAPI.GetPortStats(port)
}

// GetTopUsedPorts gets top used ports
func (a *App) GetTopUsedPorts(limit int) ([]map[string]interface{}, error) {
	return a.statsAPI.GetTopUsedPorts(limit)
}

// GetUsageHistory gets usage history within a time range
func (a *App) GetUsageHistory(start, end string) ([]map[string]interface{}, error) {
	return a.statsAPI.GetUsageHistory(start, end)
}
