package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	name         = flag.String("name", "", "Name of the metric (required)")
	envFile      = flag.String("env", ".env", "Path to env file (defaults to .env)")
	output       = flag.String("output", "stats.txt", "File to export the stats to")
	resticBinary = flag.String("restic-bin", "restic", "Location of the restic binary to use (defaults to loading the one in your PATH)")
	errorLog     *log.Logger
)

func init() {
	flag.Parse()

	log.SetOutput(os.Stdout)
	errorLog = log.New(os.Stderr, "", log.LstdFlags)
}

func collectMetrics(name string, env map[string]string) *prometheus.Registry {
	snapshot := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "restic_snapshot_timestamp",
	}, []string{"name"})
	registry := prometheus.NewRegistry()
	registry.Register(snapshot)

	restic := Restic{
		Binary: *resticBinary,
		Name:   name,
		Env:    env,
	}
	timestamp, err := restic.SnapshotTimestamp()
	if err != nil {
		errorLog.Printf("[%s] <ERR> %s", name, err)
	}
	snapshot.WithLabelValues(name).Set(float64(timestamp))

	return registry
}

func main() {
	if *name == "" {
		errorLog.Fatal("Name must be specified")
	}

	env, err := godotenv.Read(*envFile)
	if err != nil {
		errorLog.Fatal(err)
	}

	registry := collectMetrics(*name, env)
	err = prometheus.WriteToTextfile(*output, registry)
	if err != nil {
		errorLog.Fatal(err)
	}
}
