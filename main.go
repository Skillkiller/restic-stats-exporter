package main

import (
	"log"
	"net/http"
	"os"
	"restic-stats-exporter/snapshot"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.Println("Starting restic statistics exporter...")
	checkEnv("RESTIC_REPOSITORY")

	prometheus.MustRegister(&snapshot.Collector{})

	addr := ":2112"
	log.Printf("Starting metrics HTTP server on %s ...", addr)

	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func checkEnv(name string) {
	val, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("%s environment variable is not set", name)
	}

	if val == "" {
		log.Fatalf("%s environment variable is empty", name)
	}
}
