package main

import (
	"log/slog"
	"net/http"
	"os"
	"restic-stats-exporter/snapshot"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	slog.Info("Starting restic statistics exporter...")
	checkEnv("RESTIC_REPOSITORY")

	resticExecutablePath := getEnvWithDefault("RSE_RESTIC_EXECUTABLE_PATH", "restic")

	prometheus.MustRegister(snapshot.NewSnapshotCollector(resticExecutablePath))

	addr := ":2112"
	slog.Info("Starting metrics HTTP server", "addr", addr)

	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error("HTTP server failed", "error", err)
		os.Exit(1)
	}
}

func checkEnv(name string) {
	val, ok := os.LookupEnv(name)
	if !ok {
		slog.Error("Environment variable is not set", "name", name)
		os.Exit(1)
	}

	if val == "" {
		slog.Error("Environment variable is empty", "name", name)
		os.Exit(1)
	}
}

func getEnvWithDefault(name string, defaultValue string) string {
	val, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return val
}
