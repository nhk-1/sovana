package collector

import (
	"bufio"
	"errors"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"sovana/internal/metrics"
)

// Collector periodically reads system state and emits metrics.
type Collector struct {
	interval time.Duration
	out      chan metrics.Metric
	quit     chan struct{}
}

// New creates a collector producing metrics on the provided channel.
func New(interval time.Duration) *Collector {
	return &Collector{
		interval: interval,
		out:      make(chan metrics.Metric),
		quit:     make(chan struct{}),
	}
}

// Start launches metric sampling in a background goroutine.
func (c *Collector) Start() <-chan metrics.Metric {
	go func() {
		defer close(c.out)
		hostname, _ := os.Hostname()
		for {
			select {
			case <-c.quit:
				return
			case <-time.After(c.interval):
				metric, err := collectSnapshot(hostname)
				if err != nil {
					// Metric collection errors are non-fatal; we skip this tick.
					continue
				}
				c.out <- metric
			}
		}
	}()
	return c.out
}

// Stop signals the collector to stop.
func (c *Collector) Stop() { close(c.quit) }

func collectSnapshot(hostname string) (metrics.Metric, error) {
	cpuUsage, err := readCPUUsage()
	if err != nil {
		return metrics.Metric{}, err
	}
	memUsed, memTotal, err := readMemory()
	if err != nil {
		return metrics.Metric{}, err
	}
	uptime, err := readUptime()
	if err != nil {
		return metrics.Metric{}, err
	}

	return metrics.Metric{
		Hostname:      hostname,
		CPUUsage:      cpuUsage,
		MemoryUsed:    memUsed,
		MemoryTotal:   memTotal,
		UptimeSeconds: uptime,
		Timestamp:     time.Now().UTC(),
	}, nil
}

func readCPUUsage() (float64, error) {
	firstTotal, firstIdle, err := readCPUStat()
	if err != nil {
		return 0, err
	}
	time.Sleep(200 * time.Millisecond)
	secondTotal, secondIdle, err := readCPUStat()
	if err != nil {
		return 0, err
	}

	totalDelta := secondTotal - firstTotal
	idleDelta := secondIdle - firstIdle
	if totalDelta <= 0 {
		return 0, errors.New("invalid CPU delta")
	}
	usage := (1.0 - float64(idleDelta)/float64(totalDelta)) * 100
	if usage < 0 {
		usage = 0
	}
	if usage > 100 {
		usage = 100
	}
	return usage, nil
}

func readCPUStat() (total uint64, idle uint64, err error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return 0, 0, errors.New("/proc/stat missing")
	}
	fields := strings.Fields(scanner.Text())
	if len(fields) < 5 || fields[0] != "cpu" {
		return 0, 0, errors.New("unexpected /proc/stat format")
	}

	var values []uint64
	for _, val := range fields[1:] {
		n, parseErr := strconv.ParseUint(val, 10, 64)
		if parseErr != nil {
			return 0, 0, parseErr
		}
		values = append(values, n)
	}

	for _, v := range values {
		total += v
	}
	idle = values[3] // idle time
	return total, idle, nil
}

func readMemory() (used uint64, total uint64, err error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var memTotal, memAvailable uint64
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "MemTotal:":
			memTotal, err = parseKiB(fields[1])
		case "MemAvailable:":
			memAvailable, err = parseKiB(fields[1])
		}
		if err != nil {
			return 0, 0, err
		}
		if memTotal > 0 && memAvailable > 0 {
			break
		}
	}
	if memTotal == 0 {
		return 0, 0, errors.New("unable to read MemTotal")
	}
	memUsed := memTotal - memAvailable
	return memUsed, memTotal, nil
}

func parseKiB(value string) (uint64, error) {
	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return val * 1024, nil // convert to bytes
}

func readUptime() (float64, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return 0, errors.New("unexpected uptime format")
	}
	return strconv.ParseFloat(fields[0], 64)
}

// NumCPU exposes runtime.NumCPU to keep collector self-contained.
func NumCPU() int { return runtime.NumCPU() }
