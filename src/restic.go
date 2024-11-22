package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type StatsResponse struct {
	TotalSize      int `json:"total_size"`
	TotalFileCount int `json:"total_file_count"`
}

type SnapshotResponse struct {
	Time string `json:"time"`
}

type Restic struct {
	Binary string
	Name   string
	Env    map[string]string
}

func (restic Restic) Run(arguments []string, target interface{}) error {
	arguments = append(arguments, "--json")

	log.Printf("[%s] %s %s", restic.Name, restic.Binary, arguments)
	command := exec.Command(restic.Binary, arguments...)
	command.Env = os.Environ()

	for key, value := range restic.Env {
		command.Env = append(
			command.Env,
			fmt.Sprintf("%s=%s", key, value),
		)
	}

	output, err := command.Output()
	if err != nil {
		return err
	}

	err = json.Unmarshal(output, target)
	if err != nil {
		return err
	}

	return nil
}

func (restic Restic) SnapshotTimestamp() (int64, error) {
	snapshots := make([]SnapshotResponse, 0)

	args := []string{"snapshots", "latest"}
	if *tag != "" {
		args = append(args, "--tag", *tag)
	}

	err := restic.Run(args, &snapshots)
	if err != nil {
		return -1, err
	}

	if len(snapshots) == 0 {
		return 0, fmt.Errorf("no snapshots found")
	}

	time, err := time.Parse(time.RFC3339Nano, snapshots[0].Time)
	return time.Unix(), err
}
