package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"sovana/internal/metrics"
)

// Agent ships locally collected metrics to the backend using a worker pool.
type Agent struct {
	endpoint string
	client   *http.Client
	workers  int
	logger   *log.Logger
}

// New creates a new Agent.
func New(endpoint string, workers int, logger *log.Logger) *Agent {
	if logger == nil {
		logger = log.Default()
	}
	if workers < 1 {
		workers = 2
	}
	return &Agent{
		endpoint: endpoint,
		client:   &http.Client{Timeout: 5 * time.Second},
		workers:  workers,
		logger:   logger,
	}
}

// Run consumes metrics from the channel and delivers them concurrently.
func (a *Agent) Run(ctx context.Context, metricsCh <-chan metrics.Metric) {
	jobs := make(chan metrics.Metric)
	var wg sync.WaitGroup

	// start workers
	for i := 0; i < a.workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for m := range jobs {
				if err := a.sendMetric(ctx, m); err != nil {
					a.logger.Printf("worker %d: %v", id, err)
				}
			}
		}(i)
	}

	for {
		select {
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			return
		case m, ok := <-metricsCh:
			if !ok {
				close(jobs)
				wg.Wait()
				return
			}
			jobs <- m
		}
	}
}

func (a *Agent) sendMetric(ctx context.Context, m metrics.Metric) error {
	payload, err := json.Marshal(m)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status %s", resp.Status)
	}
	return nil
}
