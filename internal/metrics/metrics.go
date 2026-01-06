package metrics

import "time"

// Metric represents system health at a point in time.
type Metric struct {
	Hostname      string    `json:"hostname"`
	CPUUsage      float64   `json:"cpu_usage"`
	MemoryUsed    uint64    `json:"memory_used"`
	MemoryTotal   uint64    `json:"memory_total"`
	UptimeSeconds float64   `json:"uptime_seconds"`
	Timestamp     time.Time `json:"timestamp"`
}

// Validate ensures all fields are sane before ingesting.
func (m *Metric) Validate() error {
	if m.CPUUsage < 0 || m.CPUUsage > 100 {
		return ErrInvalidCPU
	}
	if m.MemoryUsed > m.MemoryTotal {
		return ErrInvalidMemory
	}
	if m.UptimeSeconds < 0 {
		return ErrInvalidUptime
	}
	return nil
}

var (
	// ErrInvalidCPU is returned when CPU percentage is outside [0,100].
	ErrInvalidCPU = Err("cpu_usage must be between 0 and 100")
	// ErrInvalidMemory is returned when memory numbers are inconsistent.
	ErrInvalidMemory = Err("memory_used cannot exceed memory_total")
	// ErrInvalidUptime is returned when uptime is negative.
	ErrInvalidUptime = Err("uptime_seconds cannot be negative")
)

// Err is a lightweight error type for validation errors.
type Err string

func (e Err) Error() string { return string(e) }
