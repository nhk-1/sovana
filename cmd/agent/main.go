package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sovana/internal/agent"
	"sovana/internal/collector"
)

func main() {
	server := flag.String("server", "http://localhost:8080/metrics", "metrics endpoint")
	interval := flag.Duration("interval", 5*time.Second, "collection interval")
	workers := flag.Int("workers", 4, "number of HTTP workers")
	flag.Parse()

	logger := log.New(os.Stdout, "agent ", log.LstdFlags)

	collector := collector.New(*interval)
	metricsCh := collector.Start()

	ag := agent.New(*server, *workers, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go ag.Run(ctx, metricsCh)

	<-ctx.Done()
	collector.Stop()
	logger.Println("agent stopped")
}
